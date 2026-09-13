package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamesInterned(t *testing.T) {
	t.Parallel()

	a := NewNames("b", "a", "a")
	b := NewNames("a", "b")
	assert.Equal(t, a, b)
	assert.True(t, a.Equal(b))
	assert.Equal(t, []string{"a", "b"}, a.OrderedNames())
	assert.Equal(t, "|a, b|", a.String())

	assert.Equal(t, EmptyNames, NewNames())
	assert.True(t, Names{}.Equal(EmptyNames))

	tup := NewTuple(NewAttr("b", NewNumber(1)), NewAttr("a", NewNumber(2)))
	assert.Equal(t, a, tup.Names())
	assert.Equal(t, arrayItemNames, NewArrayItemTuple(0, None).Names())
	assert.Equal(t, arrayItemNames, NewNames("@", ArrayItemAttr))
}

func TestNamesAlgebra(t *testing.T) {
	t.Parallel()

	abc := NewNames("a", "b", "c")
	bcd := NewNames("b", "c", "d")
	bc := NewNames("b", "c")
	assert.Equal(t, bc, abc.Intersect(bcd))
	assert.Equal(t, NewNames("a"), abc.Minus(bcd))
	assert.True(t, bc.IsSubsetOf(abc))
	assert.False(t, abc.IsSubsetOf(bc))
	assert.True(t, EmptyNames.IsSubsetOf(abc))
	assert.Equal(t, abc, abc.With("b"))
	assert.Equal(t, abc, abc.Without("z"))
	assert.Equal(t, NewNames("a", "b", "c", "d"), abc.With("d"))
	assert.Equal(t, NewNames("a", "c"), abc.Without("b"))
}
