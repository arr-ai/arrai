package rel

import (
	"context"
	"fmt"
	"reflect"

	"github.com/arr-ai/hash"
	"github.com/arr-ai/wbnf/parser"

	"github.com/arr-ai/arrai/pkg/fu"
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

func asString(values ...Value) String {
	n := len(values)
	tuples := make([]StringCharTuple, 0, n)
	minAt := int(^uint(0) >> 1)
	maxAt := -minAt - 1
	for _, v := range values {
		t := v.(StringCharTuple)
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
	return String{s: str, offset: minAt, holes: len(str) - n}
}

// AsString returns String and the empty set as String or false otherwise.
func AsString(v Value) (String, bool) {
	switch s := v.(type) {
	case String:
		return s, true
	case Set:
		if !s.IsTrue() {
			return String{}, true
		}
	}
	return String{}, false
}

// Hash computes a hash for a String.
func (s String) Hash(seed uintptr) uintptr {
	// TODO: implement a []rune-friendly hash function.
	return hash.String(string(s.s), seed)
}

// Equal tests two Sets for equality. Any other type returns false.
func (s String) Equal(v interface{}) bool {
	t, is := v.(String)
	return is && s.EqualString(t)
}

func (s String) EqualString(t String) bool {
	if s.offset != t.offset || s.holes != t.holes || len(s.s) != len(t.s) {
		return false
	}
	for i, r := range s.s {
		if r != t.s[i] {
			return false
		}
	}
	return true
}

// String returns a string representation of a String.
func (s String) String() string {
	return string(s.s)
}

func (s String) Format(f fmt.State, verb rune) {
	if verb == 's' {
		fu.WriteString(f, string(s.s))
	} else {
		reprString(s, f)
	}
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

func (String) getSetBuilder() setBuilder {
	return newGenericTypeSetBuilder()
}

func (String) getBucket() fmt.Stringer {
	return genericType
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

func (s String) with(at int, char rune) Set {
	i := s.index(at)
	switch {
	case 0 <= i && i < len(s.s) && s.s[i] == char:
		return s
	case i == len(s.s):
		return String{s: append(s.s, char), offset: s.offset, holes: s.holes}
	case at == s.offset-1:
		return String{
			s:      append(append(make([]rune, 0, 1+len(s.s)), char), s.s...),
			offset: s.offset - 1,
			holes:  s.holes,
		}
	}
	// TODO: Support adding holes and doubling up chars, removing the need to
	// call newGenericSetFromSet here.
	return newGenericSetFromSet(s).With(NewStringCharTuple(at, char))
}

// With returns the original String with given value added. Iff the value was
// already present, the original String is returned.
func (s String) With(value Value) Set {
	if t, ok := value.(StringCharTuple); ok {
		return s.with(t.at, t.char)
	}
	return toUnionSetWithItem(s, value)
}

// Without returns the original String without the given value. Iff the value
// was already absent, the original String is returned.
func (s String) Without(value Value) Set {
	if t, ok := value.(StringCharTuple); ok {
		i := s.index(t.at)
		switch {
		case i == 0:
			s = String{s: s.s[:i], offset: s.offset, holes: s.holes}
		case i == len(s.s)-1:
			s = String{s: s.s[i : len(s.s)-1], offset: s.offset, holes: s.holes}
		case 0 < i && i < len(s.s)-1:
			if t.char == s.s[i] {
				newS := make([]rune, len(s.s))
				copy(newS, s.s)
				newS[i] = -1
				s = String{s: newS, offset: s.offset, holes: s.holes + 1}
			}
		}
	}
	if s.Count() == 0 {
		return None
	}
	return s
}

// Map maps values per f.
func (s String) Map(f func(v Value) (Value, error)) (Set, error) {
	b := NewSetBuilder()
	for e := s.Enumerator().(*stringEnumerator); e.MoveNext(); {
		v, err := f(e.Current())
		if err != nil {
			return nil, err
		}
		b.Add(v)
	}
	return b.Finish()
}

// Where returns a new String with all the Values satisfying predicate p.
func (s String) Where(p func(v Value) (bool, error)) (Set, error) {
	b := NewSetBuilder()
	for e := s.Enumerator().(*stringEnumerator); e.MoveNext(); {
		value := e.Current()
		matches, err := p(value)
		if err != nil {
			return nil, err
		}
		if matches {
			b.Add(value)
		}
	}
	return b.Finish()
}

func (s String) CallAll(_ context.Context, arg Value, b SetBuilder) error {
	if n, ok := arg.(Number); ok {
		if i, is := n.Int(); is {
			i -= s.offset
			if 0 <= i && i < len(s.s) {
				b.Add(NewNumber(float64(s.s[i])))
			}
		}
	}
	return nil
}

func (String) unionSetSubsetBucket() string {
	return StringCharTuple{}.getBucket().String()
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
	return &stringEnumerator{s: s, i: -1}
}

func (s String) ArrayEnumerator() ValueEnumerator {
	return &stringValueEnumerator{s.Enumerator().(*stringEnumerator)}
}

// StringEnumerator represents an enumerator over a String.
type stringEnumerator struct {
	s String
	i int
}

// MoveNext moves the enumerator to the next Value.
func (e *stringEnumerator) MoveNext() bool {
	for e.i < len(e.s.s)-1 {
		e.i++
		if e.s.s[e.i] >= 0 {
			return true
		}
	}
	return false
}

// Current returns the enumerator's current Value.
func (e *stringEnumerator) Current() Value {
	return NewStringCharTuple(e.s.offset+e.i, e.s.s[e.i])
}

type stringValueEnumerator struct {
	*stringEnumerator
}

func (e *stringValueEnumerator) Current() Value {
	return NewNumber(float64(e.s.s[e.i]))
}
