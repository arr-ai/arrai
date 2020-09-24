package rel

import (
	"context"
	"fmt"
	"reflect"

	"github.com/arr-ai/hash"
	"github.com/arr-ai/wbnf/parser"
)

// BytesByteAttr is the standard name for the value-attr of a character tuple.
const BytesByteAttr = "@byte"

// Bytes is a set of Values.
type Bytes struct {
	b      []byte
	offset int
}

// NewBytes constructs a byte array as a relation.
func NewBytes(b []byte) Set {
	if len(b) == 0 {
		return None
	}
	return Bytes{b: b}
}

// NewOffsetBytes constructs an offset byte array as a relation.
func NewOffsetBytes(b []byte, offset int) Set {
	if len(b) == 0 {
		return None
	}
	return Bytes{b: b, offset: offset}
}

// TODO: support byte arrays with holes.
func AsBytes(s Set) (Bytes, bool) { //nolint:dupl
	if b, ok := s.(Bytes); ok {
		return b, true
	}
	n := s.Count()
	if n == 0 {
		return Bytes{}, true
	}
	tuples := make(bytesByteTupleArray, 0, n)
	minAt := int(^uint(0) >> 1)
	maxAt := -minAt - 1
	for i := s.Enumerator(); i.MoveNext(); {
		t, is := i.Current().(BytesByteTuple)
		if !is {
			return Bytes{}, false
		}
		if minAt > t.at {
			minAt = t.at
		}
		if maxAt < t.at {
			maxAt = t.at
		}
		tuples = append(tuples, t)
	}
	bytes := make([]byte, maxAt-minAt+1)
	for _, t := range tuples {
		bytes[t.at-minAt] = t.byteval
	}
	return Bytes{b: bytes, offset: minAt}, true
}

type bytesByteTupleArray []BytesByteTuple

func (a bytesByteTupleArray) Len() int {
	return len(a)
}

func (a bytesByteTupleArray) Less(i, j int) bool {
	return a[i].at < a[j].at
}

func (a bytesByteTupleArray) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Bytes returns the bytes of b. The caller must not modify the contents.
func (b Bytes) Bytes() []byte {
	return b.b
}

// Hash computes a hash for a Bytes.
func (b Bytes) Hash(seed uintptr) uintptr {
	// TODO: implement a []byte-friendly hash function.
	return hash.String(string(b.b), seed)
}

// Equal tests two Sets for equality. Any other type returns false.
func (b Bytes) Equal(v interface{}) bool {
	switch x := v.(type) {
	case Bytes:
		if len(b.b) != len(x.b) || b.offset != x.offset {
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
func (b Bytes) Eval(ctx context.Context, _ Scope) (Value, error) {
	return b, nil
}

// Source returns a scanner locating the Bytes's source code.
func (b Bytes) Source() parser.Scanner {
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
func (b Bytes) Export(_ context.Context) interface{} {
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
func (b Bytes) Map(f func(v Value) (Value, error)) (Set, error) {
	result := None
	for e := b.Enumerator().(*BytesEnumerator); e.MoveNext(); {
		v, err := f(e.CurrentBytesByteTuple())
		if err != nil {
			return nil, err
		}
		result = result.With(v)
	}
	return result, nil
}

// Where returns a new Bytes with all the Values satisfying predicate p.
func (b Bytes) Where(p func(v Value) (bool, error)) (Set, error) {
	result := Set(b)
	for e := b.Enumerator().(*BytesEnumerator); e.MoveNext(); {
		value := e.CurrentBytesByteTuple()
		match, err := p(value)
		if err != nil {
			return nil, err
		}
		if !match {
			result = result.Without(value)
		}
	}
	return result, nil
}

func (b Bytes) CallAll(_ context.Context, arg Value) (Set, error) {
	n, ok := arg.(Number)
	if !ok {
		return nil, fmt.Errorf("arg to CallAll must be a number, not %s", ValueTypeAsString(arg))
	}
	i := int(n.Float64()) - b.offset
	if i < 0 || i >= len(b.Bytes()) {
		return None, nil
	}
	return NewSet(NewNumber(float64(string(b.b)[i])))
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

// CurrentBytesByteTuple returns the enumerator'b current Value.
func (e *BytesEnumerator) CurrentBytesByteTuple() BytesByteTuple {
	return NewBytesByteTuple(e.i, e.b[e.i])
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
