package rel

import (
	"github.com/arr-ai/hash/hash128"

	"context"
	"fmt"
	"reflect"

	"github.com/arr-ai/frozen"
	"github.com/arr-ai/wbnf/parser"

	"github.com/arr-ai/arrai/pkg/fu"
)

// GenericSet is a set of Values.
type GenericSet struct {
	set frozen.Set[Value]

	// canonical records that CanonicalSet has already verified this set
	// stays a GenericSet (every element in the generic bucket), so
	// re-canonicalising it is a no-op. A union of two canonical sets is
	// canonical — their elements' buckets don't change — which lets
	// accumulator folds (acc | {x}) skip an O(n) rebuild per step.
	canonical bool
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
		if s.canonical {
			return s
		}
		b := NewSetBuilder()
		for e := s.Enumerator(); e.MoveNext(); {
			b.Add(e.Current())
		}
		result, err := b.Finish()
		if err != nil {
			panic(err)
		}
		if g, ok := result.(GenericSet); ok {
			g.canonical = true
			return g
		}
		return result
	}
	return s
}

func newGenericSetFromSet(s Set) Set {
	sb := frozen.SetBuilder[Value]{}
	for e := s.Enumerator(); e.MoveNext(); {
		sb.Add(e.Current())
	}
	return newSetFromFrozenSet(sb.Finish())
}

func newSetFromFrozenSet(s frozen.Set[Value]) Set {
	switch s.Count() {
	case 0:
		return None
	case 1:
		if s.Has(EmptyTuple) {
			return True
		}
	}
	return GenericSet{set: s}
}

// NewBool constructs a bool as a relation.
func NewBool(b bool) Set {
	if b {
		return True
	}
	return False
}

// Hash computes a hash for a genericSet.
// Delegates to frozen.Set, which uses the tree's H0 for seeds 0/1 instead of
// re-walking every element (important when GenericSets nest inside other sets).
func (s GenericSet) Hash(seed uintptr) uintptr {
	return s.set.Hash(seed)
}

// Hash128 returns the 128-bit hash maintained by the underlying frozen set.
func (s GenericSet) Hash128() hash128.H128 {
	return s.set.Hash128()
}

// Equal tests two Sets for equality. Any other type returns false.
func (s GenericSet) Equal(v Value) bool {
	if t, ok := v.(GenericSet); ok {
		return s.set.Equal(t.set)
	}
	return false
}

// String returns a string representation of a genericSet.
func (s GenericSet) String() string {
	return fu.String(s)
}

// Format formats a genericSet.
func (s GenericSet) Format(f fmt.State, verb rune) {
	// {} == none
	if !s.IsTrue() {
		if verb == 'v' {
			fu.WriteString(f, sEmptySet)
		}
		return
	}

	// {()} == true
	if s.Count() == 1 {
		e := s.Enumerator()
		e.MoveNext()
		if tuple, ok := e.Current().(Tuple); ok && !tuple.IsTrue() {
			fu.WriteString(f, sTrue)
			return
		}
	}

	reprOrderableSet(f, s)
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
	a, b := s.OrderedValues(), x.OrderedValues()
	for {
		if !a.MoveNext() {
			return b.MoveNext()
		}
		if !b.MoveNext() {
			return false
		}
		ai := a.Current()
		bi := b.Current()
		if ai.Less(bi) {
			return true
		}
		if bi.Less(ai) {
			return false
		}
	}
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
		// Adding a generic-bucket value cannot change any element's bucket,
		// so canonicality is preserved.
		return GenericSet{set: s.set.With(v), canonical: s.canonical}
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
func (s GenericSet) Where(p func(v Value) (bool, error)) (Set, error) {
	var failure whereErr
	set := s.set.Where(func(elem Value) bool {
		if failure.failed() {
			return false
		}
		match, err := p(elem)
		if err != nil {
			failure.set(err)
			return false
		}
		return match
	})
	if err := failure.get(); err != nil {
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
				func(a, b Value) bool {
					return a.(Tuple).MustGet("@").(Number) < b.(Tuple).MustGet("@").(Number)
				},
			),
		},
	}
}

// genericSetEnumerator represents an enumerator over a genericSet.
type genericSetEnumerator struct {
	i frozen.Iterator[Value]
}

// MoveNext moves the enumerator to the next Value.
func (e *genericSetEnumerator) MoveNext() bool {
	return e.i.Next()
}

// Current returns the enumerator's current Value.
func (e *genericSetEnumerator) Current() Value {
	return e.i.Value()
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
func (s GenericSet) OrderedValues() ValueEnumerator {
	return OrderedValueEnumerator(s.Enumerator(), ValueLess)
}
