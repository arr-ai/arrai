package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShapeInterning(t *testing.T) {
	t.Parallel()

	a := NewTuple(NewAttr("x", NewNumber(1)), NewAttr("y", NewNumber(2))).(*GenericTuple)
	b := NewTuple(NewAttr("y", NewNumber(3)), NewAttr("x", NewNumber(4))).(*GenericTuple)
	c := NewTuple(NewAttr("x", NewNumber(1))).(*GenericTuple)
	assert.Same(t, a.shape, b.shape, "same attribute set, any order, one shape")
	assert.NotSame(t, a.shape, c.shape)
	assert.Equal(t, []string{"x", "y"}, TupleOrderedNames(a))

	// Transitions land on the interned shape and are memoised.
	d := c.With("y", NewNumber(9)).(*GenericTuple)
	assert.Same(t, a.shape, d.shape)
	e := a.Without("y").(*GenericTuple)
	assert.Same(t, c.shape, e.shape)
	s1, at1 := c.shape.With("y")
	s2, at2 := c.shape.With("y")
	assert.Same(t, s1, s2)
	assert.Equal(t, at1, at2)
	assert.Equal(t, 1, at1)
}

func TestShapedTupleOperations(t *testing.T) {
	t.Parallel()

	tup := NewTuple(NewAttr("b", NewNumber(2)), NewAttr("a", NewNumber(1)), NewAttr("c", NewNumber(3)))
	v, ok := tup.Get("b")
	require.True(t, ok)
	assert.Equal(t, NewNumber(2), v)
	_, ok = tup.Get("z")
	assert.False(t, ok)
	assert.Equal(t, 3, tup.Count())
	assert.Equal(t, "(a: 1, b: 2, c: 3)", tup.String())

	// With replaces in place or inserts in order; the original is untouched.
	w := tup.With("b", NewNumber(20)).With("aa", NewNumber(0))
	assert.Equal(t, "(a: 1, aa: 0, b: 20, c: 3)", w.String())
	assert.Equal(t, "(a: 1, b: 2, c: 3)", tup.String())
	assert.Equal(t, "(a: 1, c: 3)", tup.Without("b").String())
	assert.Same(t, tup, tup.Without("zzz"), "removing an absent attr is a no-op")

	// View attributes strip their counterpart.
	viewed := tup.With("&b", NewNumber(0))
	_, has := viewed.Get("b")
	assert.False(t, has)
	_, has = viewed.Get("&b")
	assert.True(t, has)
	_, has = viewed.With("b", NewNumber(1)).Get("&b")
	assert.False(t, has)

	// Enumeration is in name order and allocation-light.
	var names []string
	for e := tup.Enumerator(); e.MoveNext(); {
		n, _ := e.Current()
		names = append(names, n)
	}
	assert.Equal(t, []string{"a", "b", "c"}, names)

	// Equality: same shape positional; different attribute sets unequal;
	// generic vs specialised representations still agree.
	assert.True(t, tup.Equal(NewTuple(NewAttr("c", NewNumber(3)), NewAttr("a", NewNumber(1)), NewAttr("b", NewNumber(2)))))
	assert.False(t, tup.Equal(tup.Without("c")))
	assert.False(t, tup.Equal(tup.With("b", NewNumber(0))))
	gen := newGenericTuple(NewAttr("@", NewNumber(1)), NewAttr(ArrayItemAttr, NewNumber(2)))
	assert.True(t, gen.Equal(NewArrayItemTuple(1, NewNumber(2))))
	assert.Equal(t, gen.Hash128(), NewArrayItemTuple(1, NewNumber(2)).Hash128())

	// Builder: last Put wins, and canonicalisation still applies.
	var b TupleBuilder
	b.Put("k", NewNumber(1))
	b.Put("k", NewNumber(2))
	assert.Equal(t, "(k: 2)", b.Finish().String())
	var b2 TupleBuilder
	b2.Put(ArrayItemAttr, NewString([]rune("x")))
	b2.Put("@", NewNumber(0))
	assert.IsType(t, ArrayItemTuple{}, b2.Finish())

	mapped, err := tup.Map(func(v Value) (Value, error) { return NewNumber(v.(Number).Float64() * 10), nil })
	require.NoError(t, err)
	assert.Equal(t, "(a: 10, b: 20, c: 30)", mapped.String())
	assert.Equal(t, "(a: 1, c: 3)", tup.Project(NewNames("a", "c")).String())
	assert.Nil(t, tup.Project(NewNames("a", "zz")))
}

func TestRelationRowsInflateWithoutCopy(t *testing.T) {
	t.Parallel()

	rel := mustRel(t,
		NewTuple(NewAttr("b", NewNumber(2)), NewAttr("a", NewNumber(1))),
		NewTuple(NewAttr("a", NewNumber(3)), NewAttr("b", NewNumber(4))),
	).(Relation)
	assert.True(t, rel.direct, "builder rows are in shape order")
	seen := 0
	for e := rel.Enumerator(); e.MoveNext(); {
		g := e.Current().(*GenericTuple)
		assert.Same(t, rel.shape, g.shape)
		seen++
	}
	assert.Equal(t, 2, seen)

	// A derived relation with another layout still inflates correctly.
	where, err := rel.Where(func(v Value) (bool, error) { return v.(Tuple).MustGet("a").Equal(NewNumber(3)), nil })
	require.NoError(t, err)
	assert.True(t, where.Equal(mustRel(t, NewTuple(NewAttr("a", NewNumber(3)), NewAttr("b", NewNumber(4))))))
	assert.Equal(t, 1, where.Count())
}

// Array memoises its hash, which depends on both its values and its offset.
// Every operation that derives a differing array must not serve the
// original's cached hash.
func TestArrayHashCacheInvalidation(t *testing.T) {
	t.Parallel()

	a := NewArray(NewNumber(1), NewNumber(2), NewNumber(3))
	h := a.Hash128() // populate the cache before deriving
	assert.Equal(t, h, a.Hash128(), "stable")
	assert.Equal(t, h, NewArray(NewNumber(1), NewNumber(2), NewNumber(3)).Hash128(), "equal arrays agree")

	shifted := a.(Array).Shift(5)
	assert.NotEqual(t, h, shifted.Hash128(), "offset participates in the hash")
	assert.Equal(t, NewOffsetArray(5, NewNumber(1), NewNumber(2), NewNumber(3)).Hash128(), shifted.Hash128())

	withItem := a.With(NewArrayItemTuple(3, NewNumber(4)))
	assert.NotEqual(t, h, withItem.Hash128())
	assert.Equal(t, NewArray(NewNumber(1), NewNumber(2), NewNumber(3), NewNumber(4)).Hash128(), withItem.Hash128())

	filtered, err := a.Where(func(v Value) (bool, error) {
		return !v.(ArrayItemTuple).item.Equal(NewNumber(2)), nil
	})
	require.NoError(t, err)
	assert.NotEqual(t, h, filtered.Hash128())

	without := a.Without(NewArrayItemTuple(0, NewNumber(1)))
	assert.NotEqual(t, h, without.Hash128())
	assert.Equal(t, NewOffsetArray(1, NewNumber(2), NewNumber(3)).Hash128(), without.Hash128())

	// Sets of arrays rely on the hash agreeing with equality.
	s, err := NewSet(a, shifted, withItem, filtered, without,
		NewArray(NewNumber(1), NewNumber(2), NewNumber(3)))
	require.NoError(t, err)
	assert.Equal(t, 5, s.Count(), "the duplicate collapses, the rest do not")
	for _, v := range []Value{a, shifted, withItem, filtered, without} {
		assert.True(t, s.Has(v))
	}
}
