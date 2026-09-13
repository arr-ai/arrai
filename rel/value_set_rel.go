package rel

import (
	"slices"

	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/arr-ai/wbnf/parser"

	"github.com/arr-ai/arrai/pkg/fu"
)

// Relation is a Set that only contains Tuples, all of which map the same keys.
type Relation struct {
	attrs   NamesSlice
	p       valueProjector
	rows    *positionalRelation // TODO: experiment with column table
	attrMap map[string]int      // cached mapIndices(attrs, p)

	// attrSet is the interned attribute set of every row; layout[i] is the
	// row position of attrSet.names[i]. When layout is the identity a row's
	// values already are the tuple's values, so inflating a row is a wrap,
	// not a copy.
	attrSet Names
	layout  []int
	direct  bool
}

func newRelation(attrs NamesSlice, p valueProjector, rows *positionalRelation) Relation {
	r := Relation{attrs: attrs, p: p, rows: rows, attrMap: mapIndices(attrs, p)}
	r.attrSet = internNames(attrs.GetSorted())
	r.layout = make([]int, len(r.attrSet.names))
	r.direct = true
	for i, name := range r.attrSet.names {
		r.layout[i] = r.attrMap[name]
		if r.layout[i] != i {
			r.direct = false
		}
	}
	r.seedAtKey()
	return r
}

// seedAtKey records @ as a candidate key for array/dict/string/bytes
// encodings, whose @ is unique by construction (🎯T20).
func (r Relation) seedAtKey() {
	at, has := r.attrMap["@"]
	if !has {
		return
	}
	for _, n := range []string{ArrayItemAttr, DictValueAttr, StringCharAttr, BytesByteAttr} {
		if _, ok := r.attrMap[n]; ok {
			r.rows.seedKey(valueProjector{at})
			return
		}
	}
}

// tupleInflateHook, when set, is called on every Relation.tuple inflation.
// Tests use it to prove a store-native path does not box rows.
var tupleInflateHook func()

// tuple inflates a row to a Tuple.
func (r Relation) tuple(row Values) Tuple {
	if tupleInflateHook != nil {
		tupleInflateHook()
	}
	if fastPaths && r.direct && len(row) == len(r.layout) {
		return newShapedTuple(r.attrSet, row)
	}
	vals := make([]Value, len(r.layout))
	for i, j := range r.layout {
		vals[i] = row[j]
	}
	return newShapedTuple(r.attrSet, vals)
}

func mapIndices(n NamesSlice, indices valueProjector) map[string]int {
	if len(n) != len(indices) {
		panic(fmt.Errorf("names and indices are not the same length: %v and %v", n, indices))
	}
	m := make(map[string]int, len(n))
	for i, name := range n {
		m[name] = indices[i]
	}
	return m
}

func (r Relation) newBody(rows *positionalRelation) Set {
	if !rows.IsTrue() {
		return None
	}
	r.rows = rows
	return r
}

func (r Relation) AttrsName() NamesSlice {
	return r.attrs
}

func (r Relation) getIndices(names NamesSlice) []int {
	mapping := make(map[string]int)
	for i, name := range r.attrs {
		mapping[name] = i
	}
	indices := make([]int, 0, len(names))
	for _, name := range names {
		index, has := mapping[name]
		if !has {
			panic(fmt.Errorf("name %s not found in relation %v", name, r))
		}
		indices = append(indices, index)
	}
	return indices
}

func (r Relation) Count() int {
	// TODO: handle laziness
	return r.rows.Count()
}

// hasShape reports whether t has exactly this relation's attributes.
// Names are interned, so for a GenericTuple this is an identity comparison;
// other Tuple kinds fall back to comparing the attribute sets.
func (r Relation) hasShape(t Tuple) bool {
	if g, is := t.(*GenericTuple); is && fastPaths {
		return g.attrSet() == r.attrSet
	}
	return r.attrs.EqualTupleAttrs(t)
}

