package rel

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/arr-ai/hash/hash128"
)

var (
	truePosRel  = &positionalRelation{store: &colStore{width: 0, n: 1}, n: 1}
	falsePosRel = &positionalRelation{store: &colStore{width: 0}}

	posRelSalt = hash128.String("github.com/arr-ai/arrai/rel.positionalRelation")
)

// positionalRelation is a set of distinct rows, all the same width: a view
// of a colStore. sel == nil views the store's first n committed rows; a
// non-nil sel views the rows at those arena ids, kept strictly ascending so
// two views of one store compare by their id sequences and membership can
// binary-search. Group-by indexes and the set hash are computed lazily per
// view and memoised in meta.
type positionalRelation struct {
	store *colStore
	arena []Value // snapshot; rows below this view's reach never change
	n     int
	sel   []uint32
	// parent is the view this one was derived from (Where/Without/selView).
	// A group-by on this view filters parent's memoised index when present
	// instead of scanning the store again (🎯T18).
	parent *positionalRelation

	once sync.Once
	meta *positionalRelationMetadata
}

type positionalRelationMetadata struct {
	sync.Mutex
	groups map[string]*groupIndex
	hash   hash128.H128
	hashed bool
	// keys are memoKeys of projectors that grouped into n distinct buckets:
	// empirically discovered candidate keys (🎯T20).
	keys []string
	// zones[i] is the cached numeric min/max of column i, or nil if not yet
	// computed. A stored zoneMap with ok=false means the column is not numeric.
	zones []*zoneMap
	// shapeHashes memoises shapeHash per tuple shape: the layout-independent
	// half of Relation.Hash128.
	shapeHashes map[*Shape]hash128.H128
	// plans is the S5 fact-keyed cache of index-answered where results.
	plans    map[string]planEntry
	planHits int
}

type planEntry struct {
	guard func() bool
	value Value
}

func (r *positionalRelation) getMeta() *positionalRelationMetadata {
	r.once.Do(func() {
		r.meta = &positionalRelationMetadata{}
	})
	return r.meta
}

// newPositionalRelation builds a relation from rows, dropping duplicates.
func newPositionalRelation(width int, rows ...Values) *positionalRelation {
	b := newStoreBuilder(width, len(rows), true)
	for _, v := range rows {
		b.add(v)
	}
	return b.finish()
}

// rowAt returns the i'th row of the view, as a slice of the arena.
func (r *positionalRelation) rowAt(i int) Values {
	if r.sel != nil {
		i = int(r.sel[i])
	}
	return rowOf(r.arena, r.store.width, i)
}

// arenaID returns the arena id of the i'th row of the view.
func (r *positionalRelation) arenaID(i int) uint32 {
	if r.sel != nil {
		return r.sel[i]
	}
	return uint32(i)
}

// selView views the same store restricted to the given ascending arena ids,
// which must all be visible in this view. The slice is retained.
func (r *positionalRelation) selView(ids []uint32) *positionalRelation {
	return &positionalRelation{store: r.store, arena: r.arena, n: len(ids), sel: ids, parent: r}
}

func (r *positionalRelation) Count() int {
	return r.n
}

func (r *positionalRelation) Width() int {
	return r.store.width
}

func (r *positionalRelation) IsEmpty() bool {
	return r.n == 0
}

func (r *positionalRelation) IsTrue() bool {
	return r.n != 0
}

func (r *positionalRelation) IsLiteralTrue() bool {
	return r.n == 1 && r.store.width == 0
}

func (r *positionalRelation) Has(v Values) bool {
	_, has := r.find(v)
	return has
}

// find returns the arena id of the view's row equal to v.
func (r *positionalRelation) find(v Values) (uint32, bool) {
	if len(v) != r.store.width || r.n == 0 {
		return 0, false
	}
	if r.store.width == 0 {
		return 0, true // the only width-0 row is the empty row, and n > 0
	}
	s := r.store
	h := v.Hash128()
	s.mu.Lock()
	s.ensureIndexLocked()
	id, ok := s.findLocked(v, h, len(r.arena)/s.width)
	s.mu.Unlock()
	if !ok {
		return 0, false
	}
	if r.sel != nil {
		if _, found := slices.BinarySearch(r.sel, id); !found {
			return 0, false
		}
	}
	return id, true
}

