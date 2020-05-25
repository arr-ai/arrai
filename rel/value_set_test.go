package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsString(t *testing.T) {
	t.Parallel()

	generic := NewString([]rune("this is a test")).Map(func(v Value) Value { return v })
	stringified, isString := AsString(generic)
	require.True(t, isString)
	assert.True(t, stringified.Equal(NewString([]rune("this is a test"))), stringified.String())

	// generic = NewOffsetString([]rune("this is a test"), 100).Map(func(v Value) Value { return v })
	// stringified, isString = AsString(generic)
	// require.True(t, isString)
	// assert.True(t, stringified.Equal(NewOffsetString([]rune("this is a test"), 100)))

	// generic = NewString([]rune("")).Map(func(v Value) Value { return v })
	// stringified, isString = AsString(generic)
	// require.True(t, isString)
	// assert.True(t, stringified.Equal(NewString([]rune(""))))
}

func TestGenericSetCallAll(t *testing.T) {
	t.Parallel()

	set := NewSet(
		NewTuple(
			NewAttr("@", NewNumber(1)),
			NewAttr("@fooo", NewNumber(42)),
		),
		NewTuple(
			NewAttr("@", NewNumber(1)),
			NewAttr("@baar", NewNumber(24)),
		),
		NewTuple(
			NewAttr("@", NewNumber(2)),
			NewAttr("@foo", NewNumber(3)),
			NewAttr("@bar", NewNumber(4)),
		),
		NewTuple(
			NewAttr("@", NewNumber(3)),
		),
		NewTuple(
			NewAttr("@", NewNumber(4)),
			NewAttr("random", NewNumber(5)),
		),
	)

	AssertEqualValues(t, NewSet(NewNumber(42), NewNumber(24)), set.CallAll(NewNumber(1)))
	assert.Panics(t, func() { set.CallAll(NewNumber(2)) })
	assert.Panics(t, func() { set.CallAll(NewNumber(3)) })
	AssertEqualValues(t, NewSet(NewNumber(5)), set.CallAll(NewNumber(4)))
	AssertEqualValues(t, None, set.CallAll(NewNumber(5)))
}