// tupleToValues converts a tuple to a row in this relation's layout. A tuple
// of the relation's own shape already holds its values in shape order, so it
// only needs permuting — or nothing at all when the layout is the identity,
// in which case the row shares the tuple's values (both are immutable).
func (r Relation) tupleToValues(t Tuple) Values {
	if g, is := t.(*GenericTuple); is && fastPaths && g.attrSet() == r.attrSet {
		if r.direct {
			return Values(g.vals)
		}
		values := make(Values, len(g.vals))
		for i, j := range r.layout {
			values[j] = g.vals[i]
		}
		return values
	}
	if len(r.attrs) != t.Count() {
		panic("tupleToValues: names and values don't have the same number")
	}
	values := make(Values, r.rows.Width())
	for i, name := range r.attrs {
		values[r.p[i]] = t.MustGet(name)
	}
	return values
}

func (r Relation) Has(v Value) bool {
	if t, is := v.(Tuple); is {
		if !r.hasShape(t) {
			return false
		}
		if r.rows.Width() == len(r.attrs) {
			return r.rows.Has(r.tupleToValues(t))
		}
		key := make([]Value, len(r.p))
		for i, name := range r.attrs {
			key[i] = t.MustGet(name)
		}
		_, ok := r.rows.groupBy(r.p).getKey(key...)
		return ok
	}
	return false
}

func (r Relation) Enumerator() ValueEnumerator {
	return &relationEnumerator{
		r: r,
		i: r.rows.Range(),
	}
}

func (r Relation) OrderedValues() ValueEnumerator {
	return OrderedValueEnumerator(r.Enumerator(), ValueLess)
}

func (r Relation) ArrayEnumerator() ValueEnumerator {
	return &relationEnumerator{
		r: r,
		i: r.rows.OrderedRange(r.p),
	}
}

func (r Relation) materializeProjection() Relation {
	if r.rows.Width() == len(r.attrs) {
		return r
	}
	id := make(valueProjector, len(r.attrs))
	for i := range id {
		id[i] = i
	}
	return newRelation(r.attrs, id, r.rows.Project(r.p))
}

func (r Relation) With(v Value) Set {
	if t, is := v.(Tuple); is && r.hasShape(t) {
		r = r.materializeProjection()
		rows := r.rows.With(r.tupleToValues(t))
		if rows == r.rows {
			return r
		}
		// The result is never empty, so newBody just swaps the rows in
		// without recomputing the relation's shape metadata.
		return r.newBody(rows)
	}
	return toUnionSetWithItem(r, v)
}

func (r Relation) Without(v Value) Set {
	if t, is := v.(Tuple); is && r.hasShape(t) {
		r = r.materializeProjection()
		values := r.tupleToValues(t)
		pr := r.rows.Without(values)
		return r.newBody(pr)
	}
	return r
}

func (r Relation) Map(f func(Value) (Value, error)) (Set, error) {
	return r.rows.Map(func(v Values) (Value, error) {
		return f(r.tuple(v))
	})
}

func (r Relation) Where(p func(Value) (bool, error)) (_ Set, err error) {
	s, err := r.rows.Where(func(v Values) (bool, error) {
		return p(r.tuple(v))
	})
	if err != nil {
		return nil, err
	}
	return r.newBody(s), nil
}

// projectColumn extracts one attribute as a set of values, with no per-row
// Tuple (🎯T29.2). Missing attr returns ok=false so the caller can fall back.
func (r Relation) projectColumn(attr string) (Set, bool) {
	index, has := r.attrMap[attr]
	if !has {
		return nil, false
	}
	b := NewSetBuilder()
	for i := 0; i < r.rows.n; i++ {
		b.Add(r.rows.rowAt(i)[index])
	}
	s, err := b.Finish()
	if err != nil {
		return nil, false
	}
	return s, true
}