// Hash128 returns the set hash of the view: the xor of its row hashes, which
// the store computes once per row.
func (r *positionalRelation) Hash128() hash128.H128 {
	m := r.getMeta()
	m.Lock()
	defer m.Unlock()
	if !m.hashed {
		h := posRelSalt
		s := r.store
		s.mu.Lock()
		s.ensureIndexLocked()
		for i := 0; i < r.n; i++ {
			h = h.Xor(s.hashes[r.arenaID(i)])
		}
		s.mu.Unlock()
		m.hash, m.hashed = h, true
	}
	return m.hash
}

// shapeHash returns the xor over the view's rows of each row's
// per-attribute hash under sh, taking row values through layout (layout[i]
// is the row position of sh's i'th attribute). It depends only on each
// row's name/value pairs — not on the relation's internal attribute order —
// which is what lets Relation.Hash128 agree with EqualRelation across
// layouts. Computed once per shape per view.
func (r *positionalRelation) shapeHash(sh *Shape, layout []int) hash128.H128 {
	m := r.getMeta()
	m.Lock()
	defer m.Unlock()
	if h, has := m.shapeHashes[sh]; has {
		return h
	}
	var h hash128.H128
	for i := 0; i < r.n; i++ {
		row := r.rowAt(i)
		var rh hash128.H128
		for k, j := range layout {
			rh = rh.Xor(hashAttr(sh.nameH[k], row[j]))
		}
		h = h.Xor(rh)
	}
	if m.shapeHashes == nil {
		m.shapeHashes = map[*Shape]hash128.H128{}
	}
	m.shapeHashes[sh] = h
	return h
}

func (r *positionalRelation) EqualPositionalRelation(r2 *positionalRelation) bool {
	if r.n != r2.n || r.store.width != r2.store.width {
		return false
	}
	if r.store == r2.store {
		// Store rows are distinct, so equal rows have equal ids, and both id
		// sequences ascend: the views are equal iff the sequences are.
		for i := 0; i < r.n; i++ {
			if r.arenaID(i) != r2.arenaID(i) {
				return false
			}
		}
		return true
	}
	if r.Hash128() != r2.Hash128() {
		return false
	}
	for i := 0; i < r.n; i++ {
		if !r2.Has(r.rowAt(i)) {
			return false
		}
	}
	return true
}

// With returns a view including row v. When this view is the frontier of its
// store — nothing appended beyond it — the row is appended in place, so a
// chain of Withs extends one arena instead of copying per step.
func (r *positionalRelation) With(v Values) *positionalRelation {
	if r.store.width != len(v) {
		panic(fmt.Errorf("positionalRelation.With: row width %d != %d", len(v), r.store.width))
	}
	if r.store.width == 0 {
		return truePosRel // the empty row is the only possible row
	}
	if r.sel == nil {
		s := r.store
		h := v.Hash128()
		s.mu.Lock()
		s.ensureIndexLocked()
		if _, found := s.findLocked(v, h, r.n); found {
			s.mu.Unlock()
			return r
		}
		if s.n == r.n {
			s.arena = append(s.arena, v...)
			s.hashes = append(s.hashes, h)
			s.insertLocked(h, uint32(s.n))
			s.n++
			nr := &positionalRelation{store: s, arena: s.arena[:s.n*s.width], n: s.n}
			s.mu.Unlock()
			return nr
		}
		s.mu.Unlock()
	} else if _, found := r.find(v); found {
		return r
	}
	// A sibling view has already extended the store, or this view is a
	// selection: copy out into a fresh store, which the result then owns, so
	// further Withs on it extend in place. Rows are distinct and v is known
	// absent.
	w := r.store.width
	arena := make([]Value, 0, (r.n+1)*w)
	if r.sel == nil {
		arena = append(arena, r.arena[:r.n*w]...)
	} else {
		for i := 0; i < r.n; i++ {
			arena = append(arena, r.rowAt(i)...)
		}
	}
	arena = append(arena, v...)
	return &positionalRelation{
		store: &colStore{width: w, arena: arena, n: r.n + 1},
		arena: arena,
		n:     r.n + 1,
	}
}

// Without returns a view excluding row v: the same store behind a selection.
func (r *positionalRelation) Without(v Values) *positionalRelation {
	id, found := r.find(v)
	if !found {
		return r
	}
	if r.store.width == 0 {
		return falsePosRel
	}
	sel := make([]uint32, 0, r.n-1)
	for i := 0; i < r.n; i++ {
		if a := r.arenaID(i); a != id {
			sel = append(sel, a)
		}
	}
	return r.selView(sel)
}

