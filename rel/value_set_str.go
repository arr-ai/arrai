package rel

import "encoding/json"

// CharAttr is the standard name for the value-attr of a character tuple.
const CharAttr = "@char"

// String is a set of Values.
type String struct {
	s []rune
}

// NewString constructs an array as a relation.
func NewString(s []rune) Set {
	if len(s) == 0 {
		return None
	}
	return &String{s}
}

// Hash computes a hash for a String.
func (s *String) Hash(seed uintptr) uintptr {
	h := seed
	for e := s.Enumerator(); e.MoveNext(); {
		h ^= e.Current().Hash(0)
	}
	return h
}

// Equal tests two Sets for equality. Any other type returns false.
func (s *String) Equal(v interface{}) bool {
	switch x := v.(type) {
	case *String:
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
func (s *String) String() string {
	j, err := json.Marshal(string(s.s))
	if err != nil {
		panic(err)
	}
	return string(j)
}

// Eval returns the string.
func (s *String) Eval(_, _ *Scope) (Value, error) {
	return s, nil
}

// Kind returns a number that is unique for each major kind of Value.
func (s *String) Kind() int {
	return 204
}

// Bool returns true iff the tuple has attributes.
func (s *String) Bool() bool {
	if len(s.s) == 0 {
		panic("Empty string not allowed (should be == None)")
	}
	return true
}

// Less returns true iff v is not a number or tuple, or v is a tuple and t
// precedes v in a lexicographical comparison of their name/value pairs.
func (s *String) Less(v Value) bool {
	if s.Kind() != v.Kind() {
		return s.Kind() < v.Kind()
	}

	return string(s.s) < string(v.(*String).s)
}

// Negate returns {(negateTag): s}.
func (s *String) Negate() Value {
	return NewTuple(NewAttr(negateTag, s))
}

// Export exports a String as a string.
func (s *String) Export() interface{} {
	return s.s
}

// Count returns the number of elements in the String.
func (s *String) Count() int {
	return len(s.s)
}

// Has returns true iff the given Value is in the String.
func (s *String) Has(value Value) bool {
	if pos, char, ok := isStringTuple(value); ok && pos < uint(len(s.s)) {
		return char == s.s[pos]
	}
	return false
}

// With returns the original String with given value added. Iff the value was
// already present, the original String is returned.
func (s *String) With(value Value) Set {
	if pos, char, ok := isStringTuple(value); ok {
		if pos == uint(len(s.s)) {
			return &String{append(s.s, char)}
		}
	}
	return newSetFromSet(s).With(value)
}

// Without returns the original String without the given value. Iff the value was
// already absent, the original String is returned.
func (s *String) Without(value Value) Set {
	if pos, char, ok := isStringTuple(value); ok {
		if pos < uint(len(s.s)) && char == s.s[pos] {
			if pos == uint(len(s.s)-1) {
				return &String{s.s[:len(s.s)-1]}
			}
			return newSetFromSet(s).Without(value)
		}
	}
	return s
}

// Map maps values per f.
func (s *String) Map(f func(v Value) Value) Set {
	result := NewSet()
	for e := s.Enumerator(); e.MoveNext(); {
		result = result.With(f(e.Current()))
	}
	return result
}

// Where returns a new String with all the Values satisfying predicate p.
func (s *String) Where(p func(v Value) bool) Set {
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
func (s *String) Call(arg Value) Value {
	return NewNumber(float64(string(s.s)[int(arg.(Number).Float64())]))
}

// StringEnumerator represents an enumerator over a String.
type StringEnumerator struct {
	s []rune
	i uint
}

// MoveNext moves the enumerator to the next Value.
func (e *StringEnumerator) MoveNext() bool {
	e.i++
	return e.i < uint(len(e.s))
}

// Current returns the enumerator's current Value.
func (e *StringEnumerator) Current() Value {
	return newStringTuple(e.i, e.s[e.i])
}

// Enumerator returns an enumerator over the Values in the String.
func (s *String) Enumerator() ValueEnumerator {
	return &StringEnumerator{s.s, ^uint(0)}
}

func newStringTuple(pos uint, char rune) Tuple {
	return NewTuple(
		NewAttr("@", NewNumber(float64(pos))),
		NewAttr(CharAttr, NewNumber(float64(char))),
	)
}

func isStringTuple(v Value) (index uint, char rune, is bool) {
	is = NewTupleMatcher(
		map[string]Matcher{
			"@":      MatchInt(func(i int) { index = uint(i) }),
			CharAttr: MatchInt(func(i int) { char = rune(i) }),
		},
		Lit(EmptyTuple),
	).Match(v)
	return
}
