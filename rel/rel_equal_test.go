package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustRel(t *testing.T, tuples ...Tuple) Set {
	t.Helper()
	b := NewSetBuilder()
	for _, tup := range tuples {
		b.Add(tup)
	}
	s, err := b.Finish()
	require.NoError(t, err)
	return s
}

// Relations compare equal by content regardless of the attribute layout
// they were built with, and the fast path (same layout) agrees with the
// canonicalising path.
func TestRelationEqualLayouts(t *testing.T) {
	t.Parallel()

	ab := func(a, b float64) Tuple { return NewTuple(NewAttr("a", NewNumber(a)), NewAttr("b", NewNumber(b))) }
	ba := func(a, b float64) Tuple { return NewTuple(NewAttr("b", NewNumber(b)), NewAttr("a", NewNumber(a))) }

	r1 := mustRel(t, ab(1, 2), ab(3, 4))
	r2 := mustRel(t, ab(3, 4), ab(1, 2))
	r3 := mustRel(t, ba(1, 2), ba(3, 4))
	assert.True(t, r1.Equal(r2), "same layout, different insertion order")
	assert.True(t, r1.Equal(r3), "different attribute order")
	assert.True(t, r3.Equal(r1))
	assert.Equal(t, r1.Hash128(), r3.Hash128())

	assert.False(t, r1.Equal(mustRel(t, ab(1, 2))), "different count")
	assert.False(t, r1.Equal(mustRel(t, ab(1, 2), ab(3, 5))), "different row")
	assert.False(t, r1.Equal(mustRel(t,
		NewTuple(NewAttr("a", NewNumber(1)), NewAttr("c", NewNumber(2))),
		NewTuple(NewAttr("a", NewNumber(3)), NewAttr("c", NewNumber(4))))), "different attrs")

	// Projection-derived relations (non-identity layout) still compare.
	proj, err := r1.(Relation).Where(func(v Value) (bool, error) { return true, nil })
	require.NoError(t, err)
	assert.True(t, proj.Equal(r3))
}

// Failed pattern matches carry lazily formatted messages that read the same
// as before.
func TestLazyPatternErrors(t *testing.T) {
	t.Parallel()

	err := lazyErrorf("couldn't find %s in tuple %s", "x", NewTuple(NewAttr("y", NewNumber(1))))
	assert.EqualError(t, err, "couldn't find x in tuple (y: 1)")
}
