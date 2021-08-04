package rel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arr-ai/arrai/pkg/arraictx"
)

func TestNewNumber(t *testing.T) {
	t.Parallel()
	n := NewNumber(0)
	assert.NotNil(t, n)
}

func TestNumberEqual(t *testing.T) {
	t.Parallel()
	a := NewNumber(42)
	b := NewNumber(42)
	c := NewNumber(43)
	assert.True(t, a.Equal(b))
	assert.True(t, b.Equal(a))
	assert.False(t, a.Equal(c))
	assert.False(t, c.Equal(a))
}

func TestNumberHash(t *testing.T) {
	t.Parallel()
	a := NewNumber(42)
	b := NewNumber(42)
	c := NewNumber(43)
	assert.Equal(t, a.Hash(0), b.Hash(0), "%s.Hash(0) vs %s.Hash(0)", a, b)
	assert.NotEqual(t, a.Hash(0), c.Hash(0), "%s.Hash(0) vs %s.Hash(0)", a, c)
	assert.NotEqual(t, a.Hash(0), b.Hash(1), "%s.Hash(0) vs %s.Hash(1)", a, b)
	assert.NotEqual(t, a.Hash(0), c.Hash(1), "%s.Hash(0) vs %s.Hash(1)", a, c)
}

func TestNumberBool(t *testing.T) {
	t.Parallel()
	assert.False(t, NewNumber(0).IsTrue())
	assert.False(t, NewNumber(0.0).IsTrue())
	assert.False(t, NewNumber(-1.0*0.0).IsTrue())
	assert.True(t, NewNumber(-1).IsTrue())
	assert.True(t, NewNumber(0.5).IsTrue())
	assert.True(t, NewNumber(-0.05).IsTrue())
}

func TestNumberLess(t *testing.T) {
	t.Parallel()
	n0 := NewNumber(0)
	n42 := NewNumber(42)
	assert.False(t, n0.Less(n0))
	assert.True(t, n0.Less(n42))
	assert.False(t, n42.Less(n42))
	assert.False(t, n42.Less(n0))
}

func TestNumberExport(t *testing.T) {
	t.Parallel()
	for _, n := range []float64{0, -1, 0.5, -0.05} {
		number := NewNumber(n)
		assert.Equal(t, n, number.Export(arraictx.InitRunCtx(context.Background())), "%s.Export()", number)
	}
}

func TestNumberString(t *testing.T) {
	t.Parallel()
	for _, s := range []struct {
		repr   string
		number Number
	}{
		{"0", NewNumber(0)},
		{"-1", NewNumber(-1)},
		{"0.5", NewNumber(0.5)},
		{"-0.05", NewNumber(-0.05)},
	} {
		assert.Equal(t, s.repr, s.number.String(), "%s.String()", s.number)
	}
}
