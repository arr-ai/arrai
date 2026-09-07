package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Every tuple representation hashes with one formula, so equal tuples hash
// equal regardless of which concrete type holds them.
func TestHash128TupleRepresentationsAgree(t *testing.T) {
	t.Parallel()

	// Build GenericTuples directly: NewTuple would canonicalise (@, @item)
	// and friends into the specialised kinds, which is what this test checks
	// against.
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
		assert.Equal(t, c.a.Hash128(), c.b.Hash128(), c.name)
		assert.Equal(t, c.a.Hash(0), c.b.Hash(0), c.name)
		assert.Equal(t, c.a.Hash(7), c.b.Hash(7), c.name)
	}
	assert.NotEqual(t,
		generic(NewAttr("a", NewNumber(1))).Hash128(),
		generic(NewAttr("b", NewNumber(1))).Hash128(), "different attr names")
	assert.NotEqual(t,
		generic(NewAttr("a", NewNumber(1))).Hash128(),
		generic(NewAttr("a", NewNumber(2))).Hash128(), "different attr values")
}

func TestHash128SeededDerivation(t *testing.T) {
	t.Parallel()
	v := NewNumber(42)
	assert.NotEqual(t, v.Hash(0), v.Hash(1))
	assert.Equal(t, v.Hash(3), v.Hash(3))
	assert.Equal(t, v.Hash128().Seeded(5), v.Hash(5))
}

func TestHash128CollectionsAreOrderIndependent(t *testing.T) {
	t.Parallel()
	mustSet := func(vs ...Value) Set {
		s, err := NewSet(vs...)
		require.NoError(t, err)
		return s
	}
	s1 := mustSet(NewNumber(1), NewNumber(2), NewNumber(3))
	s2 := mustSet(NewNumber(3), NewNumber(2), NewNumber(1))
	assert.Equal(t, s1.Hash128(), s2.Hash128())

	d1, err := NewDict(false, NewDictEntryTuple(NewNumber(1), NewNumber(2)), NewDictEntryTuple(NewNumber(3), NewNumber(4)))
	require.NoError(t, err)
	d2, err := NewDict(false, NewDictEntryTuple(NewNumber(3), NewNumber(4)), NewDictEntryTuple(NewNumber(1), NewNumber(2)))
	require.NoError(t, err)
	assert.Equal(t, d1.Hash128(), d2.Hash128())
	d3, err := NewDict(false, NewDictEntryTuple(NewNumber(1), NewNumber(2)), NewDictEntryTuple(NewNumber(3), NewNumber(5)))
	require.NoError(t, err)
	assert.NotEqual(t, d1.Hash128(), d3.Hash128(), "values participate")

	// A nested dict hashes the same however it is reached, and relations
	// hash the same as each other when equal.
	outer1 := NewTuple(NewAttr("d", d1))
	outer2 := NewTuple(NewAttr("d", d2))
	assert.Equal(t, outer1.Hash128(), outer2.Hash128())

	r1 := mustSet(NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewNumber(2))))
	r2 := mustSet(NewTuple(NewAttr("b", NewNumber(2)), NewAttr("a", NewNumber(1))))
	assert.Equal(t, r1.Hash128(), r2.Hash128())
}

func TestHash128StringsAndBytes(t *testing.T) {
	t.Parallel()
	assert.Equal(t, NewString([]rune("héllo")).Hash128(), NewString([]rune("héllo")).Hash128())
	assert.NotEqual(t, NewString([]rune("héllo")).Hash128(), NewString([]rune("hello")).Hash128())
	assert.NotEqual(t, NewString([]rune("1")).Hash128(), NewNumber(1).Hash128())
	assert.NotEqual(t, NewBytes([]byte("a")).Hash128(), NewString([]rune("a")).Hash128())
}
