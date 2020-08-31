package rel //nolint:dupl

import (
	"context"
	"fmt"
	"reflect"

	"github.com/arr-ai/hash"
	"github.com/arr-ai/wbnf/parser"
)

// StringCharTuple represents a tuple of the form (@: at, @char: char).
type StringCharTuple struct {
	at   int
	char rune
}

// NewStringCharTuple constructs a CharTuple.
func NewStringCharTuple(at int, char rune) StringCharTuple {
	return StringCharTuple{at: at, char: char}
}

func newCharTupleFromTuple(t Tuple) (StringCharTuple, bool) {
	var at int
	var char rune
	m := NewTupleMatcher(
		map[string]Matcher{
			"@":            MatchInt(func(i int) { at = i }),
			StringCharAttr: MatchInt(func(i int) { char = rune(i) }),
		},
		Lit(EmptyTuple),
	)
	if m.Match(t) {
		return NewStringCharTuple(at, char), true
	}
	return StringCharTuple{}, false
}

func maybeNewCharTupleFromTuple(t Tuple) Tuple {
	if t, ok := newCharTupleFromTuple(t); ok {
		return t
	}
	return t
}

func (t StringCharTuple) asGenericTuple() Tuple {
	return newTuple(
		NewIntAttr("@", t.at),
		NewIntAttr(StringCharAttr, int(t.char)),
	)
}

// Hash computes a hash for a CharTuple.
func (t StringCharTuple) Hash(seed uintptr) uintptr {
	h := hash.Int(t.at, seed)
	return hash.Int32(t.char, h)
}

// Equal tests two Tuples for equality. Any other type returns false.
func (t StringCharTuple) Equal(v interface{}) bool {
	if u, ok := v.(StringCharTuple); ok {
		return t == u
	}
	return false
}

// String returns a string representation of a Tuple.
func (t StringCharTuple) String() string {
	return fmt.Sprintf("(@: %d, %s: %d)", t.at, StringCharAttr, t.char)
}

// Eval returns the tuple.
func (t StringCharTuple) Eval(ctx context.Context, local Scope) (Value, error) {
	return t, nil
}

// Source returns a scanner locating the StringCharTuple's source code.
func (t StringCharTuple) Source() parser.Scanner {
	return *parser.NewScanner("")
}

var stringCharTupleKind = registerKind(301, reflect.TypeOf((*StringCharTuple)(nil)))

// Kind returns a number that is unique for each major kind of Value.
func (t StringCharTuple) Kind() int {
	return stringCharTupleKind
}

// Bool returns true iff the tuple has attributes.
func (t StringCharTuple) IsTrue() bool {
	return true
}

// Less returns true iff v is not a number or tuple, or v is a tuple and t
// precedes v in a lexicographical comparison of their name/value pairs.
func (t StringCharTuple) Less(v Value) bool {
	if t.Kind() != v.Kind() {
		return t.Kind() < v.Kind()
	}
	u := v.(StringCharTuple)
	if t.at != u.at {
		return t.at < u.at
	}
	return t.char < u.char
}

func (t StringCharTuple) Negate() Value {
	return StringCharTuple{at: -t.at, char: -t.char}
}

// Export exports a Tuple.
func (t StringCharTuple) Export(_ context.Context) interface{} {
	return map[string]interface{}{
		"@":            t.at,
		StringCharAttr: t.char,
	}
}

// Count returns how many attributes are in the Tuple.
func (t StringCharTuple) Count() int {
	return 2
}

// Get returns the Value associated with a name, and true iff it was found.
func (t StringCharTuple) Get(name string) (Value, bool) {
	switch name {
	case "@":
		return NewNumber(float64(t.at)), true
	case StringCharAttr:
		return NewNumber(float64(t.char)), true
	}
	return nil, false
}

// MustGet returns e.Get(name) or panics if an error occurs.
func (t StringCharTuple) MustGet(name string) Value {
	if v, has := t.Get(name); has {
		return v
	}
	panic(fmt.Errorf("%q not found", name))
}

// With returns a Tuple with all name/Value pairs in t (except the one for the
// given name, if present) with the addition of the given name/Value pair.
func (t StringCharTuple) With(name string, value Value) Tuple {
	return maybeNewCharTupleFromTuple(t.asGenericTuple().With(name, value))
}

// Without returns a Tuple with all name/Value pairs in t exception the one of
// the given name.
func (t StringCharTuple) Without(name string) Tuple {
	if name == "@" || name == StringCharAttr {
		return t.asGenericTuple().Without(name)
	}
	return t
}

func (t StringCharTuple) Map(f func(Value) (Value, error)) (Tuple, error) { //nolint:dupl
	at, err := f(NewNumber(float64(t.at)))
	if err != nil {
		return nil, err
	}
	char, err := f(NewNumber(float64(t.char)))
	if err != nil {
		return nil, err
	}
	if at, ok := at.(Number); ok {
		if at, is := at.Int(); is {
			if char, ok := char.(Number); ok {
				if char, is := char.Int(); is {
					return NewStringCharTuple(at, rune(char)), nil
				}
			}
		}
	}
	return t.asGenericTuple().Map(f)
}

// HasName returns true iff the Tuple has an attribute with the given name.
func (t StringCharTuple) HasName(name string) bool {
	return name == "@" || name == StringCharAttr
}

// Names returns the attribute names.
func (t StringCharTuple) Names() Names {
	return NewNames("@", StringCharAttr)
}

// Project returns a tuple with the given names from this tuple, or nil if any
// name wasn't found.
func (t StringCharTuple) Project(names Names) Tuple {
	if names.Has("@") && names.Has(StringCharAttr) {
		return t
	}
	return t.asGenericTuple().Project(names)
}

// Enumerator returns an enumerator over the Values in the CharTuple.
func (t StringCharTuple) Enumerator() AttrEnumerator {
	return &stringCharTupleEnumerator{t: t, i: -1}
}

type stringCharTupleEnumerator struct {
	t     StringCharTuple
	i     int
	name  string
	value Value
}

func (e *stringCharTupleEnumerator) MoveNext() bool {
	if e.i == 1 {
		return false
	}
	e.i++
	switch e.i {
	case 0:
		e.name = "@"
		e.value = NewNumber(float64(e.t.at))
	case 1:
		e.name = StringCharAttr
		e.value = NewNumber(float64(e.t.char))
	}
	return true
}

func (e *stringCharTupleEnumerator) Current() (string, Value) {
	return e.name, e.value
}
