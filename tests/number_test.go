package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arr-ai/arrai/rel"
)

// TestNewNumber tests NewNumber.
func TestNewNumber(t *testing.T) {
	n := rel.NewNumber(0)
	assert.NotNil(t, n)
}

// TestNumberEqual tests rel.Number.Equal.
func TestNumberEqual(t *testing.T) {
	a := rel.NewNumber(42)
	b := rel.NewNumber(42)
	c := rel.NewNumber(43)
	assert.True(t, a.Equal(b))
	assert.True(t, b.Equal(a))
	assert.False(t, a.Equal(c))
	assert.False(t, c.Equal(a))
}

// TestNumberHash tests rel.Number.Hash.
func TestNumberHash(t *testing.T) {
	a := rel.NewNumber(42)
	b := rel.NewNumber(42)
	c := rel.NewNumber(43)
	assert.Equal(t, a.Hash(0), b.Hash(0), "%s.Hash(0) vs %s.Hash(0)", a, b)
	assert.NotEqual(t, a.Hash(0), c.Hash(0), "%s.Hash(0) vs %s.Hash(0)", a, c)
	assert.NotEqual(t, a.Hash(0), b.Hash(1), "%s.Hash(0) vs %s.Hash(1)", a, b)
	assert.NotEqual(t, a.Hash(0), c.Hash(1), "%s.Hash(0) vs %s.Hash(1)", a, c)
}

// TestNumberBool tests rel.Number.Bool.
func TestNumberBool(t *testing.T) {
	assert.False(t, rel.NewNumber(0).Bool())
	assert.False(t, rel.NewNumber(0.0).Bool())
	assert.False(t, rel.NewNumber(-0.0).Bool())
	assert.True(t, rel.NewNumber(-1).Bool())
	assert.True(t, rel.NewNumber(0.5).Bool())
	assert.True(t, rel.NewNumber(-0.05).Bool())
}

// TestNumberLess tests rel.Number.Less.
func TestNumberLess(t *testing.T) {
	n0 := rel.NewNumber(0)
	n42 := rel.NewNumber(42)
	assert.False(t, n0.Less(n0))
	assert.True(t, n0.Less(n42))
	assert.False(t, n42.Less(n42))
	assert.False(t, n42.Less(n0))
}

// TestNumberExport tests rel.Number.Export.
func TestNumberExport(t *testing.T) {
	for _, n := range []float64{0, -1, 0.5, -0.05} {
		number := rel.NewNumber(n)
		assert.Equal(t, n, number.Export(), "%s.Export()", number)
	}
}

// TestNumberString tests rel.Number.String.
func TestNumberString(t *testing.T) {
	for _, s := range []struct {
		repr   string
		number *rel.Number
	}{
		{"0", rel.NewNumber(0)},
		{"-1", rel.NewNumber(-1)},
		{"0.5", rel.NewNumber(0.5)},
		{"-0.05", rel.NewNumber(-0.05)},
	} {
		assert.Equal(t, s.repr, s.number.String(), "%s.String()", s.number)
	}
}
