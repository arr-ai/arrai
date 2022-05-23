package rel

import (
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
	attrs NamesSlice
	p     valueProjector
	rows  *positionalRelation // TODO: experiment with column table
}

func newRelation(attrs NamesSlice, p valueProjector, rows *positionalRelation) Relation {
	return Relation{attrs, p, rows}
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

func valuesToTuple(val Values, names map[string]int) Tuple {
	m := make([]Attr, 0, len(val))
	for name, index := range names {
		m = append(m, NewAttr(name, val[index]))
	}
	return NewTuple(m...)
}

func (r Relation) Enumerator() ValueEnumerator {
	return &relationEnumerator{
		attrs: mapIndices(r.attrs, r.p),
		i:     r.rows.Range(),
	}
}

func (r Relation) OrderedValues() ValueEnumerator {
	return OrderedValueEnumerator(r.Enumerator(), ValueLess)
}

func (r Relation) ArrayEnumerator() ValueEnumerator {
	return &relationEnumerator{
		attrs: mapIndices(r.attrs, r.p),
		i:     r.rows.OrderedRange(r.p),
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
	m := mapIndices(r.attrs, r.p)
	return r.rows.Map(func(v Values) (Value, error) {
		return f(valuesToTuple(v, m))
	})
}

func (r Relation) Where(p func(Value) (bool, error)) (_ Set, err error) {
	m := mapIndices(r.attrs, r.p)
	s, err := r.rows.Where(func(v Values) (bool, error) {
		return p(valuesToTuple(v, m))
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
}

func newRelationBuilder(names []string, cap int) *relationBuilder {
	m := make(map[string]int, len(names))
	for i, n := range names {
		m[n] = i
	}
	return &relationBuilder{
		prb:     &positionalRelationBuilder{sb: frozen.NewSetBuilder(cap)},
		mapping: m,
		names:   names,
	}
}

func (r *relationBuilder) Add(v Value) {
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
	indices := mapIndices(r.attrs, r.p)
	for _, n := range names {
		if i, has := indices[n]; has {
			projection = append(projection, i)
			continue
		}
		panic(fmt.Errorf("attribute %q does not exist in Relation %s", n, r))
	}
	return projection
}

func (r Relation) Equal(i interface{}) bool {
	if r2, is := i.(Relation); is {
		return r.EqualRelation(r2)
	}
	return false
}

func (r Relation) canonicalRelation() *positionalRelation {
	names := make(NamesSlice, len(r.attrs))
	copy(names, r.attrs)
	sort.Strings(names)
	m := mapIndices(r.attrs, r.p)
	projection := make(valueProjector, 0, len(r.attrs))
	for _, name := range names {
		projection = append(projection, m[name])
	}
	isContiguous := projection.isContiguous()
	return &positionalRelation{
		set: r.rows.set.Map(func(elem interface{}) interface{} {
			if isContiguous {
				return elem.(Values)[projection[0] : projection[len(projection)-1]+1]
			}
			return elem.(Values).project(projection).values()
		}),
	}
}

func (r Relation) EqualRelation(r2 Relation) bool {
	if !r.attrs.EqualNamesSlice(r2.attrs) {
		return false
	}
	return r.canonicalRelation().set.EqualSet(r2.canonicalRelation().set)
}

func (r Relation) Hash(seed uintptr) uintptr {
	var h uintptr
	for i := r.Enumerator(); i.MoveNext(); {
		h ^= i.Current().Hash(seed)
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
	attrs map[string]int
	i     *positionalRelationValuesEnumerator
}

func (r *relationEnumerator) MoveNext() bool {
	return r.i.Next()
}

func (r *relationEnumerator) Current() Value {
	return valuesToTuple(r.i.Values(), r.attrs)
}