// projectCanonical builds an Array/String/Bytes/Dict from two source columns
// whose dest names are a reserved @-shape, without boxing source rows
// (🎯T29.3).
func (r Relation) projectCanonical(dst, src []string) (Set, bool) {
	if !isCanonicalTupleShape(dst) || len(dst) != len(src) {
		return nil, false
	}
	srcOf := make(map[string]string, 2)
	valName := ""
	for i, d := range dst {
		srcOf[d] = src[i]
		if d != "@" {
			valName = d
		}
	}
	atCol := r.getAttrIndex(srcOf["@"])
	valCol := r.getAttrIndex(srcOf[valName])
	if atCol < 0 || valCol < 0 {
		return nil, false
	}
	b := NewSetBuilder()
	for i := 0; i < r.rows.n; i++ {
		row := r.rows.rowAt(i)
		t, ok := canonicalTupleFromCols(valName, row[atCol], row[valCol])
		if !ok {
			return nil, false
		}
		b.Add(t)
	}
	s, err := b.Finish()
	if err != nil {
		return nil, false
	}
	return s, true
}

func canonicalTupleFromCols(valName string, at, val Value) (Value, bool) {
	switch valName {
	case DictValueAttr:
		return NewDictEntryTuple(at, val), true
	case ArrayItemAttr:
		n, ok := at.(Number)
		if !ok {
			return nil, false
		}
		i, ok := n.Int()
		if !ok {
			return nil, false
		}
		return NewArrayItemTuple(i, val), true
	case StringCharAttr:
		n, ok := at.(Number)
		if !ok {
			return nil, false
		}
		i, ok := n.Int()
		if !ok {
			return nil, false
		}
		c, ok := val.(Number)
		if !ok {
			return nil, false
		}
		return NewStringCharTuple(i, rune(c.Float64())), true
	case BytesByteAttr:
		n, ok := at.(Number)
		if !ok {
			return nil, false
		}
		i, ok := n.Int()
		if !ok {
			return nil, false
		}
		c, ok := val.(Number)
		if !ok {
			return nil, false
		}
		return NewBytesByteTuple(i, byte(c.Float64())), true
	default:
		return nil, false
	}
}

func (r Relation) projectDots(dst, src []string) (Set, bool) {
	if len(dst) != len(src) || len(src) == 0 {
		return nil, false
	}
	p := make(valueProjector, len(src))
	for i, name := range src {
		idx := r.getAttrIndex(name)
		if idx < 0 {
			return nil, false
		}
		p[i] = idx
	}
	out := make(NamesSlice, len(dst))
	copy(out, dst)
	if fastPaths && (r.rows.hasCandidateKey(p) || r.rows.groupBy(p).count() == r.rows.n) {
		// Injective: share the store so Count is the source cardinality
		// without copying rows (🎯T20).
		return newRelation(out, p, r.rows), true
	}
	rows := r.rows.Project(p)
	id := make(valueProjector, len(dst))
	for i := range id {
		id[i] = i
	}
	return newRelation(out, id, rows), true
}

func (r Relation) hasOnlyAttrs(names []string) bool {
	return r.attrSet == NewNames(names...)
}

func (r Relation) projector(names Names) (valueProjector, bool) {
	ns := names.OrderedNames()
	p := make(valueProjector, len(ns))
	for i, n := range ns {
		p[i] = r.getAttrIndex(n)
		if p[i] < 0 {
			return nil, false
		}
	}
	return p, true
}