// Where returns the view of rows satisfying p: the same store behind a
// selection, with no rows copied. Large views evaluate the predicate in
// parallel; concatenating the per-range selections in range order keeps the
// selection ascending.
func (r *positionalRelation) Where(p func(Values) (bool, error)) (*positionalRelation, error) {
	if fastPaths {
		if ranges := parallelRanges(r.n); ranges != nil {
			return r.whereParallel(ranges, p)
		}
	}
	sel := make([]uint32, 0, r.n)
	for i := 0; i < r.n; i++ {
		match, err := p(r.rowAt(i))
		if err != nil {
			return nil, err
		}
		if match {
			sel = append(sel, r.arenaID(i))
		}
	}
	if len(sel) == r.n {
		return r, nil
	}
	return r.selView(sel), nil
}

func (r *positionalRelation) whereParallel(
	ranges [][2]int, p func(Values) (bool, error),
) (*positionalRelation, error) {
	sels := make([][]uint32, len(ranges))
	errs := make([]error, len(ranges))
	runRanges(ranges, func(w, lo, hi int) {
		sel := make([]uint32, 0, hi-lo)
		for i := lo; i < hi; i++ {
			match, err := p(r.rowAt(i))
			if err != nil {
				errs[w] = err
				return
			}
			if match {
				sel = append(sel, r.arenaID(i))
			}
		}
		sels[w] = sel
	})
	if err := firstErr(errs); err != nil {
		return nil, err
	}
	total := 0
	for _, s := range sels {
		total += len(s)
	}
	if total == r.n {
		return r, nil
	}
	sel := make([]uint32, 0, total)
	for _, s := range sels {
		sel = append(sel, s...)
	}
	return r.selView(sel), nil
}

// Project returns rows projected through p. A projection that visits every
// column exactly once permutes distinct rows into distinct rows; any other
// projection deduplicates.
func (r *positionalRelation) Project(p valueProjector) *positionalRelation {
	if p.isIdentity(r.store.width) {
		return r
	}
	dedup := !p.isPermutation(r.store.width)
	if dedup && fastPaths && r.groupBy(p).count() == r.n {
		dedup = false
	}
	b := newStoreBuilder(len(p), r.n, dedup)
	for i := 0; i < r.n; i++ {
		b.addProjection(r.rowAt(i), p)
	}
	return b.finish()
}

func (r *positionalRelation) Map(f func(Values) (Value, error)) (Set, error) {
	sb := NewSetBuilder()
	for i := 0; i < r.n; i++ {
		val, err := f(r.rowAt(i))
		if err != nil {
			return nil, err
		}
		sb.Add(val)
	}
	return sb.Finish()
}

func (r *positionalRelation) String() string {
	if r.n == 0 {
		return "{}"
	}
	sb := strings.Builder{}
	sb.WriteString("{")
	identity := make(valueProjector, r.store.width)
	for i := range identity {
		identity[i] = i
	}
	for i, e := 0, r.OrderedRange(identity); e.Next(); i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(e.Values().String())
	}
	sb.WriteString("}")
	return sb.String()
}

func (r *positionalRelation) Range() *positionalRelationValuesEnumerator {
	return &positionalRelationValuesEnumerator{r: r, i: -1}
}

// OrderedRange enumerates the view's rows ordered by their projection
// through p.
func (r *positionalRelation) OrderedRange(p valueProjector) *positionalRelationValuesEnumerator {
	order := make([]uint32, r.n)
	for i := range order {
		order[i] = uint32(i)
	}
	slices.SortFunc(order, func(i, j uint32) int {
		a, b := r.rowAt(int(i)).project(p), r.rowAt(int(j)).project(p)
		if a.Less(b) {
			return -1
		} else if b.Less(a) {
			return 1
		}
		return 0
	})
	return &positionalRelationValuesEnumerator{r: r, order: order, i: -1}
}

// groupBy returns the memoised group-by index for p over this view.
func (r *positionalRelation) groupBy(p valueProjector) *groupIndex {
	m := r.getMeta()
	key := p.memoKey()
	m.Lock()
	if g, has := m.groups[key]; has {
		m.Unlock()
		return g
	}
	m.Unlock()

	var g *groupIndex
	if fastPaths {
		for anc := r.parent; anc != nil; anc = anc.parent {
			if pg := anc.cachedGroup(p); pg != nil {
				g = pg.filterToView(r)
				break
			}
		}
	}
	if g == nil {
		g = newGroupIndex(r, p)
	}
	m.Lock()
	if existing, has := m.groups[key]; has {
		m.Unlock()
		return existing
	}
	if m.groups == nil {
		m.groups = map[string]*groupIndex{}
	}
	m.groups[key] = g
	if g.filtered {
		if g.base != nil && g.base.rel.hasCandidateKey(p) {
			m.keys = appendKey(m.keys, key)
		}
	} else if g.count() == r.n {
		m.keys = appendKey(m.keys, key)
	}
	m.Unlock()
	return g
}

