package rel

import (
	"context"
	"fmt"
	"reflect"

	"github.com/arr-ai/hash"
	"github.com/arr-ai/wbnf/parser"
)

// StringCharAttr is the standard name for the value-attr of a character tuple.
const StringCharAttr = "@char"

// String is a set of Values.
type String struct {
	s      []rune
	offset int
	holes  int
}

// NewString constructs a string as a relation.
func NewString(s []rune) Set {
	return NewOffsetString(s, 0)
}

// NewOffsetString constructs an offset string as a relation.
func NewOffsetString(s []rune, offset int) Set {
	if len(s) == 0 {
		return None
	}
	holes := 0
	for _, r := range s {
		if r < 0 {
			holes++
		}
	}
	return String{s: s, offset: offset, holes: holes}
}

func AsString(s Set) (String, bool) { //nolint:dupl
	if s, ok := s.(String); ok {
		return s, true
	}
	n := s.Count()
	if n == 0 {
		return String{}, true
	}
	tuples := make(stringCharTupleArray, 0, n)
	minAt := int(^uint(0) >> 1)
	maxAt := -minAt - 1
	for i := s.Enumerator(); i.MoveNext(); {
		t, is := i.Current().(StringCharTuple)
		if !is {
			return String{}, false
		}
		if minAt > t.at {
			minAt = t.at
		}
		if maxAt < t.at {
			maxAt = t.at
		}
		tuples = append(tuples, t)
	}
	str := make([]rune, maxAt-minAt+1)
	for i := range str {
		str[i] = -1
	}
	for _, t := range tuples {
		str[t.at-minAt] = t.char
	}
	return String{s: str, offset: minAt, holes: len(str) - n}, true
}

type stringCharTupleArray []StringCharTuple

func (a stringCharTupleArray) Len() int {
	return len(a)
}

func (a stringCharTupleArray) Less(i, j int) bool {
	return a[i].at < a[j].at
}

func (a stringCharTupleArray) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Hash computes a hash for a String.
func (s String) Hash(seed uintptr) uintptr {
	// TODO: implement a []rune-friendly hash function.
	return hash.String(string(s.s), seed)
}

// Equal tests two Sets for equality. Any other type returns false.
func (s String) Equal(v interface{}) bool {
	switch x := v.(type) {
	case String:
		if len(s.s) != len(x.s) || s.offset != x.offset {
			return false
		}
		for i, c := range s.s {
			if c != x.s[i] {
				return false
			}
		}
		return true
	case Set:
		return newSetFromSet(s).Equal(x)
	}
	return false
}

// String returns a string representation of a String.
func (s String) String() string {
	return string(s.s)
}

// Eval returns the string.
func (s String) Eval(ctx context.Context, _ Scope) (Value, error) {
	return s, nil
}

// Source returns a scanner locating the String's source code.
func (s String) Source() parser.Scanner {
	return *parser.NewScanner("")
}

var stringKind = registerKind(204, reflect.TypeOf(String{}))

// Kind returns a number that is unique for each major kind of Value.
func (s String) Kind() int {
	return stringKind
}

// Bool returns true iff the tuple has attributes.
func (s String) IsTrue() bool {
	if len(s.s) == 0 {
		panic("Empty string not allowed (should be == None)")
	}
	return true
}

// Less returns true iff v is not a number or tuple, or v is a tuple and t
// precedes v in a lexicographical comparison of their name/value pairs.
func (s String) Less(v Value) bool {
	if s.Kind() != v.Kind() {
		return s.Kind() < v.Kind()
	}

	return s.String() < v.(String).String()
}

// Negate returns {(negateTag): s}.
func (s String) Negate() Value {
	return NewTuple(NewAttr(negateTag, s))
}

// Export exports a String as a string.
func (s String) Export(_ context.Context) interface{} {
	return string(s.s)
}

// Count returns the number of elements in the String.
func (s String) Count() int {
	return len(s.s) - s.holes
}

// Has returns true iff the given Value is in the String.
func (s String) Has(value Value) bool {
	if t, ok := value.(StringCharTuple); ok {
		if s.offset <= t.at && t.at < s.offset+len(s.s) {
			return t.char == s.s[t.at-s.offset]
		}
	}
	return false
}

