package rel

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/arr-ai/frozen"
)

// Array is an ordered collection of Values.
type Array struct {
	values []Value
	offset int
}

// NewArray constructs an array as a relation.
func NewArray(values ...Value) Set {
	return NewOffsetArray(0, values...)
}

// NewArray constructs an array as a relation.
func NewOffsetArray(offset int, values ...Value) Set {
	if len(values) == 0 {
		return None
	}
	return Array{values: values, offset: offset}
}

func AsArray(s Set) (Array, bool) {
	switch s := s.(type) {
	case Array:
		return s, true
	case Set:
		return Array{}, !s.IsTrue()
	}
	return Array{}, false
}

func asArray(s Set) (Array, bool) {
	if i := s.Enumerator(); i.MoveNext() {
		t, is := i.Current().(ArrayItemTuple)
		if !is {
			return Array{}, false
		}

		middleIndex := s.Count()
		items := make([]Value, 2*middleIndex)
		items[middleIndex] = t.item
		anchorOffset, minOffset := t.at, t.at
		lowestIndex, highestIndex := middleIndex, middleIndex
		for i.MoveNext() {
			if t, is = i.Current().(ArrayItemTuple); !is {
				return Array{}, false
			}
			if t.at < minOffset {
				minOffset = t.at
			}
			sliceIndex := middleIndex - (anchorOffset - t.at)
			items[sliceIndex] = t.item

			if sliceIndex < lowestIndex {
				lowestIndex = sliceIndex
			}

			if sliceIndex > highestIndex {
				highestIndex = sliceIndex
			}
		}

		return Array{values: items[lowestIndex : highestIndex+1], offset: minOffset}, true
	}
	return Array{}, true
}

// Values returns the slice of values in the array.
func (a Array) Values() []Value {
	return a.values
}

// Hash computes a hash for a Array.
func (a Array) Hash(seed uintptr) uintptr {
	h := seed
	for e := a.Enumerator(); e.MoveNext(); {
		h ^= e.Current().Hash(seed)
	}
	return h
}

// Equal tests two Sets for equality. Any other type returns false.
func (a Array) Equal(v interface{}) bool {
	switch x := v.(type) {
	case Array:
		if len(a.values) != len(x.values) {
			return false
		}
		for i, c := range a.values {
			if !c.Equal(x.values[i]) {
				return false
			}
		}
		return true
	case Set:
		return newSetFromSet(a).Equal(x)
	}
	return false
}

