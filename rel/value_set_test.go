package rel

import (
	"context"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsString(t *testing.T) {
	t.Parallel()

	generic, err := NewString([]rune("this is a test")).Map(func(v Value) (Value, error) { return v, nil })
	require.NoError(t, err)

	stringified, isString := generic.(String)
	require.True(t, isString, "%v", generic)
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

func TestUnionSetCallAll(t *testing.T) {
	t.Parallel()

	// unionset of relational sets
	set := MustNewSet(
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
		),
		NewTuple(
			NewAttr("@", NewNumber(4)),
			NewAttr("random", NewNumber(5)),
		),
	)
	ctx := arraictx.InitRunCtx(context.Background())
	AssertEqualValues(t, MustNewSet(NewNumber(42), NewNumber(24)), mustCallAll(ctx, set, NewNumber(1)))
	assert.Panics(t, func() { mustCallAll(ctx, set.With(NewTuple(NewAttr("@", NewNumber(1)))), NewNumber(2)) })
	assert.Panics(t, func() { mustCallAll(ctx, set.With(NewTuple(NewAttr("hi", NewNumber(1)))), NewNumber(3)) })
	AssertEqualValues(t, MustNewSet(NewNumber(5)), mustCallAll(ctx, set, NewNumber(4)))
	AssertEqualValues(t, None, mustCallAll(ctx, set, NewNumber(5)))
}
