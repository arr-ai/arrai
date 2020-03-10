package rel

import (
	"log"
	"reflect"
	"strings"

	"github.com/arr-ai/frozen"
)

const ArrayItemAttr = "@item"

// func isArrayTuple(v Value) (index int, item Value, is bool) {
// 	is = NewTupleMatcher(
// 		map[string]Matcher{
// 			"@":           MatchInt(func(i int) { index = i }),
// 			ArrayItemAttr: Bind(&item),
// 		},
// 		Lit(EmptyTuple),
// 	).Match(v)
// 	return
// }

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
	if s, ok := s.(Array); ok {
		return s, true
	}
	if i := s.Enumerator(); i.MoveNext() {
		match := arrayTupleMatcher()
		tupleOffset, item, is := match(i.Current())
		if !is {
			return Array{}, false
		}

		middleIndex := s.Count()
		items := make([]Value, 2*middleIndex)
		items[middleIndex] = item
		anchorOffset, minOffset := tupleOffset, tupleOffset
		lowestIndex, highestIndex := middleIndex, middleIndex
		for i.MoveNext() {
			if tupleOffset, item, is = match(i.Current()); !is {
				return Array{}, false
			}
			if tupleOffset < minOffset {
				minOffset = tupleOffset
			}
			sliceIndex := middleIndex - (anchorOffset - tupleOffset)
			items[sliceIndex] = item

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
		panic("Empty string not allowed (should be == None)")
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

// Negate returns {@neg: b}.
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
	match := arrayTupleMatcher()
	log.Print(value)
	if pos, item, ok := match(value); ok {
		if a.offset <= pos && pos < a.offset+len(a.values) {
			return item == a.values[pos-a.offset]
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
	return newSetFromSet(a).With(newArrayTuple(index, item))
}

// With returns the original Array with given value added. Iff the value was
// already present, the original Array is returned.
func (a Array) With(value Value) Set {
	match := arrayTupleMatcher()
	if index, item, ok := match(value); ok {
		return a.with(index, item)
	}
	return newSetFromSet(a).With(value)
}

// Without returns the original Array without the given value. Iff the value
// was already absent, the original Array is returned.
func (a Array) Without(value Value) Set {
	if pos, item, ok := arrayTupleMatcher()(value); ok {
		if i := a.index(pos); i >= 0 && item == a.values[i] {
			if pos == a.offset {
				return Array{values: a.values[1:], offset: a.offset + 1}
			}
			if pos == a.offset+i {
				return Array{values: a.values[:i], offset: a.offset}
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
	result := Set(a)
	for e := a.Enumerator(); e.MoveNext(); {
		value := e.Current()
		if !p(value) {
			result = result.Without(value)
		}
	}
	return result
}

// Call ...
func (a Array) Call(arg Value) Value {
	i := int(arg.(Number).Float64())
	return a.values[i-a.offset]
}

func (a Array) index(pos int) int {
	pos -= a.offset
	if 0 <= pos && pos < len(a.values) {
		return pos
	}
	return -1
}

// Enumerator returns an enumerator over the Values in the Array.
func (a Array) Enumerator() ValueEnumerator {
	return &ArrayValueEnumerator{a: a, i: -1}
}

func (a Array) AsString() (String, bool) {
	return String{}, false
}

func (a Array) ArrayEnumerator() (OffsetValueEnumerator, bool) {
	return &ArrayOffsetValueEnumerator{a: a, i: -1}, true
}

func newArrayTuple(index int, v Value) Tuple {
	return NewTuple(
		NewAttr("@", NewNumber(float64(index))),
		NewAttr(ArrayItemAttr, v),
	)
}

func arrayTupleMatcher() func(v Value) (index int, item Value, matches bool) {
	var index int
	var item Value
	m := NewTupleMatcher(
		map[string]Matcher{
			"@":           MatchInt(func(i int) { index = i }),
			ArrayItemAttr: Let(func(v Value) { item = v }),
		},
		Lit(EmptyTuple),
	)
	return func(v Value) (int, Value, bool) {
		matches := m.Match(v)
		return index, item, matches
	}
}

// ArrayValueEnumerator represents an enumerator over a Array.
type ArrayValueEnumerator struct {
	a Array
	i int
}

// MoveNext moves the enumerator to the next Value.
func (e *ArrayValueEnumerator) MoveNext() bool {
	e.i++
	return e.i < len(e.a.values)
}

// Current returns the enumerator's current Value.
func (e *ArrayValueEnumerator) Current() Value {
	return newArrayTuple(e.a.offset+e.i, e.a.values[e.i])
}

// ArrayOffsetValueEnumerator represents an enumerator over a Array.
type ArrayOffsetValueEnumerator struct {
	a Array
	i int
}

// MoveNext moves the enumerator to the next Value.
func (e *ArrayOffsetValueEnumerator) MoveNext() bool {
	e.i++
	return e.i < len(e.a.values)
}

// Current returns the enumerator's current Value.
func (e *ArrayOffsetValueEnumerator) Current() Value {
	return e.a.values[e.i]
}

// Current returns the offset of the enumerator's current Value.
func (e *ArrayOffsetValueEnumerator) Offset() int {
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

// type arrayEnumerator struct {
// 	values []Value
// 	i      int
// }

// func (e *arrayEnumerator) MoveNext() bool {
// 	if e.i >= len(e.values)-1 {
// 		return false
// 	}
// 	e.i++
// 	return true
// }

// func (e *arrayEnumerator) Current() Value {
// 	return e.values[e.i]
// }

// func stringSet(b Set) Set {
// 	if b, ok := b.(Array); ok {
// 		return b
// 	}
// 	if !b.IsTrue() {
// 		return b
// 	}

// 	var result Array
// 	matcher := arrayTupleMatcher(func(index int, b Value) {
// 		result = result.with(index, b)
// 	})
// }
