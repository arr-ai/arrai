package rel

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/arr-ai/frozen"
	"github.com/arr-ai/wbnf/parser"
)

// GenericSet is a set of Values.
type GenericSet struct {
	set frozen.Set
}

// genericSet equivalents for Boolean true and false
var (
	None  = Set(EmptySet{})
	False = None
	True  = Set(TrueSet{})

	stringCharTupleType = reflect.TypeOf(StringCharTuple{})
	bytesByteTupleType  = reflect.TypeOf(BytesByteTuple{})
	arrayItemTupleType  = reflect.TypeOf(ArrayItemTuple{})
	dictEntryTupleType  = reflect.TypeOf(DictEntryTuple{})
	// TODO: uncomment when relation is implemented
	// genericTupleType    = reflect.TypeOf(&GenericTuple{})
)

func CanonicalSet(s Set) Set {
	if s, ok := s.(GenericSet); ok {
		b := NewSetBuilder()
		for e := s.Enumerator(); e.MoveNext(); {
			b.Add(e.Current())
		}
		result, err := b.Finish()
		if err != nil {
			panic(err)
		}
		return result
	}
	return s
}

func newGenericSetFromSet(s Set) Set {
	sb := frozen.SetBuilder{}
	for e := s.Enumerator(); e.MoveNext(); {
		sb.Add(e.Current())
	}
	return newSetFromFrozenSet(sb.Finish())
}

func newSetFromFrozenSet(s frozen.Set) Set {
	switch s.Count() {
	case 0:
		return None
	case 1:
		if s.Has(EmptyTuple) {
			return True
		}
	}
	return GenericSet{s}
}

// NewBool constructs a bool as a relation.
func NewBool(b bool) Set {
	if b {
		return True
	}
	return False
}

// Hash computes a hash for a genericSet.
func (s GenericSet) Hash(seed uintptr) uintptr {
	h := seed
	for e := s.Enumerator(); e.MoveNext(); {
		h ^= e.Current().Hash(0)
	}
	return h
}

// Equal tests two Sets for equality. Any other type returns false.
func (s GenericSet) Equal(v interface{}) bool {
	if t, ok := v.(GenericSet); ok {
		return s.set.Equal(t.set)
	}
	return false
}

// String returns a string representation of a genericSet.
func (s GenericSet) String() string {
	// {} == none
	if !s.IsTrue() {
		return "{}"
	}

	// {()} == true
	if s.Count() == 1 {
		e := s.Enumerator()
		e.MoveNext()
		if tuple, ok := e.Current().(Tuple); ok && !tuple.IsTrue() {
			return sTrue
		}
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	for i, value := range s.OrderedValues() {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(value.String())
	}
	buf.WriteString("}")
	return buf.String()
}

// Eval returns the set.
func (s GenericSet) Eval(ctx context.Context, local Scope) (Value, error) {
	return s, nil
}

// Source returns a scanner locating the GenericSet's source code.
func (s GenericSet) Source() parser.Scanner {
	return *parser.NewScanner("")
}

var genericSetKind = registerKind(200, reflect.TypeOf(Function{}))

// Kind returns a number that is unique for each major kind of Value.
func (s GenericSet) Kind() int {
	return genericSetKind
}

// Bool returns true iff the tuple has attributes.
func (s GenericSet) IsTrue() bool {
	return s.Count() > 0
}

// Less returns true iff v.Kind() < genericSet.Kind() or v is a
// genericSet and t precedes v in a lexicographical comparison of their
// sorted values.
func (s GenericSet) Less(v Value) bool {
	if s.Kind() != v.Kind() {
		return s.Kind() < v.Kind()
	}

	x := v.(GenericSet)
	a := s.OrderedValues()
	b := x.OrderedValues()
	n := len(a)
	if n > len(b) {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i].Less(b[i]) {
			return true
		}
		if b[i].Less(a[i]) {
			return false
		}
	}
	return len(a) < len(b)
}

// Negate returns {(negateTag): s}.
func (s GenericSet) Negate() Value {
	if !s.IsTrue() {
		return s
	}
	return NewTuple(NewAttr(negateTag, s))
}

// Export exports a genericSet as an array of exported Values.
func (s GenericSet) Export(ctx context.Context) interface{} {
	if s.set.IsEmpty() {
		return []interface{}{}
	}
	result := make([]interface{}, 0, s.set.Count())
	for e := s.Enumerator(); e.MoveNext(); {
		result = append(result, e.Current().Export(ctx))
	}
	return result
}