func (s String) with(index int, char rune) Set {
	switch {
	case s.index(index) == len(s.s):
		return String{s: append(s.s, char), offset: s.offset, holes: s.holes}
	case index == s.offset-1:
		return String{
			s:      append(append(make([]rune, 0, 1+len(s.s)), char), s.s...),
			offset: s.offset - 1,
			holes:  s.holes,
		}
	}
	return newSetFromSet(s).With(NewStringCharTuple(index, char))
}

// With returns the original String with given value added. Iff the value was
// already present, the original String is returned.
func (s String) With(value Value) Set {
	if t, ok := value.(StringCharTuple); ok {
		return s.with(t.at, t.char)
	}
	return newSetFromSet(s).With(value)
}

// Without returns the original String without the given value. Iff the value
// was already absent, the original String is returned.
func (s String) Without(value Value) Set {
	if t, ok := value.(StringCharTuple); ok {
		i := s.index(t.at)
		switch {
		case i == 0:
			return String{s: s.s[:i], offset: s.offset, holes: s.holes}
		case i == len(s.s)-1:
			return String{s: s.s[i : len(s.s)-1], offset: s.offset, holes: s.holes}
		case 0 < i && i < len(s.s)-1:
			if t.char == s.s[i] {
				newS := make([]rune, len(s.s))
				copy(newS, s.s)
				newS[i] = -1
				return String{s: newS, offset: s.offset, holes: s.holes + 1}
			}
		}
	}
	return s
}

// Map maps values per f.
func (s String) Map(f func(v Value) (Value, error)) (Set, error) {
	result := None
	for e := s.Enumerator().(*stringValueEnumerator); e.MoveNext(); {
		v, err := f(e.currentStringCharTuple())
		if err != nil {
			return nil, err
		}
		result = result.With(v)
	}
	return result, nil
}

// Where returns a new String with all the Values satisfying predicate p.
func (s String) Where(p func(v Value) (bool, error)) (Set, error) {
	values := make([]Value, 0, s.Count())
	for e := s.Enumerator().(*stringValueEnumerator); e.MoveNext(); {
		value := e.currentStringCharTuple()
		matches, err := p(value)
		if err != nil {
			return nil, err
		}
		if matches {
			values = append(values, value)
		}
	}
	return NewSet(values...)
}

func (s String) CallAll(_ context.Context, arg Value) (Set, error) {
	n, ok := arg.(Number)
	if !ok {
		return nil, fmt.Errorf("arg to CallAll must be a number, not %s", ValueTypeAsString(arg))
	}
	i := int(n.Float64()) - s.offset
	if i < 0 || i >= len(s.s) {
		return None, nil
	}
	return NewSet(NewNumber(float64(string(s.s)[i])))
}

func (s String) index(pos int) int {
	pos -= s.offset
	if 0 <= pos && pos <= len(s.s) {
		return pos
	}
	return -1
}

// Enumerator returns an enumerator over the Values in the String.
func (s String) Enumerator() ValueEnumerator {
	return &stringValueEnumerator{s: s, i: -1}
}

func (s String) ArrayEnumerator() (OffsetValueEnumerator, bool) {
	return &stringOffsetValueEnumerator{stringValueEnumerator{s: s, i: -1}}, true
}

// StringEnumerator represents an enumerator over a String.
type stringValueEnumerator struct {
	s String
	i int
}

// MoveNext moves the enumerator to the next Value.
func (e *stringValueEnumerator) MoveNext() bool {
	for e.i < len(e.s.s)-1 {
		e.i++
		if e.s.s[e.i] >= 0 {
			return true
		}
	}
	return false
}

// Current returns the enumerator's current Value.
func (e *stringValueEnumerator) Current() Value {
	return NewStringCharTuple(e.s.offset+e.i, e.s.s[e.i])
}

// This version avoids mallocs.
func (e *stringValueEnumerator) currentStringCharTuple() StringCharTuple {
	return NewStringCharTuple(e.s.offset+e.i, e.s.s[e.i])
}

type stringOffsetValueEnumerator struct {
	stringValueEnumerator
}

func (e *stringOffsetValueEnumerator) Current() Value {
	return NewNumber(float64(e.s.s[e.i]))
}

func (e *stringOffsetValueEnumerator) Offset() int {
	return e.s.offset + e.i
}