func (r *positionalRelation) cachedGroup(p valueProjector) *groupIndex {
	m := r.getMeta()
	m.Lock()
	defer m.Unlock()
	if m.groups == nil {
		return nil
	}
	return m.groups[p.memoKey()]
}

func appendKey(keys []string, key string) []string {
	for _, k := range keys {
		if k == key {
			return keys
		}
	}
	return append(keys, key)
}

func (r *positionalRelation) hasArenaID(id uint32) bool {
	if r.sel != nil {
		_, ok := slices.BinarySearch(r.sel, id)
		return ok
	}
	return int(id) < r.n
}

func (r *positionalRelation) planGet(key string, guard func() bool) (Value, bool) {
	m := r.getMeta()
	m.Lock()
	defer m.Unlock()
	e, ok := m.plans[key]
	if !ok || e.value == nil || e.guard == nil || !e.guard() || !guard() {
		return nil, false
	}
	m.planHits++
	return e.value, true
}

func (r *positionalRelation) planPut(key string, guard func() bool, v Value) {
	m := r.getMeta()
	m.Lock()
	defer m.Unlock()
	if m.plans == nil {
		m.plans = map[string]planEntry{}
	}
	m.plans[key] = planEntry{guard: guard, value: v}
}

func (r *positionalRelation) planHitCount() int {
	m := r.getMeta()
	m.Lock()
	defer m.Unlock()
	return m.planHits
}

func (r *positionalRelation) seedKey(p valueProjector) {
	m := r.getMeta()
	m.Lock()
	defer m.Unlock()
	m.keys = appendKey(m.keys, p.memoKey())
}

func (r *positionalRelation) hasCandidateKey(p valueProjector) bool {
	m := r.getMeta()
	m.Lock()
	defer m.Unlock()
	key := p.memoKey()
	for _, k := range m.keys {
		if k == key {
			return true
		}
	}
	return false
}

func (p valueProjector) memoKey() string {
	var sb strings.Builder
	for i, x := range p {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(x))
	}
	return sb.String()
}

// isPermutation reports whether p visits every column of a width-wide row
// exactly once, in some order.
func (p valueProjector) isPermutation(width int) bool {
	if len(p) != width {
		return false
	}
	seen := make([]bool, width)
	for _, i := range p {
		if i < 0 || i >= width || seen[i] {
			return false
		}
		seen[i] = true
	}
	return true
}

func createMode(leftKey, rightKey, leftOutput, rightOutput valueProjector) CombineOp {
	if len(leftKey) != len(rightKey) {
		panic(fmt.Errorf("keys are not of the same length: %v and %v", leftKey, rightKey))
	}
	if (!leftKey.isSubProjection(leftOutput) && leftOutput.hasCommonIndices(leftKey)) ||
		(!rightKey.isSubProjection(rightOutput) && rightOutput.hasCommonIndices(rightKey)) {
		panic(fmt.Errorf("partial key output: %v %v %v %v", leftKey, rightKey, leftOutput, rightOutput))
	}
	var mode CombineOp
	if !leftOutput.isSubProjection(leftKey) {
		mode |= OnlyOnLHS
	}
	if !rightOutput.isSubProjection(rightKey) {
		mode |= OnlyOnRHS
	}
	// only one side should include the key indices
	if leftOutput.hasCommonIndices(leftKey) != rightOutput.hasCommonIndices(rightKey) {
		mode |= InBoth
	}
	return mode
}

func (r *positionalRelation) Join(
	r2 *positionalRelation,
	leftKey, rightKey, leftOutput, rightOutput valueProjector,
) *positionalRelation {
	mode := createMode(leftKey, rightKey, leftOutput, rightOutput)
	switch mode {
	case AllPairs, OnlyOnLHS | OnlyOnRHS: // <&>, <->
		return r.JoinKeepEverything(r2, leftKey, rightKey, leftOutput, rightOutput)
	case OnlyOnLHS, OnlyOnLHS | InBoth: // <--, <&-
		return joinOneSide(r, r2.groupBy(rightKey), leftKey, leftOutput)
	case OnlyOnRHS, OnlyOnRHS | InBoth: // -->, -&>
		return joinOneSide(r2, r.groupBy(leftKey), rightKey, rightOutput)
	case InBoth: // -&-
		return r.JoinCommonOnly(r2, leftKey, rightKey, leftOutput, rightOutput)
	case 0: // ---
		return r.JoinIfCommonExist(r2, leftKey, rightKey)
	default:
		panic(fmt.Errorf("unhandled mode %v", mode))
	}
}

