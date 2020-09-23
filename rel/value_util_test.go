package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueTypeAsString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `number`, ValueTypeAsString(NewNumber(0)))
	assert.Equal(t, `array`, ValueTypeAsString(NewArray(NewNumber(0))))
	assert.Equal(t, `bytes`, ValueTypeAsString(NewBytes([]byte{0})))
	assert.Equal(t, `closure`, ValueTypeAsString(NewClosure(Scope{}, nil)))
	assert.Equal(t, `dict`, ValueTypeAsString(
		NewDict(false, NewDictEntryTuple(NewNumber(0), NewNumber(0)))))
	assert.Equal(t, `expr closure`, ValueTypeAsString(NewExprClosure(Scope{}, nil)))
	assert.Equal(t, `string`, ValueTypeAsString(NewString([]rune(" "))))
	assert.Equal(t, `set`, ValueTypeAsString(None))
	assert.Equal(t, `array item tuple`, ValueTypeAsString(NewArrayItemTuple(0, nil)))
	assert.Equal(t, `bytes byte tuple`, ValueTypeAsString(NewBytesByteTuple(0, 0)))
	assert.Equal(t, `dict entry tuple`, ValueTypeAsString(NewDictEntryTuple(nil, nil)))
	assert.Equal(t, `string char tuple`, ValueTypeAsString(NewStringCharTuple(0, ' ')))
	assert.Equal(t, `tuple`, ValueTypeAsString(NewTuple()))

	// Less obvious cases resulting from implementation details.
	assert.Equal(t, `set`, ValueTypeAsString(NewString([]rune{})))
	assert.Equal(t, `set`, ValueTypeAsString(NewDict(false)))
	assert.Equal(t, `set`, ValueTypeAsString(NewArray()))
	assert.Equal(t, `set`, ValueTypeAsString(NewBytes([]byte{})))
}
