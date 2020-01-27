package rel

import (
	"reflect"
)

// CharAttr is the standard name for the value-attr of a character tuple.
const CharAttr = "@char"

// String is a set of Values.
type String struct {
	s      []rune
	offset int
}

// NewString constructs an array as a relation.
func NewString(s []rune) Set {
	if len(s) == 0 {
		return None
	}
	return String{s: s}
}

// NewString constructs an array as a relation.
func NewOffsetString(s []rune, offset int) Set {
	if len(s) == 0 {
		return None
	}
	return String{s: s, offset: offset}
}

// Hash computes a hash for a String.
func (s String) Hash(seed uintptr) uintptr {
	h := seed
	for e := s.Enumerator(); e.MoveNext(); {
		h ^= e.Current().Hash(seed)
	}
	return h
}

// Equal tests two Sets for equality. Any other type returns false.
func (s String) Equal(v interface{}) bool {
	switch x := v.(type) {
	case String:
		if len(s.s) != len(x.s) {
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
func (s String) Eval(_, _ Scope) (Value, error) {
	return s, nil
}

var stringKind = registerKind(204, reflect.TypeOf(String{}))

// Kind returns a number that is unique for each major kind of Value.
func (s String) Kind() int {
	return stringKind
}

// Bool returns true iff the tuple has attributes.
func (s String) Bool() bool {
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

	return string(s.s) < string(v.(*String).s)
}

// Negate returns {(negateTag): s}.
func (s String) Negate() Value {
	return NewTuple(NewAttr(negateTag, s))
}

// Export exports a String as a string.
func (s String) Export() interface{} {
	return s.s
}

// Count returns the number of elements in the String.
func (s String) Count() int {
	return len(s.s)
}

// Has returns true iff the given Value is in the String.
func (s String) Has(value Value) bool {
	if pos, char, ok := isStringTuple(value); ok {
		if s.offset <= pos && pos < s.offset+len(s.s) {
			return char == s.s[pos-s.offset]
		}
	}
	return false
}

func (s String) with(index int, char rune) Set {
	if s.index(index) == len(s.s) {
		return String{s: append(s.s, char), offset: s.offset}
	} else if index == s.offset-1 {
		return String{
			s:      append(append(make([]rune, 0, 1+len(s.s)), char), s.s...),
			offset: s.offset - 1,
		}
	}
	return newSetFromSet(s).With(newStringTuple(index, char))
}

// With returns the original String with given value added. Iff the value was
// already present, the original String is returned.
func (s String) With(value Value) Set {
	if index, char, ok := isStringTuple(value); ok {
		return s.with(index, char)
	}
	return newSetFromSet(s).With(value)
}

// Without returns the original String without the given value. Iff the value
// was already absent, the original String is returned.
func (s String) Without(value Value) Set {
	if pos, char, ok := isStringTuple(value); ok {
		if i := s.index(pos); i >= 0 && char == s.s[i] {
			if pos == s.offset+i {
				return String{s: s.s[:i], offset: s.offset}
			}
			return newSetFromSet(s).Without(value)
		}
	}
	return s
}

// Map maps values per f.
func (s String) Map(f func(v Value) Value) Set {
	result := NewSet()
	for e := s.Enumerator(); e.MoveNext(); {
		result = result.With(f(e.Current()))
	}
	return result
}

// Where returns a new String with all the Values satisfying predicate p.
func (s String) Where(p func(v Value) bool) Set {
	result := Set(s)
	for e := s.Enumerator(); e.MoveNext(); {
		value := e.Current()
		if !p(value) {
			result = result.Without(value)
		}
	}
	return result
}

// Call ...
func (s String) Call(arg Value) Value {
	i := int(arg.(Number).Float64())
	return NewNumber(float64(string(s.s)[i-s.offset]))
}

func (s String) index(pos int) int {
	pos -= s.offset
	if 0 <= pos && pos < len(s.s) {
		return pos
	}
	return -1
}

// StringEnumerator represents an enumerator over a String.
type StringEnumerator struct {
	s []rune
	i int
}

// MoveNext moves the enumerator to the next Value.
func (e *StringEnumerator) MoveNext() bool {
	e.i++
	return e.i < len(e.s)
}

// Current returns the enumerator's current Value.
func (e *StringEnumerator) Current() Value {
	return newStringTuple(e.i, e.s[e.i])
}

// Enumerator returns an enumerator over the Values in the String.
func (s String) Enumerator() ValueEnumerator {
	return &StringEnumerator{s.s, -1}
}

func newStringTuple(index int, char rune) Tuple {
	return NewTuple(
		NewAttr("@", NewNumber(float64(index))),
		NewAttr(CharAttr, NewNumber(float64(char))),
	)
}

func isStringTuple(v Value) (index int, char rune, is bool) {
	is = stringTupleMatcher(func(i int, c rune) { index = i; char = c }).Match(v)
	return
}

func stringTupleMatcher(match func(index int, char rune)) TupleMatcher {
	n := 0
	var index int
	var char rune
	check := func() {
		if n == 1 {
			match(index, char)
		}
		n++
	}
	return NewTupleMatcher(
		map[string]Matcher{
			"@":      MatchInt(func(i int) { index = i; check() }),
			CharAttr: MatchInt(func(i int) { char = rune(i); check() }),
		},
		Lit(EmptyTuple),
	)
}

// func stringSet(s Set) Set {
// 	if s, ok := s.(String); ok {
// 		return s
// 	}
// 	if !s.Bool() {
// 		return s
// 	}

// 	var result String
// 	matcher := stringTupleMatcher(func(index int, char rune) {
// 		result = result.with(index, char)
// 	})
// }