// nestOnStore groups on the arena and keeps nested rows as a view (🎯T29.7).
func (r Relation) nestOnStore(nestAttrs Names, dest string, valuesOnly bool) (Set, bool) {
	if !nestAttrs.IsSubsetOf(r.attrSet) || !nestAttrs.IsTrue() {
		return nil, false
	}
	keyNames := r.attrSet.Minus(nestAttrs)
	if !keyNames.IsTrue() {
		return nil, false
	}
	keyProj, ok := r.projector(keyNames)
	if !ok {
		return nil, false
	}
	nestProj, ok := r.projector(nestAttrs)
	if !ok {
		return nil, false
	}
	g := r.rows.groupBy(keyProj)
	outNames := append(keyNames.OrderedNames(), dest)
	b := newStoreBuilder(len(outNames), g.count(), true)
	id := make(valueProjector, len(nestProj))
	for i := range id {
		id[i] = i
	}
	nestOut := NamesSlice(nestAttrs.OrderedNames())
	okAll := true
	g.each(func(bucket *groupBucket) {
		if !okAll {
			return
		}
		src := g.row(bucket.rep)
		row := make(Values, len(outNames))
		for i, idx := range keyProj {
			row[i] = src[idx]
		}
		view := r.rows.selView(bucket.rows)
		if valuesOnly {
			vr := newRelation(r.attrs, r.p, view)
			col, ok := vr.projectColumn(nestAttrs.OrderedNames()[0])
			if !ok {
				okAll = false
				return
			}
			row[len(outNames)-1] = col
		} else {
			var payload Value = newRelation(nestOut, id, view.Project(nestProj))
			if isCanonicalTupleShape(nestOut) {
				if v, ok := payload.(Relation).projectCanonical([]string(nestOut), []string(nestOut)); ok {
					payload = v
				}
			}
			row[len(outNames)-1] = payload
		}
		b.add(row)
	})
	if !okAll {
		return nil, false
	}
	return newRelation(outNames, makeIdentity(len(outNames)), b.finish()), true
}

func makeIdentity(n int) valueProjector {
	p := make(valueProjector, n)
	for i := range p {
		p[i] = i
	}
	return p
}

// orderByColumn argsorts rows by one arena column (🎯T29.6). Output tuples are
// inflated only when building the result array.
func (r Relation) orderByColumn(attr string) ([]Value, bool) {
	index, has := r.attrMap[attr]
	if !has {
		return nil, false
	}
	type kv struct {
		key Value
		id  uint32
	}
	pairs := make([]kv, r.rows.n)
	for i := 0; i < r.rows.n; i++ {
		pairs[i] = kv{key: r.rows.rowAt(i)[index], id: r.rows.arenaID(i)}
	}
	slices.SortStableFunc(pairs, func(a, b kv) int {
		if a.key.Less(b.key) {
			return -1
		}
		if b.key.Less(a.key) {
			return 1
		}
		return 0
	})
	values := make([]Value, len(pairs))
	width := r.rows.store.width
	for i, p := range pairs {
		values[i] = r.tuple(rowOf(r.rows.arena, width, int(p.id)))
	}
	return values, true
}

func (r Relation) getAttrIndex(attr string) int {
	for i, a := range r.attrs {
		if a == attr {
			return r.p[i]
		}
	}
	return -1
}

func (r Relation) CallAll(_ context.Context, v Value, sb SetBuilder) error {
	atIndex := r.getAttrIndex("@")
	if atIndex == -1 || len(r.attrs) != 2 {
		return errElementsNotMatchingAt
	}
	valIndex := 1
	if atIndex == 1 {
		valIndex = 0
	}

	for i := r.rows.Range(); i.Next(); {
		vals := i.Values()
		if vals[atIndex].Equal(v) {
			sb.Add(vals[valIndex])
		}
	}
	return nil
}

func (r Relation) unionSetSubsetBucket() string {
	// sort to ensure that identity of Relations are the same no matter the order of names.
	return r.attrs.GetSorted().String()
}

var relationKind = registerKind(211, reflect.TypeOf(Relation{}))

func (r Relation) Kind() int {
	return relationKind
}

func (r Relation) IsTrue() bool {
	return !r.rows.IsEmpty()
}

func (r Relation) Less(v Value) bool {
	if r.Kind() != v.Kind() {
		return r.Kind() < v.Kind()
	}
	r2 := v.(Relation)
	if r.attrs.LessNamesSlice(r2.attrs) && !r.attrs.EqualNamesSlice(r2.attrs) {
		return true
	}
	if r.Count() != r2.Count() {
		return r.Count() < r2.Count()
	}

	for i, j := r.ArrayEnumerator(), r2.ArrayEnumerator(); i.MoveNext() && j.MoveNext(); {
		left, right := i.Current(), j.Current()
		if left.Less(right) {
			return true
		}
		if right.Less(left) {
			return false
		}
	}
	return false
}

