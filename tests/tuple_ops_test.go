package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/rel"
)

// TestTrivialMerge test rel.Merge({}, {}).
func TestTrivialMerge(t *testing.T) {
	a := rel.EmptyTuple
	b := rel.EmptyTuple
	assert.True(t, rel.EmptyTuple.Equal(rel.Merge(a, b)))
}

// TestOneSidedMerge tests rel.Merge({}, {a:42, b:43}) and vice-versa.
func TestOneSidedMerge(t *testing.T) {
	a := rel.EmptyTuple
	b := rel.NewTuple(
		rel.Attr{"a", rel.NewNumber(42)},
		rel.Attr{"b", rel.NewNumber(43)},
	)
	assert.True(t, b.Equal(rel.Merge(a, b)))
	assert.True(t, b.Equal(rel.Merge(b, a)))
}

// TestEqualMerge tests rel.Merge({a:42, b:43}, {a:42, b:43}).
func TestEqualMerge(t *testing.T) {
	a := rel.NewTuple(
		rel.Attr{"a", rel.NewNumber(42)},
		rel.Attr{"b", rel.NewNumber(43)},
	)
	b := rel.NewTuple(
		rel.Attr{"a", rel.NewNumber(42)},
		rel.Attr{"b", rel.NewNumber(43)},
	)
	assert.True(t, a.Equal(rel.Merge(a, b)))
}

// TestMixedMerge tests rel.Merge({a:42, b:43}, {b:43, c:44}).
func TestMixedMerge(t *testing.T) {
	a := rel.NewTuple(
		rel.Attr{"a", rel.NewNumber(42)},
		rel.Attr{"b", rel.NewNumber(43)},
	)
	b := rel.NewTuple(
		rel.Attr{"b", rel.NewNumber(43)},
		rel.Attr{"c", rel.NewNumber(44)},
	)
	c := rel.NewTuple(
		rel.Attr{"a", rel.NewNumber(42)},
		rel.Attr{"b", rel.NewNumber(43)},
		rel.Attr{"c", rel.NewNumber(44)},
	)
	assert.True(t, c.Equal(rel.Merge(a, b)))
}

// TestFailedMerge tests rel.Merge({a:42, b:43}, {b:432, c:44}).
func TestFailedMerge(t *testing.T) {
	a := rel.NewTuple(
		rel.Attr{"a", rel.NewNumber(42)},
		rel.Attr{"b", rel.NewNumber(43)},
	)
	b := rel.NewTuple(
		rel.Attr{"b", rel.NewNumber(432)},
		rel.Attr{"c", rel.NewNumber(44)},
	)
	assert.Nil(t, rel.Merge(a, b))
}
