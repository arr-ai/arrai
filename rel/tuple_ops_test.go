package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrivialMerge(t *testing.T) {
	t.Parallel()
	a := EmptyTuple
	b := EmptyTuple
	assert.True(t, EmptyTuple.Equal(Merge(a, b)))
}

func TestOneSidedMerge(t *testing.T) {
	t.Parallel()
	a := EmptyTuple
	b := NewTuple(
		Attr{"a", NewNumber(42)},
		Attr{"b", NewNumber(43)},
	)
	assert.True(t, b.Equal(Merge(a, b)))
	assert.True(t, b.Equal(Merge(b, a)))
}

func TestEqualMerge(t *testing.T) {
	t.Parallel()
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

func TestMixedMerge(t *testing.T) {
	t.Parallel()
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

func TestFailedMerge(t *testing.T) {
	t.Parallel()
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