func (GenericSet) getSetBuilder() setBuilder {
	return newGenericTypeSetBuilder()
}

func (GenericSet) getBucket() fmt.Stringer {
	return genericType
}

// Count returns the number of elements in the genericSet.
func (s GenericSet) Count() int {
	return s.set.Count()
}

// Has returns true iff the given Value is in the genericSet.
func (s GenericSet) Has(value Value) bool {
	return s.set.Has(value)
}

// With returns the original genericSet with given value added. Iff the value was
// already present, the original genericSet is returned.
func (s GenericSet) With(v Value) Set {
	if v.getBucket() == genericType {
		return GenericSet{s.set.With(v)}
	}
	return toUnionSetWithItem(s, v)
}

// Without returns the original genericSet without the given value. Iff the value was
// already absent, the original genericSet is returned.
func (s GenericSet) Without(value Value) Set {
	return newSetFromFrozenSet(s.set.Without(value))
}

// Map maps values per f.
func (s GenericSet) Map(f func(v Value) (Value, error)) (Set, error) {
	sb := NewSetBuilder()
	for e := s.Enumerator(); e.MoveNext(); {
		v, err := f(e.Current())
		if err != nil {
			return nil, err
		}
		sb.Add(v)
	}
	return sb.Finish()
}

// Where returns a new genericSet with all the Values satisfying predicate p.
func (s GenericSet) Where(p func(v Value) (bool, error)) (_ Set, err error) {
	set := s.set.Where(func(elem interface{}) bool {
		if err != nil {
			return false
		}
		if elem == nil {
			return false
		}
		match, err2 := p(elem.(Value))
		if err2 != nil {
			err = err2
			return false
		}
		return match
	})
	if err != nil {
		return nil, err
	}
	return newSetFromFrozenSet(set), nil
}

var errElementsNotMatchingAt = fmt.Errorf("cannot call sets with elements not matching (@: _, _: _)")

func (s GenericSet) CallAll(_ context.Context, arg Value, b SetBuilder) error {
	// generic set does not hold tuples
	return errElementsNotMatchingAt
}

func (GenericSet) unionSetSubsetBucket() string {
	return genericType.String()
}

// Enumerator returns an enumerator over the Values in the genericSet.
func (s GenericSet) Enumerator() ValueEnumerator {
	return &genericSetEnumerator{s.set.Range()}
}

// Any return any value from the set.
func (s GenericSet) Any() Value {
	for e := s.Enumerator(); e.MoveNext(); {
		return e.Current()
	}
	panic("Any(): empty set")
}

type genericSetValueEnumerator struct {
	*genericSetEnumerator
}

func (g *genericSetValueEnumerator) Current() Value {
	// TODO: this is just a placeholder to replicate old functionality
	// this should return anything that's not @
	return g.i.Value().(Tuple).MustGet(ArrayItemAttr)
}

func (s GenericSet) ArrayEnumerator() ValueEnumerator {
	return &genericSetValueEnumerator{
		&genericSetEnumerator{
			s.set.OrderedRange(
				func(a, b interface{}) bool {
					return a.(Tuple).MustGet("@").(Number) < b.(Tuple).MustGet("@").(Number)
				},
			),
		},
	}
}

// genericSetEnumerator represents an enumerator over a genericSet.
type genericSetEnumerator struct {
	i frozen.Iterator
}

// MoveNext moves the enumerator to the next Value.
func (e *genericSetEnumerator) MoveNext() bool {
	return e.i.Next()
}

// Current returns the enumerator's current Value.
func (e *genericSetEnumerator) Current() Value {
	return e.i.Value().(Value)
}

// ValueList represents a []Value for use in sort.Sort().
type ValueList []Value

func (vl ValueList) Len() int {
	return len(vl)
}

func (vl ValueList) Less(i, j int) bool {
	return vl[i].Less(vl[j])
}

func (vl ValueList) Swap(i, j int) {
	vl[i], vl[j] = vl[j], vl[i]
}

// OrderedValues returns Values in a canonical ordering.
func (s GenericSet) OrderedValues() []Value {
	a := make([]Value, s.Count())
	i := 0
	for e := s.Enumerator(); e.MoveNext(); {
		a[i] = e.Current()
		i++
	}
	sort.Sort(ValueList(a))
	return a
}
