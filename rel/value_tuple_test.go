package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
