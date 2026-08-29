package rel

import (
	"github.com/arr-ai/hash/hash128"

	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/arr-ai/frozen"
	"github.com/arr-ai/wbnf/parser"

	"github.com/arr-ai/arrai/pkg/fu"
)

// Relation is a Set that only contains Tuples, all of which map the same keys.
type Relation struct {
	attrs   NamesSlice
	p       valueProjector
	rows    *positionalRelation // TODO: experiment with column table
	attrMap map[string]int      // cached mapIndices(attrs, p)

	// shape is the tuple shape of every row; layout[i] is the row position of
	// shape.names[i]. When layout is the identity a row's values already are
	// the tuple's values, so inflating a row is a wrap, not a copy.
	shape  *Shape
	layout []int
	direct bool
}

func newRelation(attrs NamesSlice, p valueProjector, rows *positionalRelation) Relation {
	r := Relation{attrs: attrs, p: p, rows: rows, attrMap: mapIndices(attrs, p)}
	r.shape = shapeOf(attrs.GetSorted())
	r.layout = make([]int, len(r.shape.names))
	r.direct = true
	for i, name := range r.shape.names {
		r.layout[i] = r.attrMap[name]
		if r.layout[i] != i {
			r.direct = false
		}
	}
	return r
}

// tuple inflates a row to a Tuple.
func (r Relation) tuple(row Values) Tuple {
	if r.direct && len(row) == len(r.layout) {
		return newShapedTuple(r.shape, row)
	}
	vals := make([]Value, len(r.layout))
	for i, j := range r.layout {
		vals[i] = row[j]
	}
	return newShapedTuple(r.shape, vals)
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

func (r Relation) tupleToValues(t Tuple) Values {
	if len(r.attrs) != t.Count() {
		panic("tupleToValues: names and values don't have the same number")
	}
	values := make(Values, len(r.attrs))
	for i, name := range r.attrs {
		values[r.p[i]] = t.MustGet(name)
	}
	return values
}

func (r Relation) Has(v Value) bool {
	if t, is := v.(Tuple); is {
		return r.attrs.EqualTupleAttrs(t) && r.rows.Has(r.tupleToValues(t))
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

func (r Relation) With(v Value) Set {
	if t, is := v.(Tuple); is && r.attrs.EqualTupleAttrs(t) {
		return newRelation(r.attrs, r.p, r.rows.With(r.tupleToValues(t)))
	}
	return toUnionSetWithItem(r, v)
}

func (r Relation) Without(v Value) Set {
	if t, is := v.(Tuple); is && r.attrs.EqualTupleAttrs(t) {
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
	shape   *Shape // set when names are in shape order, enabling zero-copy Add
}

func newRelationBuilder(names []string, cap int) *relationBuilder {
	m := make(map[string]int, len(names))
	for i, n := range names {
		m[n] = i
	}
	b := &relationBuilder{
		prb:     &positionalRelationBuilder{sb: frozen.NewSetBuilder[any](cap)},
		mapping: m,
		names:   names,
	}
	if sort.StringsAreSorted(names) {
		b.shape = shapeOf(names)
	}
	return b
}

func (r *relationBuilder) Add(v Value) {
	if g, ok := v.(*GenericTuple); ok && r.shape != nil && g.shape == r.shape {
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
	sort.Strings(names)
	projection := make(valueProjector, 0, len(r.attrs))
	for _, name := range names {
		projection = append(projection, r.attrMap[name])
	}
	isContiguous := projection.isContiguous()
	return &positionalRelation{
		set: frozen.SetMap(r.rows.set, func(elem any) any {
			if isContiguous {
				return elem.(Values)[projection[0] : projection[len(projection)-1]+1]
			}
			return elem.(Values).project(projection).values()
		}),
	}
}

func (r Relation) EqualRelation(r2 Relation) bool {
	if r.rows.Count() != r2.rows.Count() || !r.attrs.EqualNamesSlice(r2.attrs) {
		return false
	}
	// Rows are positional; when both relations lay their attributes out
	// identically the row sets compare directly (frozen checks the sets'
	// hashes first). Only differing layouts need canonicalising.
	if r.sameLayout(r2) {
		return r.rows.set.Equal(r2.rows.set)
	}
	return r.canonicalRelation().set.Equal(r2.canonicalRelation().set)
}

// sameLayout reports whether r and r2 store each attribute at the same
// position.
func (r Relation) sameLayout(r2 Relation) bool {
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

func (r Relation) Hash(seed uintptr) uintptr {
	return r.Hash128().Seeded(seed)
}

// Hash128 computes the 128-bit hash of a Relation: the xor over attribute
// names, xor'd with the xor over rows of each row's own name/value hash
// (the same per-attribute formula GenericTuple uses). Row hashes must be
// computed this way, and not taken from the positional row set's own hash,
// because that mixes values in strict positional order — which would make
// Hash128 disagree with EqualRelation whenever two equal relations differ
// only in their internal attribute ordering (e.g. from different
// join/projection code paths), violating the hash/equals contract.
func (r Relation) Hash128() hash128.H128 {
	// The attribute-name half is a property of the shape. Rows hash per
	// attribute through the shape's cached name hashes, in shape order via
	// layout, so Hash128 agrees with EqualRelation regardless of the
	// relation's internal attribute ordering.
	h := relationSalt.Xor(r.shape.namesH)
	for e := r.rows.Range(); e.Next(); {
		row := e.Values()
		var rh hash128.H128
		for i, j := range r.layout {
			rh = rh.Xor(hashAttr(r.shape.nameH[i], row[j]))
		}
		h = h.Xor(rh)
	}
	return h
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
