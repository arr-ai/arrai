package rel

import (
	"context"
	"fmt"
	"reflect"

	"github.com/arr-ai/hash"
	"github.com/arr-ai/wbnf/parser"
)

const ArrayItemAttr = "@item"

// ArrayItemTuple represents a tuple of the form (@: at, @item: item).
type ArrayItemTuple struct {
	at   int
	item Value
}

// NewArrayItemTuple constructs a CharTuple.
func NewArrayItemTuple(at int, item Value) ArrayItemTuple {
	return ArrayItemTuple{at: at, item: item}
}

func newArrayItemTupleFromTuple(t Tuple) (ArrayItemTuple, bool) {
	var at int
	var item Value
	m := NewTupleMatcher(
		map[string]Matcher{
			"@":           MatchInt(func(i int) { at = i }),
			ArrayItemAttr: Let(func(v Value) { item = v }),
		},
		Lit(EmptyTuple),
	)
	if m.Match(t) {
		return NewArrayItemTuple(at, item), true
	}
	return ArrayItemTuple{}, false
}

func maybeNewArrayItemTupleFromTuple(t Tuple) Tuple {
	if t, ok := newArrayItemTupleFromTuple(t); ok {
		return t
	}
	return t
}

func (t ArrayItemTuple) asGenericTuple() Tuple {
	return newTuple(
		NewIntAttr("@", t.at),
		NewAttr(ArrayItemAttr, t.item),
	)
}

// Hash computes a hash for a CharTuple.
func (t ArrayItemTuple) Hash(seed uintptr) uintptr {
	h := hash.Int(t.at, seed)
	return t.item.Hash(h)
}

// Equal tests two Tuples for equality. Any other type returns false.
func (t ArrayItemTuple) Equal(v interface{}) bool {
	if u, ok := v.(ArrayItemTuple); ok {
		return t.at == u.at && t.item.Equal(u.item)
	}
	return false
}

// String returns a string representation of a Tuple.
func (t ArrayItemTuple) String() string {
	return fmt.Sprintf("(@: %d, %s: %v)", t.at, ArrayItemAttr, t.item)
}

// Eval returns the tuple.
func (t ArrayItemTuple) Eval(ctx context.Context, local Scope) (Value, error) {
	return t, nil
}

// Source returns a scanner locating the ArrayItemTuple's source code.
func (t ArrayItemTuple) Source() parser.Scanner {
	return *parser.NewScanner("")
}

var arrayItemTupleKind = registerKind(302, reflect.TypeOf((*ArrayItemTuple)(nil)))

// Kind returns a number that is unique for each major kind of Value.
func (t ArrayItemTuple) Kind() int {
	return arrayItemTupleKind
}

// Bool returns true iff the tuple has attributes.
func (t ArrayItemTuple) IsTrue() bool {
	return true
}

// Less returns true iff v is not a number or tuple, or v is a tuple and t
// precedes v in a lexicographical comparison of their name/value pairs.
func (t ArrayItemTuple) Less(v Value) bool {
	if t.Kind() != v.Kind() {
		return t.Kind() < v.Kind()
	}
	u := v.(ArrayItemTuple)
	if t.at != u.at {
		return t.at < u.at
	}
	return t.item.Less(u.item)
}

func (t ArrayItemTuple) Negate() Value {
	return ArrayItemTuple{at: -t.at, item: t.item.Negate()}
}

// Export exports a Tuple.
func (t ArrayItemTuple) Export(ctx context.Context) interface{} {
	return map[string]interface{}{
		"@":           t.at,
		ArrayItemAttr: t.item.Export(ctx),
	}
}

// Count returns how many attributes are in the Tuple.
func (t ArrayItemTuple) Count() int {
	return 2
}

// Get returns the Value associated with a name, and true iff it was found.
func (t ArrayItemTuple) Get(name string) (Value, bool) {
	switch name {
	case "@":
		return NewNumber(float64(t.at)), true
	case ArrayItemAttr:
		return t.item, true
	}
	return nil, false
}

// MustGet returns e.Get(name) or panics if an error occurs.
func (t ArrayItemTuple) MustGet(name string) Value {
	if v, has := t.Get(name); has {
		return v
	}
	panic(fmt.Errorf("%q not found", name))
}

// With returns a Tuple with all name/Value pairs in t (except the one for the
// given name, if present) with the addition of the given name/Value pair.
func (t ArrayItemTuple) With(name string, value Value) Tuple {
	return maybeNewArrayItemTupleFromTuple(t.asGenericTuple().With(name, value))
}

// Without returns a Tuple with all name/Value pairs in t exception the one of
// the given name.
func (t ArrayItemTuple) Without(name string) Tuple {
	if name == "@" || name == ArrayItemAttr {
		return t.asGenericTuple().Without(name)
	}
	return t
}

func (t ArrayItemTuple) Map(f func(Value) (Value, error)) (Tuple, error) {
	at, err := f(NewNumber(float64(t.at)))
	if err != nil {
		return nil, err
	}
	item, err := f(t.item)
	if err != nil {
		return nil, err
	}
	if at, ok := at.(Number); ok {
		if at, is := at.Int(); is {
			return NewArrayItemTuple(at, item), nil
		}
	}
	return t.asGenericTuple().Map(f)
}

// HasName returns true iff the Tuple has an attribute with the given name.
func (t ArrayItemTuple) HasName(name string) bool {
	return name == "@" || name == ArrayItemAttr
}

// Names returns the attribute names.
func (t ArrayItemTuple) Names() Names {
	return NewNames("@", ArrayItemAttr)
}

// Project returns a tuple with the given names from this tuple, or nil if any
// name wasn't found.
func (t ArrayItemTuple) Project(names Names) Tuple {
	if names.Has("@") && names.Has(ArrayItemAttr) {
		return t
	}
	return t.asGenericTuple().Project(names)
}

// Enumerator returns an enumerator over the Values in the CharTuple.
func (t ArrayItemTuple) Enumerator() AttrEnumerator {
	return &arrayItemTupleEnumerator{t: t, i: -1}
}

type arrayItemTupleEnumerator struct {
	t     ArrayItemTuple
	i     int
	name  string
	value Value
}

func (e *arrayItemTupleEnumerator) MoveNext() bool {
	if e.i == 1 {
		return false
	}
	e.i++
	switch e.i {
	case 0:
		e.name = "@"
		e.value = NewNumber(float64(e.t.at))
	case 1:
		e.name = ArrayItemAttr
		e.value = e.t.item
	}
	return true
}

func (e *arrayItemTupleEnumerator) Current() (string, Value) {
	return e.name, e.value
}
