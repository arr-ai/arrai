package rel

import (
	"context"
	"fmt"
	"reflect"

	"github.com/arr-ai/wbnf/parser"
)

const DictValueAttr = "@value"

// DictEntryTuple represents a tuple of the form (@: at, @value: item).
type DictEntryTuple struct {
	at    Value
	value Value
}

// NewDictEntryTuple constructs a CharTuple.
func NewDictEntryTuple(at Value, item Value) DictEntryTuple {
	return DictEntryTuple{at: at, value: item}
}

func newDictEntryTupleFromTuple(t Tuple) (DictEntryTuple, bool) {
	var at Value
	var item Value
	m := NewTupleMatcher(
		map[string]Matcher{
			"@":           Let(func(i Value) { at = i }),
			DictValueAttr: Let(func(v Value) { item = v }),
		},
		Lit(EmptyTuple),
	)
	if m.Match(t) {
		return NewDictEntryTuple(at, item), true
	}
	return DictEntryTuple{}, false
}

func maybeNewDictEntryTupleFromTuple(t Tuple) Tuple {
	if t, ok := newDictEntryTupleFromTuple(t); ok {
		return t
	}
	return t
}

func (t DictEntryTuple) asGenericTuple() Tuple {
	return newTuple(
		NewAttr("@", t.at),
		NewAttr(DictValueAttr, t.value),
	)
}

// Hash computes a hash for a CharTuple.
func (t DictEntryTuple) Hash(seed uintptr) uintptr {
	return t.value.Hash(t.at.Hash(seed))
}

// Equal tests two Tuples for equality. Any other type returns false.
func (t DictEntryTuple) Equal(v interface{}) bool {
	if u, ok := v.(DictEntryTuple); ok {
		return t.at.Equal(u.at) && t.value.Equal(u.value)
	}
	return false
}

// String returns a string representation of a Tuple.
func (t DictEntryTuple) String() string {
	return fmt.Sprintf("(@: %v, %s: %v)", t.at, DictValueAttr, t.value)
}

// Eval returns the tuple.
func (t DictEntryTuple) Eval(ctx context.Context, local Scope) (Value, error) {
	return t, nil
}

// Source returns a scanner locating the DictEntryTuple's source code.
func (t DictEntryTuple) Source() parser.Scanner {
	return *parser.NewScanner("")
}

var dictValueTupleKind = registerKind(303, reflect.TypeOf((*DictEntryTuple)(nil)))

// Kind returns a number that is unique for each major kind of Value.
func (t DictEntryTuple) Kind() int {
	return dictValueTupleKind
}

// Bool returns true iff the tuple has attributes.
func (t DictEntryTuple) IsTrue() bool {
	return true
}

// Less returns true iff v is not a number or tuple, or v is a tuple and t
// precedes v in a lexicographical comparison of their name/value pairs.
func (t DictEntryTuple) Less(v Value) bool {
	if t.Kind() != v.Kind() {
		return t.Kind() < v.Kind()
	}
	u := v.(DictEntryTuple)
	if !t.at.Equal(u.at) {
		return t.at.Less(u.at)
	}
	return t.value.Less(u.value)
}

func (t DictEntryTuple) Negate() Value {
	return DictEntryTuple{at: t.at.Negate(), value: t.value.Negate()}
}

// Export exports a Tuple.
func (t DictEntryTuple) Export(ctx context.Context) interface{} {
	return map[string]interface{}{
		"@":           t.at.Export(ctx),
		DictValueAttr: t.value.Export(ctx),
	}
}

// Count returns how many attributes are in the Tuple.
func (t DictEntryTuple) Count() int {
	return 2
}

// Get returns the Value associated with a name, and true iff it was found.
func (t DictEntryTuple) Get(name string) (Value, bool) {
	switch name {
	case "@":
		return t.at, true
	case DictValueAttr:
		return t.value, true
	}
	return nil, false
}

// MustGet returns e.Get(name) or panics if an error occurs.
func (t DictEntryTuple) MustGet(name string) Value {
	if v, has := t.Get(name); has {
		return v
	}
	panic(fmt.Errorf("%q not found", name))
}

// With returns a Tuple with all name/Value pairs in t (except the one for the
// given name, if present) with the addition of the given name/Value pair.
func (t DictEntryTuple) With(name string, value Value) Tuple {
	return maybeNewDictEntryTupleFromTuple(t.asGenericTuple().With(name, value))
}

// Without returns a Tuple with all name/Value pairs in t exception the one of
// the given name.
func (t DictEntryTuple) Without(name string) Tuple {
	if name == "@" || name == DictValueAttr {
		return t.asGenericTuple().Without(name)
	}
	return t
}

func (t DictEntryTuple) Map(f func(Value) (Value, error)) (Tuple, error) {
	at, err := f(t.at)
	if err != nil {
		return nil, err
	}
	value, err := f(t.value)
	if err != nil {
		return nil, err
	}

	return NewDictEntryTuple(at, value), nil
}

// HasName returns true iff the Tuple has an attribute with the given name.
func (t DictEntryTuple) HasName(name string) bool {
	return name == "@" || name == DictValueAttr
}

// Names returns the attribute names.
func (t DictEntryTuple) Names() Names {
	return NewNames("@", DictValueAttr)
}

// Project returns a tuple with the given names from this tuple, or nil if any
// name wasn't found.
func (t DictEntryTuple) Project(names Names) Tuple {
	if names.Has("@") && names.Has(DictValueAttr) {
		return t
	}
	return t.asGenericTuple().Project(names)
}

// Enumerator returns an enumerator over the Values in the CharTuple.
func (t DictEntryTuple) Enumerator() AttrEnumerator {
	return &dictEntryTupleEnumerator{t: t, i: -1}
}

type dictEntryTupleEnumerator struct {
	t     DictEntryTuple
	i     int
	name  string
	value Value
}

func (e *dictEntryTupleEnumerator) MoveNext() bool {
	if e.i == 1 {
		return false
	}
	e.i++
	switch e.i {
	case 0:
		e.name = "@"
		e.value = e.t.at
	case 1:
		e.name = DictValueAttr
		e.value = e.t.value
	}
	return true
}

func (e *dictEntryTupleEnumerator) Current() (string, Value) {
	return e.name, e.value
}
