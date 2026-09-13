package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Every tuple representation hashes with one formula, so equal tuples hash
// equal regardless of which concrete type holds them.
func TestHashTupleRepresentationsAgree(t *testing.T) {
	t.Parallel()

	// Build GenericTuples directly: NewTuple would canonicalise (@, @item)
	// and friends into the specialised kinds, which is what this test checks
	// against.
	generic := func(attrs ...Attr) Value { return newGenericTuple(attrs...) }
	cases := []struct {
		name string
		a, b Value
	}{
		{"array item", NewArrayItemTuple(3, NewString([]rune("x"))),
			generic(NewAttr("@", NewNumber(3)), NewAttr(ArrayItemAttr, NewString([]rune("x"))))},
		{"dict entry", NewDictEntryTuple(NewNumber(1), NewNumber(2)),
			generic(NewAttr("@", NewNumber(1)), NewAttr(DictValueAttr, NewNumber(2)))},
		{"string char", NewStringCharTuple(2, 'a'),
			generic(NewAttr("@", NewNumber(2)), NewAttr(StringCharAttr, NewNumber('a')))},
		{"bytes byte", NewBytesByteTuple(0, 7),
			generic(NewAttr("@", NewNumber(0)), NewAttr(BytesByteAttr, NewNumber(7)))},
		{"attr order", generic(NewAttr("a", NewNumber(1)), NewAttr("b", NewNumber(2))),
			generic(NewAttr("b", NewNumber(2)), NewAttr("a", NewNumber(1)))},
	}
	for _, c := range cases {
		assert.Equal(t, c.a.Hash(), c.b.Hash(), c.name)
	}
	assert.NotEqual(t,
		generic(NewAttr("a", NewNumber(1))).Hash(),
		generic(NewAttr("b", NewNumber(1))).Hash(), "different attr names")
	assert.NotEqual(t,
		generic(NewAttr("a", NewNumber(1))).Hash(),
		generic(NewAttr("a", NewNumber(2))).Hash(), "different attr values")
}

func TestHashCollectionsAreOrderIndependent(t *testing.T) {
	t.Parallel()
	mustSet := func(vs ...Value) Set {
		s, err := NewSet(vs...)
		require.NoError(t, err)
		return s
	}
	s1 := mustSet(NewNumber(1), NewNumber(2), NewNumber(3))
	s2 := mustSet(NewNumber(3), NewNumber(2), NewNumber(1))
	assert.Equal(t, s1.Hash(), s2.Hash())

	d1, err := NewDict(false, NewDictEntryTuple(NewNumber(1), NewNumber(2)), NewDictEntryTuple(NewNumber(3), NewNumber(4)))
	require.NoError(t, err)
	d2, err := NewDict(false, NewDictEntryTuple(NewNumber(3), NewNumber(4)), NewDictEntryTuple(NewNumber(1), NewNumber(2)))
	require.NoError(t, err)
	assert.Equal(t, d1.Hash(), d2.Hash())
	d3, err := NewDict(false, NewDictEntryTuple(NewNumber(1), NewNumber(2)), NewDictEntryTuple(NewNumber(3), NewNumber(5)))
	require.NoError(t, err)
	assert.NotEqual(t, d1.Hash(), d3.Hash(), "values participate")

	outer1 := NewTuple(NewAttr("d", d1))
	outer2 := NewTuple(NewAttr("d", d2))
	assert.Equal(t, outer1.Hash(), outer2.Hash())

	r1 := mustSet(NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewNumber(2))))
	r2 := mustSet(NewTuple(NewAttr("b", NewNumber(2)), NewAttr("a", NewNumber(1))))
	assert.Equal(t, r1.Hash(), r2.Hash())
}

func TestHashStringsAndBytes(t *testing.T) {
	t.Parallel()
	assert.Equal(t, NewString([]rune("héllo")).Hash(), NewString([]rune("héllo")).Hash())
	assert.NotEqual(t, NewString([]rune("héllo")).Hash(), NewString([]rune("hello")).Hash())
	assert.NotEqual(t, NewString([]rune("1")).Hash(), NewNumber(1).Hash())
	assert.NotEqual(t, NewBytes([]byte("a")).Hash(), NewString([]rune("a")).Hash())
}

func TestHashPlusZeroEqualsMinusZero(t *testing.T) {
	t.Parallel()
	pos := NewNumber(0.0)
	neg := NewNumber(-1.0 * 0.0)
	assert.True(t, pos.Equal(neg))
	assert.Equal(t, pos.Hash(), neg.Hash())
}

func TestHashEmptySetDiffersFromEmptyTuple(t *testing.T) {
	t.Parallel()
	assert.NotEqual(t, EmptySet{}.Hash(), EmptyTuple.Hash(), "{} and () must not share a hash")
	assert.Equal(t, EmptySet{}.Hash(), Dict{}.Hash(), "empty Dict is {}")
	assert.NotEqual(t, True.Hash(), EmptyTuple.Hash(), "{()} is a set wrap of (), not ()")
}

func TestHashRelationPairingDoesNotCollapse(t *testing.T) {
	t.Parallel()
	alice, bob := NewString([]rune("Alice")), NewString([]rune("Bob"))
	a30 := NewTuple(NewAttr("name", alice), NewAttr("age", NewNumber(30)))
	b40 := NewTuple(NewAttr("name", bob), NewAttr("age", NewNumber(40)))
	a40 := NewTuple(NewAttr("name", alice), NewAttr("age", NewNumber(40)))
	b30 := NewTuple(NewAttr("name", bob), NewAttr("age", NewNumber(30)))
	r1 := mustRel(t, a30, b40)
	r2 := mustRel(t, a40, b30)
	assert.False(t, r1.Equal(r2))
	assert.NotEqual(t, r1.Hash(), r2.Hash(),
		"finished row wraps must distinguish swapped pairings")
}