func (r Relation) Negate() Value {
	if !r.IsTrue() {
		return r
	}
	return NewTuple(NewAttr(negateTag, r))
}

// Join joins two relation based on the keys and defined outputs. Only does natural Join.
func (r Relation) Join(r2 Relation, keys, leftOutput, rightOutput NamesSlice) Set {
	if leftOutput.hasIntersect(rightOutput) {
		panic(fmt.Errorf("relation.Join: left and right output intersect, left: %v, right: %v", leftOutput, rightOutput))
	}
	leftKeysIndices, rightKeysIndices, leftOutputIndices, rightOutputIndices :=
		r.getIndices(keys), r2.getIndices(keys), r.getIndices(leftOutput), r2.getIndices(rightOutput)
	leftKey, rightKey, leftOutputProj, rightOutputProj :=
		r.p.compose(leftKeysIndices),
		r2.p.compose(rightKeysIndices),
		r.p.compose(leftOutputIndices),
		r2.p.compose(rightOutputIndices)
	count := len(leftOutput) + len(rightOutput)
	projection := make(valueProjector, 0, count)
	for i := 0; i < count; i++ {
		projection = append(projection, i)
	}
	rows := r.rows.Join(r2.rows, leftKey, rightKey, leftOutputProj, rightOutputProj)

	if rows.IsEmpty() {
		return False
	}
	if rows.IsLiteralTrue() {
		return True
	}
	attrs := append(leftOutput, rightOutput...)
	if len(attrs) == 2 {
		at, val := 0, 1
		if attrs[val] == "@" {
			at, val = val, at
		}
		if attrs[at] == "@" {
			switch attrs[val] {
			case ArrayItemAttr, BytesByteAttr, DictValueAttr, StringCharAttr:
				sb := NewSetBuilder()
				for i := rows.Range(); i.Next(); {
					values := i.Values().project(r.p)
					sb.Add(NewTuple(NewAttr("@", values.get(at)), NewAttr(attrs[val], values.get(val))))
				}
				set, err := sb.Finish()
				if err != nil {
					panic(err)
				}
				return set
			}
		}
	}

	return newRelation(attrs, projection, rows)
}

func (r Relation) Export(ctx context.Context) interface{} {
	if r.rows.IsEmpty() {
		return []interface{}{}
	}
	result := make([]interface{}, 0, r.rows.Count())
	for e := r.Enumerator(); e.MoveNext(); {
		result = append(result, e.Current().Export(ctx))
	}
	return result
}

type relationBuilder struct {
	prb     *positionalRelationBuilder
	mapping map[string]int
	names   NamesSlice
	attrSet Names // set when names are in sorted order, enabling zero-copy Add
}

func newRelationBuilder(names []string, cap int) *relationBuilder {
	m := make(map[string]int, len(names))
	for i, n := range names {
		m[n] = i
	}
	b := &relationBuilder{
		prb:     newPositionalRelationBuilder(len(names), cap),
		mapping: m,
		names:   names,
	}
	if slices.IsSorted(names) {
		b.attrSet = internNames(names)
	}
	return b
}

func (r *relationBuilder) Add(v Value) {
	if g, ok := v.(*GenericTuple); ok && fastPaths && r.attrSet.namesRep != nil && g.names == r.attrSet {
		// The tuple's values already are the row: both are immutable.
		r.prb.Add(Values(g.vals))
		return
	}
	t := v.(Tuple)
	values := make(Values, len(r.names))
	for name, index := range r.mapping {
		values[index] = t.MustGet(name)
	}
	r.prb.Add(values)
}

func (r *relationBuilder) Finish() (Set, error) {
	indices := make([]int, len(r.names))
	for i := range r.names {
		indices[i] = i
	}
	return newRelation(r.names, indices, r.prb.Finish()), nil
}

func (r Relation) getSetBuilder() setBuilder {
	return newGenericTypeSetBuilder()
}

func (r Relation) getBucket() fmt.Stringer {
	return genericType
}

