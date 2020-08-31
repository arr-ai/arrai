package rel //nolint:dupl

import (
	"context"
	"fmt"
	"reflect"

	"github.com/arr-ai/hash"
	"github.com/arr-ai/wbnf/parser"
)

// BytesByteTuple represents a tuple of the form (@: at, @byte: byteval).
type BytesByteTuple struct {
	at      int
	byteval byte
}

// NewBytesByteTuple constructs a BytesByteTuple.
func NewBytesByteTuple(at int, byteval byte) BytesByteTuple {
	return BytesByteTuple{at: at, byteval: byteval}
}

func newBytesByteTupleFromTuple(t Tuple) (BytesByteTuple, bool) {
	var at int
	var byteval byte
	m := NewTupleMatcher(
		map[string]Matcher{
			"@":           MatchInt(func(i int) { at = i }),
			BytesByteAttr: MatchInt(func(i int) { byteval = byte(i) }),
		},
		Lit(EmptyTuple),
	)
	if m.Match(t) {
		return NewBytesByteTuple(at, byteval), true
	}
	return BytesByteTuple{}, false
}

func maybeNewBytesByteTupleFromTuple(t Tuple) Tuple {
	if t, ok := newBytesByteTupleFromTuple(t); ok {
		return t
	}
	return t
}

func (t BytesByteTuple) asGenericTuple() Tuple {
	return newTuple(
		NewIntAttr("@", t.at),
		NewIntAttr(BytesByteAttr, int(t.byteval)),
	)
}

// Hash computes a hash for a BytesByteTuple.
func (t BytesByteTuple) Hash(seed uintptr) uintptr {
	h := hash.Int(t.at, seed)
	return hash.Uint8(t.byteval, h)
}

// Equal tests two Tuples for equality. Any other type returns false.
func (t BytesByteTuple) Equal(v interface{}) bool {
	if u, ok := v.(BytesByteTuple); ok {
		return t == u
	}
	return false
}

// String returns a string representation of a Tuple.
func (t BytesByteTuple) String() string {
	return fmt.Sprintf("(@: %d, %s: %d)", t.at, BytesByteAttr, t.byteval)
}

// Eval returns the tuple.
func (t BytesByteTuple) Eval(ctx context.Context, local Scope) (Value, error) {
	return t, nil
}

// Source returns a scanner locating the BytesByteTuple's source code.
func (t BytesByteTuple) Source() parser.Scanner {
	return *parser.NewScanner("")
}

var bytesByteTupleKind = registerKind(304, reflect.TypeOf((*BytesByteTuple)(nil)))

// Kind returns a number that is unique for each major kind of Value.
func (t BytesByteTuple) Kind() int {
	return bytesByteTupleKind
}

// Bool returns true iff the tuple has attributes.
func (t BytesByteTuple) IsTrue() bool {
	return true
}

// Less returns true iff v is not a number or tuple, or v is a tuple and t
// precedes v in a lexicographical comparison of their name/value pairs.
func (t BytesByteTuple) Less(v Value) bool {
	if t.Kind() != v.Kind() {
		return t.Kind() < v.Kind()
	}
	u := v.(BytesByteTuple)
	if t.at != u.at {
		return t.at < u.at
	}
	return t.byteval < u.byteval
}

func (t BytesByteTuple) Negate() Value {
	return NewTuple(
		NewAttr("@", NewNumber(-float64(t.at))),
		NewAttr("@byte", NewNumber(-float64(t.byteval))),
	)
}

// Export exports a Tuple.
func (t BytesByteTuple) Export(_ context.Context) interface{} {
	return map[string]interface{}{
		"@":           t.at,
		BytesByteAttr: t.byteval,
	}
}

// Count returns how many attributes are in the Tuple.
func (t BytesByteTuple) Count() int {
	return 2
}

// Get returns the Value associated with a name, and true iff it was found.
func (t BytesByteTuple) Get(name string) (Value, bool) {
	switch name {
	case "@":
		return NewNumber(float64(t.at)), true
	case BytesByteAttr:
		return NewNumber(float64(t.byteval)), true
	}
	return nil, false
}

// MustGet returns e.Get(name) or panics if an error occurs.
func (t BytesByteTuple) MustGet(name string) Value {
	if v, has := t.Get(name); has {
		return v
	}
	panic(fmt.Errorf("%q not found", name))
}

// With returns a Tuple with all name/Value pairs in t (except the one for the
// given name, if present) with the addition of the given name/Value pair.
func (t BytesByteTuple) With(name string, value Value) Tuple {
	return maybeNewBytesByteTupleFromTuple(t.asGenericTuple().With(name, value))
}

// Without returns a Tuple with all name/Value pairs in t exception the one of
// the given name.
func (t BytesByteTuple) Without(name string) Tuple {
	if name == "@" || name == BytesByteAttr {
		return t.asGenericTuple().Without(name)
	}
	return t
}

func (t BytesByteTuple) Map(f func(Value) (Value, error)) (Tuple, error) { //nolint:dupl
	at, err := f(NewNumber(float64(t.at)))
	if err != nil {
		return nil, err
	}
	byteval, err := f(NewNumber(float64(t.byteval)))
	if err != nil {
		return nil, err
	}
	if at, ok := at.(Number); ok {
		if at, is := at.Int(); is {
			if byteval, ok := byteval.(Number); ok {
				if byteval, is := byteval.Int(); is {
					return NewBytesByteTuple(at, byte(byteval)), nil
				}
			}
		}
	}
	return t.asGenericTuple().Map(f)
}

// HasName returns true iff the Tuple has an attribute with the given name.
func (t BytesByteTuple) HasName(name string) bool {
	return name == "@" || name == BytesByteAttr
}

// Names returns the attribute names.
func (t BytesByteTuple) Names() Names {
	return NewNames("@", BytesByteAttr)
}

// Project returns a tuple with the given names from this tuple, or nil if any
// name wasn't found.
func (t BytesByteTuple) Project(names Names) Tuple {
	if names.Has("@") && names.Has(BytesByteAttr) {
		return t
	}
	return t.asGenericTuple().Project(names)
}

// Enumerator returns an enumerator over the Values in the BytesByteTuple.
func (t BytesByteTuple) Enumerator() AttrEnumerator {
	return &bytesByteTupleEnumerator{t: t, i: -1}
}

type bytesByteTupleEnumerator struct {
	t     BytesByteTuple
	i     int
	name  string
	value Value
}

func (e *bytesByteTupleEnumerator) MoveNext() bool {
	if e.i == 1 {
		return false
	}
	e.i++
	switch e.i {
	case 0:
		e.name = "@"
		e.value = NewNumber(float64(e.t.at))
	case 1:
		e.name = BytesByteAttr
		e.value = NewNumber(float64(e.t.byteval))
	}
	return true
}

func (e *bytesByteTupleEnumerator) Current() (string, Value) {
	return e.name, e.value
}
