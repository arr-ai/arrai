package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
)

// TestTrivialMerge test Merge({}, {}).
func TestTrivialMerge(t *testing.T) {
	a := EmptyTuple
	b := EmptyTuple
	assert.True(t, EmptyTuple.Equal(Merge(a, b)))
}

// TestOneSidedMerge tests Merge({}, {a:42, b:43}) and vice-versa.
func TestOneSidedMerge(t *testing.T) {
	a := EmptyTuple
	b := NewTuple(
		Attr{"a", NewNumber(42)},
		Attr{"b", NewNumber(43)},
	)
	assert.True(t, b.Equal(Merge(a, b)))
	assert.True(t, b.Equal(Merge(b, a)))
}

// TestEqualMerge tests Merge({a:42, b:43}, {a:42, b:43}).
func TestEqualMerge(t *testing.T) {
	a := NewTuple(
		Attr{"a", NewNumber(42)},
		Attr{"b", NewNumber(43)},
	)
	b := NewTuple(
		Attr{"a", NewNumber(42)},
		Attr{"b", NewNumber(43)},
	)
	assert.True(t, a.Equal(Merge(a, b)))
}

// TestMixedMerge tests Merge({a:42, b:43}, {b:43, c:44}).
func TestMixedMerge(t *testing.T) {
	a := NewTuple(
		Attr{"a", NewNumber(42)},
		Attr{"b", NewNumber(43)},
	)
	b := NewTuple(
		Attr{"b", NewNumber(43)},
		Attr{"c", NewNumber(44)},
	)
	c := NewTuple(
		Attr{"a", NewNumber(42)},
		Attr{"b", NewNumber(43)},
		Attr{"c", NewNumber(44)},
	)
	assert.True(t, c.Equal(Merge(a, b)))
}

// TestFailedMerge tests Merge({a:42, b:43}, {b:432, c:44}).
func TestFailedMerge(t *testing.T) {
	a := NewTuple(
		Attr{"a", NewNumber(42)},
		Attr{"b", NewNumber(43)},
	)
	b := NewTuple(
		Attr{"b", NewNumber(432)},
		Attr{"c", NewNumber(44)},
	)
	assert.Nil(t, Merge(a, b))
}