// String returns a string representation of an Array.
func (a Array) String() string {
	var sb strings.Builder
	if a.offset != 0 {
		fmt.Fprintf(&sb, `%d\`, a.offset)
	}
	sb.WriteRune('[')
	for i, v := range a.values {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.String())
	}
	sb.WriteRune(']')
	return sb.String()
}

// Eval returns the string.
func (a Array) Eval(_ Scope) (Value, error) {
	return a, nil
}

var arrayKind = registerKind(208, reflect.TypeOf(Array{}))

// Kind returns a number that is unique for each major kind of Value.
func (a Array) Kind() int {
	return arrayKind
}

// Bool returns true iff the tuple has attributes.
func (a Array) IsTrue() bool {
	if len(a.values) == 0 {
		panic("Empty array not allowed (should be == None)")
	}
	return true
}

// Less returns true iff v is not a number or tuple, or v is a tuple and t
// precedes v in a lexicographical comparison of their name/value pairs.
func (a Array) Less(v Value) bool {
	if a.Kind() != v.Kind() {
		return a.Kind() < v.Kind()
	}
	b := v.(Array)
	if a.offset != b.offset {
		return a.offset < b.offset
	}
	n := len(a.values)
	if n < len(b.values) {
		n = len(b.values)
	}
	for i, av := range a.values[:n] {
		bv := b.values[i]
		if av.Less(bv) {
			return true
		}
		if bv.Less(av) {
			return false
		}
	}
	return len(a.values) < len(b.values)
}

// Negate returns {@neg: a}.
func (a Array) Negate() Value {
	return NewTuple(NewAttr(negateTag, a))
}

// Export exports an Array as a slice.
func (a Array) Export() interface{} {
	result := make([]interface{}, 0, a.Count())
	for _, v := range a.values {
		result = append(result, v.Export())
	}
	return result
}

// Count returns the number of elements in the Array.
func (a Array) Count() int {
	return len(a.values)
}

// Has returns true iff the given Value is in the Array.
func (a Array) Has(value Value) bool {
	if t, ok := value.(ArrayItemTuple); ok {
		if a.offset <= t.at && t.at < a.offset+len(a.values) {
			return t.item == a.values[t.at-a.offset]
		}
	}
	return false
}

func (a Array) with(index int, item Value) Set {
	if a.index(index) == len(a.values) {
		return Array{values: append(a.values, item), offset: a.offset}
	} else if index == a.offset-1 {
		return Array{
			values: append(append(make([]Value, 0, 1+len(a.values)), item), a.values...),
			offset: a.offset - 1,
		}
	}
	return newSetFromSet(a).With(NewArrayItemTuple(index, item))
}

// With returns the original Array with given value added. Iff the value was
// already present, the original Array is returned.
func (a Array) With(value Value) Set {
	if t, ok := value.(ArrayItemTuple); ok {
		return a.with(t.at, t.item)
	}
	return newSetFromSet(a).With(value)
}

// Without returns the original Array without the given value. Iff the value
// was already absent, the original Array is returned.
func (a Array) Without(value Value) Set {
	if t, ok := value.(ArrayItemTuple); ok {
		if i := a.index(t.at); i >= 0 && t.item == a.values[i] {
			if t.at == a.offset {
				return Array{values: a.values[1:], offset: a.offset + 1}
			}
			if t.at == a.offset+len(a.values)-1 {
				return Array{values: a.values[:len(a.values)-1], offset: a.offset}
			}
			return newSetFromSet(a).Without(value)
		}
	}
	return a
}

// Map maps values per f.
func (a Array) Map(f func(v Value) Value) Set {
	result := NewSet()
	for e := a.Enumerator(); e.MoveNext(); {
		result = result.With(f(e.Current()))
	}
	return result
}

// Where returns a new Array with all the Values satisfying predicate p.
func (a Array) Where(p func(v Value) bool) Set {
	values := make([]Value, 0, a.Count())
	for e := a.Enumerator(); e.MoveNext(); {
		value := e.Current()
		if p(value) {
			values = append(values, value)
		}
	}
	return NewSet(values...)
}

// Call ...
func (a Array) Call(arg Value) Value {
	i := int(arg.(Number).Float64())
	return a.values[i-a.offset]
}

func (a Array) index(pos int) int {
	pos -= a.offset
	if 0 <= pos && pos <= len(a.values) {
		return pos
	}
	return -1
}

// Enumerator returns an enumerator over the Values in the Array.
func (a Array) Enumerator() ValueEnumerator {
	return &arrayValueEnumerator{a: a, i: -1}
}

func (a Array) ArrayEnumerator() (OffsetValueEnumerator, bool) {
	return &arrayOffsetValueEnumerator{arrayValueEnumerator{a: a, i: -1}}, true
}

// arrayValueEnumerator represents an enumerator over a Array.
type arrayValueEnumerator struct {
	a Array
	i int
}

// MoveNext moves the enumerator to the next Value.
func (e *arrayValueEnumerator) MoveNext() bool {
	if e.i >= len(e.a.values)-1 {
		return false
	}
	e.i++
	return true
}

// Current returns the enumerator's current Value.
func (e *arrayValueEnumerator) Current() Value {
	return NewArrayItemTuple(e.a.offset+e.i, e.a.values[e.i])
}

// arrayOffsetValueEnumerator represents an enumerator over a Array.
type arrayOffsetValueEnumerator struct {
	arrayValueEnumerator
}

// Current returns the enumerator's current Value.
func (e *arrayOffsetValueEnumerator) Current() Value {
	return e.a.values[e.i]
}

// Current returns the offset of the enumerator's current Value.
func (e *arrayOffsetValueEnumerator) Offset() int {
	return e.a.offset + e.i
}

type arrayEnumerator struct {
	i frozen.Iterator
	t Tuple
}

func (e *arrayEnumerator) MoveNext() bool {
	if e.i.Next() {
		e.t = e.i.Value().(Tuple)
		return true
	}
	return false
}

func (e *arrayEnumerator) Current() Value {
	return e.t.MustGet(ArrayItemAttr)
}

func (e *arrayEnumerator) Offset() int {
	return int(e.t.MustGet("@").(Number).Float64())
}
