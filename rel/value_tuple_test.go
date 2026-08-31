package rel

import (
	"testing"

	"github.com/arr-ai/hash/hash128"
	"github.com/stretchr/testify/assert"
)

// TestSpecializedTupleHash128AgreesWithGenericTuple guards against a
// regression like the one fixed in Relation/Dict.Hash128: GenericTuple.Equal
// is permissive across kinds (it compares against any Tuple, including
// these specialized ones), even though each specialized kind's own Equal is
// narrower (it only matches its own kind). So a specialized tuple and its
// generic equivalent must still hash the same in the direction that IS
// claimed equal — generic.Equal(specialized) — or a hash-based Set/Map can
// fail to recognize them as duplicates depending on which representation
// it happens to hold.
func TestSpecializedTupleHash128AgreesWithGenericTuple(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		t    Tuple
	}{
		{"StringCharTuple", NewStringCharTuple(2, 'x')},
		{"BytesByteTuple", NewBytesByteTuple(2, 7)},
		{"ArrayItemTuple", NewArrayItemTuple(2, NewNumber(7))},
		{"DictEntryTuple", NewDictEntryTuple(NewString([]rune("k")), NewNumber(7))},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			generic := c.t.(interface{ asGenericTuple() Tuple }).asGenericTuple()

			assert.False(t, c.t.Equal(generic),
				"specialized tuple's own Equal is narrower and doesn't match a generic tuple (by design)")
			assert.True(t, generic.Equal(c.t), "generic tuple must equal its specialized equivalent")
			assert.Equal(t, generic.(*GenericTuple).Hash128(), c.t.(interface{ Hash128() hash128.H128 }).Hash128(),
				"specialized tuple and its generic equivalent must hash the same, since generic.Equal considers them equal")
		})
	}
}

func TestTupleBuilder(t *testing.T) {
	t.Parallel()

	tb := &TupleBuilder{}
	tb.Put("@", NewNumber(0))
	tb.Put(StringCharAttr, NewNumber(1))
	tp := tb.Finish()
	assert.IsType(t, StringCharTuple{}, tp)
	assert.Equal(t, NewStringCharTuple(0, 1), tp)

	tb = &TupleBuilder{}
	tb.Put("@", NewNumber(0))
	tb.Put(BytesByteAttr, NewNumber(1))
	tp = tb.Finish()
	assert.IsType(t, BytesByteTuple{}, tp)
	assert.Equal(t, NewBytesByteTuple(0, 1), tp)

	tb = &TupleBuilder{}
	tb.Put("@", NewNumber(0))
	tb.Put(ArrayItemAttr, NewNumber(1))
	tp = tb.Finish()
	assert.IsType(t, ArrayItemTuple{}, tp)
	assert.Equal(t, NewArrayItemTuple(0, NewNumber(1)), tp)

	tb = &TupleBuilder{}
	tb.Put("@", NewNumber(0))
	tb.Put(DictValueAttr, NewNumber(1))
	tp = tb.Finish()
	assert.IsType(t, DictEntryTuple{}, tp)
	assert.Equal(t, NewDictEntryTuple(NewNumber(0), NewNumber(1)), tp)

	tb = &TupleBuilder{}
	tb.Put("@", NewNumber(0))
	tb.Put(DictValueAttr, NewNumber(1))
	tb.Put("@random", NewNumber(0))
	tp = tb.Finish()
	assert.IsType(t, &GenericTuple{}, tp)
	assert.True(
		t,
		newGenericTuple(
			NewAttr("@", NewNumber(0)),
			NewAttr(DictValueAttr, NewNumber(1)),
			NewAttr("@random", NewNumber(0)),
		).Equal(tp),
	)

	tb = &TupleBuilder{}
	tb.Put("a", NewNumber(0))
	tp = tb.Finish()
	assert.IsType(t, &GenericTuple{}, tp)
	assert.True(t, newGenericTuple(NewAttr("a", NewNumber(0))).Equal(tp))

	tb = &TupleBuilder{}
	tb.Put("@", NewNumber(0))
	tp = tb.Finish()
	assert.IsType(t, &GenericTuple{}, tp)
	assert.True(t, newGenericTuple(NewAttr("@", NewNumber(0))).Equal(tp))

	tp = (&TupleBuilder{}).Finish()
	assert.IsType(t, &GenericTuple{}, tp)
	assert.True(t, EmptyTuple.Equal(tp))
}
