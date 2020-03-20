package rel

import (
	"reflect"
)

// StringCharAttr is the standard name for the value-attr of a character tuple.
const StringCharAttr = "@char"

// String is a set of Values.
type String struct {
	s      []rune
	offset int
}

// NewString constructs an array as a relation.
func NewString(s []rune) Set {
	return NewOffsetString(s, 0)
}

// NewString constructs an array as a relation.
func NewOffsetString(s []rune, offset int) Set {
	if len(s) == 0 {
		return None
	}
	return String{s: s, offset: offset}
}

func AsString(s Set) (String, bool) {
	if s, ok := s.(String); ok {
		return s, true
	}
	if i := s.Enumerator(); i.MoveNext() {
		t, is := i.Current().(StringCharTuple)
		if !is {
			return String{}, false
		}

		middleIndex := s.Count()
		strs := make([]rune, 2*middleIndex)
		strs[middleIndex] = t.char
		anchorOffset, minOffset := t.at, t.at
		lowestIndex, highestIndex := middleIndex, middleIndex
		for i.MoveNext() {
			if t, is = i.Current().(StringCharTuple); !is {
				return String{}, false
			}
			if t.at < minOffset {
				minOffset = t.at
			}
			sliceIndex := middleIndex - (anchorOffset - t.at)
			strs[sliceIndex] = t.char

			if sliceIndex < lowestIndex {
				lowestIndex = sliceIndex
			}

			if sliceIndex > highestIndex {
				highestIndex = sliceIndex
			}
		}

		return NewOffsetString(strs[lowestIndex:highestIndex+1], minOffset).(String), true
	}
	return String{}, true
}

// Hash computes a hash for a String.
func (s String) Hash(seed uintptr) uintptr {
	// TODO: Optimize.
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
func (s String) Eval(_ Scope) (Value, error) {
	return s, nil
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
func (s String) Export() interface{} {
	return string(s.s)
}

// Count returns the number of elements in the String.
func (s String) Count() int {
	return len(s.s)
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
	if s.index(index) == len(s.s) {
		return String{s: append(s.s, char), offset: s.offset}
	} else if index == s.offset-1 {
		return String{
			s:      append(append(make([]rune, 0, 1+len(s.s)), char), s.s...),
			offset: s.offset - 1,
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
		if i := s.index(t.at); i >= 0 && t.char == s.s[i] {
			if t.at == s.offset+i {
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
	if e.i >= len(e.s.s)-1 {
		return false
	}
	e.i++
	return true
}

// Current returns the enumerator's current Value.
func (e *stringValueEnumerator) Current() Value {
	return NewStringCharTuple(e.i, e.s.s[e.i])
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