func (r Relation) Eval(ctx context.Context, local Scope) (Value, error) {
	return r, nil
}

func (r Relation) Source() parser.Scanner {
	return *parser.NewScanner("")
}

func (r Relation) String() string {
	return fu.String(r)
}

func (r Relation) Format(f fmt.State, verb rune) {
	fu.WriteString(f, "{")

	attrs := r.attrs.GetSorted()
	fu.Fprintf(f, "|%s| ", strings.Join(attrs, ", "))
	projection := r.projectionBasedOnNames(attrs)
	notFirst := false
	for i := r.rows.OrderedRange(projection); i.Next(); {
		if notFirst {
			fu.WriteString(f, ", ")
		} else {
			notFirst = true
		}
		fu.Format(i.Values().project(projection), f, verb)
	}

	fu.WriteString(f, "}")
}

func (r Relation) projectionBasedOnNames(names NamesSlice) valueProjector {
	projection := make(valueProjector, 0, len(names))
	for _, n := range names {
		if i, has := r.attrMap[n]; has {
			projection = append(projection, i)
			continue
		}
		panic(fmt.Errorf("attribute %q does not exist in Relation %s", n, r))
	}
	return projection
}

func (r Relation) Equal(i Value) bool {
	if r2, is := i.(Relation); is {
		return r.EqualRelation(r2)
	}
	return false
}

func (r Relation) canonicalRelation() *positionalRelation {
	names := make(NamesSlice, len(r.attrs))
	copy(names, r.attrs)
	slices.Sort(names)
	projection := make(valueProjector, 0, len(r.attrs))
	for _, name := range names {
		projection = append(projection, r.attrMap[name])
	}
	return r.rows.Project(projection)
}

func (r Relation) EqualRelation(r2 Relation) bool {
	if r.rows.Count() != r2.rows.Count() || !r.attrs.EqualNamesSlice(r2.attrs) {
		return false
	}
	// Rows are positional; when both relations lay their attributes out
	// identically the row sets compare directly (frozen checks the sets'
	// hashes first). Only differing layouts need canonicalising.
	if fastPaths && r.sameLayout(r2) {
		return r.rows.EqualPositionalRelation(r2.rows)
	}
	return r.canonicalRelation().EqualPositionalRelation(r2.canonicalRelation())
}

// sameLayout reports whether r and r2 store each attribute at the same
// position.
func (r Relation) sameLayout(r2 Relation) bool {
	if r.rows.Width() != r2.rows.Width() {
		return false
	}
	if len(r.attrs) != len(r2.attrs) {
		return false
	}
	for i, name := range r.attrs {
		if r2.attrMap[name] != r.p[i] {
			return false
		}
	}
	return true
}

// Hash is the set wrap of finished row hashes: each row is hashed as a
// tuple (wrap of xor of name⋈value attrs) so pairings do not collapse,
// then those wraps are xored and wrapped as a set. Layout-independent:
// the positional Mix of Values would disagree with EqualRelation when
// two equal relations differ only in attribute order.
func (r Relation) Hash() uintptr {
	return hashSet(r.rows.shapeHash(r.attrSet, r.layout))
}

// RelationValuesEnumerator enumerates the values as Values.
type RelationValuesEnumerator struct {
	i *positionalRelationValuesEnumerator
	p valueProjector
}

func (e *RelationValuesEnumerator) Next() bool {
	return e.i.Next()
}

func (e *RelationValuesEnumerator) Values() Values {
	return e.i.Values().project(e.p).values()
}

func (r Relation) OrderedValuesEnumerator(names NamesSlice) *RelationValuesEnumerator {
	p := r.projectionBasedOnNames(names)
	return &RelationValuesEnumerator{
		i: r.rows.OrderedRange(p),
		p: p,
	}
}

type relationEnumerator struct {
	r Relation
	i *positionalRelationValuesEnumerator
}

func (r *relationEnumerator) MoveNext() bool {
	return r.i.Next()
}

func (r *relationEnumerator) Current() Value {
	return r.r.tuple(r.i.Values())
}
