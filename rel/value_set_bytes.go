package rel

import (
	"reflect"

	"github.com/arr-ai/wbnf/parser"
)

// BytesByteAttr is the standard name for the value-attr of a character tuple.
const BytesByteAttr = "@byte"

// Bytes is a set of Values.
type Bytes struct {
	b      []byte
	offset int
}

// NewBytes constructs an array as a relation.
func NewBytes(b []byte) Set {
	if len(b) == 0 {
		return None
	}
	return Bytes{b: b}
}

// NewBytes constructs an array as a relation.
func NewOffsetBytes(b []byte, offset int) Set {
	if len(b) == 0 {
		return None
	}
	return Bytes{b: b, offset: offset}
}

func AsBytes(s Set) (Bytes, bool) { //nolint:dupl
	if b, ok := s.(Bytes); ok {
		return b, true
	}
	if i := s.Enumerator(); i.MoveNext() {
		t, is := i.Current().(BytesByteTuple)
		if !is {
			return Bytes{}, false
		}

		middleIndex := s.Count()
		strs := make([]byte, 2*middleIndex)
		strs[middleIndex] = t.byteval
		anchorOffset, minOffset := t.at, t.at
		lowestIndex, highestIndex := middleIndex, middleIndex
		for i.MoveNext() {
			if t, is = i.Current().(BytesByteTuple); !is {
				return Bytes{}, false
			}
			if t.at < minOffset {
				minOffset = t.at
			}
			sliceIndex := middleIndex - (anchorOffset - t.at)
			strs[sliceIndex] = t.byteval

			if sliceIndex < lowestIndex {
				lowestIndex = sliceIndex
			}

			if sliceIndex > highestIndex {
				highestIndex = sliceIndex
			}
		}

		return NewOffsetBytes(strs[lowestIndex:highestIndex+1], minOffset).(Bytes), true
	}
	return Bytes{}, true
}

// Bytes returns the bytes of b. The caller must not modify the contents.
func (b Bytes) Bytes() []byte {
	return b.b
}

// Hash computes a hash for a Bytes.
func (b Bytes) Hash(seed uintptr) uintptr {
	h := seed
	for e := b.Enumerator(); e.MoveNext(); {
		h ^= e.Current().Hash(seed)
	}
	return h
}

// Equal tests two Sets for equality. Any other type returns false.
func (b Bytes) Equal(v interface{}) bool {
	switch x := v.(type) {
	case Bytes:
		if len(b.b) != len(x.b) {
			return false
		}
		for i, c := range b.b {
			if c != x.b[i] {
				return false
			}
		}
		return true
	case Set:
		return newSetFromSet(b).Equal(x)
	}
	return false
}

// String returns a string representation of a Bytes.
func (b Bytes) String() string {
	return string(b.b)
}

// Eval returns the string.
func (b Bytes) Eval(_ Scope) (Value, error) {
	return b, nil
}

// Scanner returns the scanner of Bytes.
func (b Bytes) Scanner() parser.Scanner {
	return *parser.NewScanner("")
}

var bytesKind = registerKind(207, reflect.TypeOf(Bytes{}))

// Kind returns a number that is unique for each major kind of Value.
func (b Bytes) Kind() int {
	return bytesKind
}

// Bool returns true iff the tuple has attributes.
func (b Bytes) IsTrue() bool {
	if len(b.b) == 0 {
		panic("Empty string not allowed (should be == None)")
	}
	return true
}

// Less returns true iff v is not a number or tuple, or v is a tuple and t
// precedes v in a lexicographical comparison of their name/value pairs.
func (b Bytes) Less(v Value) bool {
	if b.Kind() != v.Kind() {
		return b.Kind() < v.Kind()
	}

	return string(b.b) < string(v.(*Bytes).b)
}

// Negate returns {(negateTag): b}.
func (b Bytes) Negate() Value {
	return NewTuple(NewAttr(negateTag, b))
}

// Export exports a Bytes as a string.
func (b Bytes) Export() interface{} {
	return string(b.b)
}

// Count returns the number of elements in the Bytes.
func (b Bytes) Count() int {
	return len(b.b)
}

// Has returns true iff the given Value is in the Bytes.
func (b Bytes) Has(value Value) bool {
	if pos, byt, ok := isBytesTuple(value); ok {
		if b.offset <= pos && pos < b.offset+len(b.b) {
			return byt == b.b[pos-b.offset]
		}
	}
	return false
}

func (b Bytes) with(index int, byt byte) Set {
	if b.index(index) == len(b.b) {
		return Bytes{b: append(b.b, byt), offset: b.offset}
	} else if index == b.offset-1 {
		return Bytes{
			b:      append(append(make([]byte, 0, 1+len(b.b)), byt), b.b...),
			offset: b.offset - 1,
		}
	}
	return newSetFromSet(b).With(newBytesTuple(index, byt))
}