// JoinKeepEverything pairs every left row with every right row sharing its
// key: <&> and <->.
func (r *positionalRelation) JoinKeepEverything(
	r2 *positionalRelation,
	leftKey, rightKey, leftOutput, rightOutput valueProjector,
) *positionalRelation {
	leftGroup, rightGroup := r.groupBy(leftKey), r2.groupBy(rightKey)
	b := newStoreBuilder(len(leftOutput)+len(rightOutput), 0, true)
	leftGroup.each(func(lb *groupBucket) {
		rb := rightGroup.bucketOfRow(leftGroup.repRow(lb), leftKey)
		if rb == nil {
			return
		}
		for _, lid := range lb.rows {
			left := leftGroup.row(lid)
			for _, rid := range rb.rows {
				b.addJoined(left, leftOutput, rightGroup.row(rid), rightOutput)
			}
		}
	})
	return b.finish()
}

// JoinIfCommonExist returns true iff the two sides share a key: ---.
func (r *positionalRelation) JoinIfCommonExist(
	r2 *positionalRelation,
	leftKey, rightKey valueProjector,
) *positionalRelation {
	if r.Count() > r2.Count() {
		r, r2 = r2, r
		leftKey, rightKey = rightKey, leftKey
	}
	group := r.groupBy(leftKey)
	for i := 0; i < r2.n; i++ {
		if group.bucketOfRow(r2.rowAt(i), rightKey) != nil {
			return truePosRel
		}
	}
	return falsePosRel
}

// JoinCommonOnly returns the shared keys, projected as one side's output:
// -&-. The output columns are a subset of the key columns, so each output
// row projects straight off a representative row of the owning side.
func (r *positionalRelation) JoinCommonOnly(
	r2 *positionalRelation,
	leftKey, rightKey, leftOutput, rightOutput valueProjector,
) *positionalRelation {
	own, other := r.groupBy(leftKey), r2.groupBy(rightKey)
	key, output := leftKey, leftOutput
	if len(output) == 0 {
		own, other = other, own
		key, output = rightKey, rightOutput
	}
	// Keys are distinct per group; the outputs stay distinct unless the
	// projection drops one of the key's columns.
	dedupe := !key.isSubProjection(output)
	b := newStoreBuilder(len(output), 0, dedupe)
	own.each(func(ob *groupBucket) {
		rep := own.repRow(ob)
		if other.bucketOfRow(rep, key) != nil {
			b.addProjection(rep, output)
		}
	})
	return b.finish()
}

// joinOneSide keeps base rows whose key the index shares: <--, <&-, -->,
// -&- variants where only one side's columns survive.
func joinOneSide(
	base *positionalRelation, intersector *groupIndex,
	key, output valueProjector,
) *positionalRelation {
	if output.isIdentity(base.Width()) {
		sel := make([]uint32, 0, base.n)
		for i := 0; i < base.n; i++ {
			if intersector.bucketOfRow(base.rowAt(i), key) != nil {
				sel = append(sel, base.arenaID(i))
			}
		}
		if len(sel) == base.n {
			return base
		}
		return base.selView(sel)
	}
	b := newStoreBuilder(len(output), 0, !output.isPermutation(base.Width()))
	for i := 0; i < base.n; i++ {
		row := base.rowAt(i)
		if intersector.bucketOfRow(row, key) != nil {
			b.addProjection(row, output)
		}
	}
	return b.finish()
}

type positionalRelationBuilder struct {
	b storeBuilder
}

func newPositionalRelationBuilder(width, capacity int) *positionalRelationBuilder {
	return &positionalRelationBuilder{b: newStoreBuilder(width, capacity, true)}
}

func (r *positionalRelationBuilder) Add(v Values) {
	r.b.add(v)
}

func (r *positionalRelationBuilder) Finish() *positionalRelation {
	return r.b.finish()
}

type positionalRelationValuesEnumerator struct {
	r     *positionalRelation
	order []uint32 // view positions when ordered; nil for natural order
	i     int
}

func (e *positionalRelationValuesEnumerator) Next() bool {
	e.i++
	return e.i < e.r.n
}

func (e *positionalRelationValuesEnumerator) Values() Values {
	if e.order != nil {
		return e.r.rowAt(int(e.order[e.i]))
	}
	return e.r.rowAt(e.i)
}
