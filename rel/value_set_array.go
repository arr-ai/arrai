package rel

import (
	"context"
	"fmt"
	"math"
	"reflect"

	"github.com/arr-ai/wbnf/parser"

	"github.com/arr-ai/arrai/pkg/fu"
)

// Array is an ordered collection of Values.
type Array struct {
	values []Value
	offset int
	count  int
}

// NewArray constructs an array as a relation.
func NewArray(values ...Value) Set {
	return NewOffsetArray(0, values...)
}

// NewOffsetArray constructs an offset array as a relation.
func NewOffsetArray(offset int, values ...Value) Set {
	// Trim holes from both ends.
	for i, v := range values {
		if v != nil {
			if i > 0 {
				offset += i
				values = values[i:]
			}
			break
		}
	}
	for i := len(values) - 1; i >= 0; i-- {
		if values[i] != nil {
			if i < len(values)-1 {
				values = values[:i+1]
			}
			break
		}
	}

	if len(values) == 0 {
		return None
	}

	// Count non-holes.
	n := 0
	for _, v := range values {
		if v != nil {
			n++
		}
	}

	return Array{values: values, offset: offset, count: n}
}

func AsArray(v Value) (Array, bool) {
	switch v := v.(type) {
	case Array:
		return v, true
	case EmptySet:
		return Array{}, true
	}
	return Array{}, false
}

func asArray(values ...Value) Array {
	minIndex := math.MaxInt32
	maxIndex := math.MinInt32
	for _, v := range values {
		t := v.(ArrayItemTuple)
		if t.at < minIndex {
			minIndex = t.at
		}
		if t.at > maxIndex {
			maxIndex = t.at
		}
	}
	items := make([]Value, maxIndex-minIndex+1)

	n := 0
	for _, v := range values {
		t := v.(ArrayItemTuple)
		if items[t.at-minIndex] == nil {
			n++
		}
		items[t.at-minIndex] = t.item
	}
	return Array{
		values: items,
		offset: minIndex,
		count:  n,
	}
}

func (a Array) clone() Array {
	values := make([]Value, len(a.values))
	copy(values, a.values)
	a.values = values
	return a
}

