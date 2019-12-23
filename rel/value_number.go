package rel

import (
	"encoding/binary"
	"strconv"

	"github.com/OneOfOne/xxhash"
)

// Number is a number.
type Number struct {
	number float64
}

// NewNumber returns a Number for the given number.
func NewNumber(number float64) *Number {
	return &Number{number}
}

// Float64 returns the value of the number.
func (n *Number) Float64() float64 {
	return n.number
}

// Hash computes a hash for a Number.
func (n *Number) Hash(seed uint32) uint32 {
	xx := xxhash.NewS32(seed)
	if err := binary.Write(xx, binary.LittleEndian, n.number); err != nil {
		panic(err)
	}
	return xx.Sum32()
}

// Equal tests two Values for equality. Any other type returns false.
func (n *Number) Equal(v interface{}) bool {
	if b, ok := v.(*Number); ok {
		return n.number == b.number
	}
	return false
}

// String returns a string representation of a Number.
func (n *Number) String() string {
	return strconv.FormatFloat(n.number, 'G', -1, 64)
}

// Eval returns the number.
func (n *Number) Eval(local, global *Scope) (Value, error) {
	return n, nil
}

// Kind returns a number that is unique for each major kind of Value.
func (n *Number) Kind() int {
	return 100
}

// Bool returns true iff the tuple has attributes.
func (n *Number) Bool() bool {
	return n.number != 0
}

// Less returns true iff v is not a number or n.number < v.number.
func (n *Number) Less(v Value) bool {
	if n.Kind() != v.Kind() {
		return n.Kind() < v.Kind()
	}
	return n.number < v.(*Number).number
}

// Negate returns -n.
func (n *Number) Negate() Value {
	if !n.Bool() {
		return n
	}
	return NewNumber(-n.number)
}

// Export exports a Number.
func (n *Number) Export() interface{} {
	return n.number
}