// With returns the original Bytes with given value added. Iff the value was
// already present, the original Bytes is returned.
func (b Bytes) With(value Value) Set {
	if index, byt, ok := isBytesTuple(value); ok {
		return b.with(index, byt)
	}
	return newSetFromSet(b).With(value)
}

// Without returns the original Bytes without the given value. Iff the value
// was already absent, the original Bytes is returned.
func (b Bytes) Without(value Value) Set {
	if pos, byt, ok := isBytesTuple(value); ok {
		if i := b.index(pos); i >= 0 && i < len(b.b) && byt == b.b[i] {
			if pos == b.offset+i {
				return Bytes{b: b.b[:i], offset: b.offset}
			}
			return newSetFromSet(b).Without(value)
		}
	}
	return b
}

// Map maps values per f.
func (b Bytes) Map(f func(v Value) Value) Set {
	result := NewSet()
	for e := b.Enumerator(); e.MoveNext(); {
		result = result.With(f(e.Current()))
	}
	return result
}

// Where returns a new Bytes with all the Values satisfying predicate p.
func (b Bytes) Where(p func(v Value) bool) Set {
	result := Set(b)
	for e := b.Enumerator(); e.MoveNext(); {
		value := e.Current()
		if !p(value) {
			result = result.Without(value)
		}
	}
	return result
}

// Call ...
func (b Bytes) Call(arg Value) Value {
	i := int(arg.(Number).Float64())
	return NewNumber(float64(string(b.b)[i-b.offset]))
}

func (b Bytes) CallSlice(start, end Value, step int, inclusive bool) (Set, error) {
	indexes, err := resolveArrayIndexes(start, end, step, b.offset, len(b.b), inclusive)
	if err != nil {
		return nil, err
	}
	slice := make([]byte, 0, len(indexes))
	for _, i := range indexes {
		slice = append(slice, b.b[i-b.offset])
	}
	return NewOffsetBytes(slice, b.offset), nil
}

func (b Bytes) index(pos int) int {
	pos -= b.offset
	if 0 <= pos && pos <= len(b.b) {
		return pos
	}
	return -1
}

// BytesEnumerator represents an enumerator over a Bytes.
type BytesEnumerator struct {
	b []byte
	i int
}

// MoveNext moves the enumerator to the next Value.
func (e *BytesEnumerator) MoveNext() bool {
	e.i++
	return e.i < len(e.b)
}

// Current returns the enumerator'b current Value.
func (e *BytesEnumerator) Current() Value {
	return newBytesTuple(e.i, e.b[e.i])
}

// Enumerator returns an enumerator over the Values in the Bytes.
func (b Bytes) Enumerator() ValueEnumerator {
	return &BytesEnumerator{b.b, -1}
}

func (b Bytes) ArrayEnumerator() (OffsetValueEnumerator, bool) {
	return &bytesEnumerator{b: b.b, offset: b.offset, i: -1}, true
}

func newBytesTuple(index int, b byte) Tuple {
	return NewTuple(
		NewIntAttr("@", index),
		NewIntAttr(BytesByteAttr, int(b)),
	)
}

func isBytesTuple(v Value) (index int, b byte, is bool) {
	is = bytesTupleMatcher(func(i int, b2 byte) { index = i; b = b2 }).Match(v)
	return
}

func bytesTupleMatcher(match func(index int, b byte)) TupleMatcher {
	n := 0
	var index int
	var b byte
	check := func() {
		if n == 1 {
			match(index, b)
		}
		n++
	}
	return NewTupleMatcher(
		map[string]Matcher{
			"@":           MatchInt(func(i int) { index = i; check() }),
			BytesByteAttr: MatchInt(func(i int) { b = byte(i); check() }),
		},
		Lit(EmptyTuple),
	)
}

type bytesEnumerator struct {
	b      []byte
	offset int
	i      int
}

func (e *bytesEnumerator) MoveNext() bool {
	if e.i >= len(e.b)-1 {
		return false
	}
	e.i++
	return true
}

func (e *bytesEnumerator) Current() Value {
	return NewNumber(float64(e.b[e.i]))
}

func (e *bytesEnumerator) Offset() int {
	return e.offset + e.i
}

// func stringSet(b Set) Set {
// 	if b, ok := b.(Bytes); ok {
// 		return b
// 	}
// 	if !b.IsTrue() {
// 		return b
// 	}

// 	var result Bytes
// 	matcher := bytesTupleMatcher(func(index int, b byte) {
// 		result = result.with(index, b)
// 	})
// }