// Values returns the slice of values in the array. Holes in the indices are
// represented by nil elements.
//
// Callers must not reassign elements of the returned slice.
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
		if len(a.values) != len(x.values) || a.offset != x.offset || a.count != x.count {
			return false
		}
		for i, c := range a.values {
			if (c != nil) != (x.values[i] != nil) || c != nil && !c.Equal(x.values[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// String returns a string representation of an Array.
func (a Array) String() string {
	return fu.String(a)
}

// String returns a string representation of an Array.
func (a Array) Format(f fmt.State, verb rune) {
	if a.offset != 0 {
		fu.Fprintf(f, `%d\`, a.offset)
	}
	fu.WriteString(f, "[")
	for i, v := range a.values {
		writeSep(f, i, ", ")
		if v != nil {
			fu.Format(v, f, 'v')
		}
	}
	fu.WriteString(f, "]")
}

// Shift increments the Array's offset
func (a Array) Shift(offset int) Array {
	a.offset += offset
	return a
}

// Eval returns the string.
func (a Array) Eval(ctx context.Context, _ Scope) (Value, error) {
	return a, nil
}

// Source returns a scanner locating the Array's source code.
func (a Array) Source() parser.Scanner {
	return *parser.NewScanner("")
}

var arrayKind = registerKind(208, reflect.TypeOf(Array{}))

// Kind returns a number that is unique for each major kind of Value.
func (a Array) Kind() int {
	return arrayKind
}

// IsTrue returns true if the tuple has attributes.
func (a Array) IsTrue() bool {
	return a.count > 0
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
	if n > len(b.values) {
		n = len(b.values)
	}
	for i, av := range a.values[:n] {
		bv := b.values[i]
		if bv == nil {
			return av != nil
		}
		if av == nil {
			return false
		}
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
func (a Array) Export(ctx context.Context) interface{} {
	result := make([]interface{}, 0, a.Count())
	for _, v := range a.values {
		if v != nil {
			result = append(result, v.Export(ctx))
		} else {
			result = append(result, nil)
		}
	}
	return result
}

func (Array) getSetBuilder() setBuilder {
	return newGenericTypeSetBuilder()
}

func (Array) getBucket() fmt.Stringer {
	return genericType
}

// Count returns the number of elements in the Array.
func (a Array) Count() int {
	return a.count
}

// Has returns true iff the given Value is in the Array.
func (a Array) Has(value Value) bool {
	if t, ok := value.(ArrayItemTuple); ok {
		if a.offset <= t.at && t.at < a.offset+len(a.values) {
			v := a.values[t.at-a.offset]
			return v != nil && v.Equal(t.item)
		}
	}
	return false
}

func (a Array) withItem(index int, item Value) Set {
	b := a
	index -= a.offset
	switch {
	case index < 0:
		b.values = make([]Value, len(a.values)-index)
		copy(b.values[-index:], a.values)
		b.offset += index
		index = 0
	case index >= len(a.values):
		b.values = make([]Value, index+1)
		copy(b.values, a.values)
	case item.Equal(a.values[index]):
		return a
	default:
		b.values = make([]Value, len(a.values))
		copy(b.values, a.values)
	}
	if b.values[index] != nil {
		panic("superimposed array items not supported yet")
	}
	b.values[index] = item
	b.count++
	return b
}

// With returns the original Array with given value added. Iff the value was
// already present, the original Array is returned.
//
// Where value is an ArrayItemTuple, its item is inserted at its at index. If
// a value is already present at the index, this is an error.
func (a Array) With(value Value) Set {
	if t, ok := value.(ArrayItemTuple); ok {
		return a.withItem(t.at, t.item)
	}
	return toUnionSetWithItem(a, value)
}

// Without returns the original Array without the given value. Iff the value
// was already absent, the original Array is returned.
//
// Value is expected to be an ArrayItemTuple; no other values can be present in
// an array, so any other type will be a no-op. Both the value's at (index) and
// item must exactly match an element of the array in order to remove it.
func (a Array) Without(value Value) Set {
	if t, ok := value.(ArrayItemTuple); ok {
		if i := t.at - a.offset; 0 <= i && i < len(a.values) {
			v := a.values[i]
			if v != nil && v.Equal(t.item) {
				if t.at == a.offset {
					return Array{
						values: a.values[1:],
						offset: a.offset + 1,
						count:  a.count - 1,
					}
				}
				if t.at == a.offset+len(a.values)-1 {
					return Array{
						values: a.values[:len(a.values)-1],
						offset: a.offset,
						count:  a.count - 1,
					}
				}
				result := a.clone()
				result.values[i] = nil
				result.count--
				if !result.IsTrue() {
					return None
				}
				return result
			}
		}
	}
	return a
}

// Map maps values per f.
func (a Array) Map(f func(v Value) (Value, error)) (Set, error) {
	b := NewSetBuilder()
	for e := a.Enumerator(); e.MoveNext(); {
		v, err := f(e.Current())
		if err != nil {
			return nil, err
		}
		b.Add(v)
	}
	return b.Finish()
}

// Where returns a new Array with all the Values satisfying predicate p.
func (a Array) Where(p func(v Value) (bool, error)) (Set, error) {
	result := a.clone()
	for i, v := range a.values {
		if v != nil {
			match, err := p(NewArrayItemTuple(a.offset+i, v))
			if err != nil {
				return nil, err
			}
			if !match {
				result.values[i] = nil
				result.count--
			}
		}
	}
	if result.count == 0 {
		return None, nil
	}
	// Trim leading nils.
	for i, v := range result.values {
		if v != nil {
			if i > 0 {
				result.values = result.values[i:]
				result.offset += i
			}
			break
		}
	}
	// Trim trailing nils.
	for i := len(result.values) - 1; i >= 0; i-- {
		if v := result.values[i]; v != nil {
			if i < len(result.values)-1 {
				result.values = result.values[:i+1]
			}
			break
		}
	}
	return result, nil
}

func (a Array) CallAll(_ context.Context, arg Value, b SetBuilder) error {
	if n, ok := arg.(Number); ok {
		if i, is := n.Int(); is {
			i -= a.offset
			if 0 <= i && i < len(a.values) {
				if v := a.values[i]; v != nil {
					b.Add(v)
				}
			}
		}
	}
	return nil
}

func (Array) unionSetSubsetBucket() string {
	return ArrayItemTuple{}.getBucket().String()
}

// Enumerator returns an enumerator over the Values in the Array.
func (a Array) Enumerator() ValueEnumerator {
	return &arrayValueEnumerator{a: a, i: -1}
}

func (a Array) ArrayEnumerator() ValueEnumerator {
	return &arrayItemEnumerator{a.Enumerator().(*arrayValueEnumerator)}
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
	for {
		e.i++
		if e.i < len(e.a.values) && e.a.values[e.i] != nil {
			break
		}
	}
	return e.i < len(e.a.values)
}

// Current returns the enumerator's current Value.
func (e *arrayValueEnumerator) Current() Value {
	return NewArrayItemTuple(e.a.offset+e.i, e.a.values[e.i])
}

// arrayItemEnumerator represents an enumerator over a Array.
type arrayItemEnumerator struct {
	*arrayValueEnumerator
}

// Current returns the enumerator's current Value.
func (e *arrayItemEnumerator) Current() Value {
	return e.a.values[e.i]
}
