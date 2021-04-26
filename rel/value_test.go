package rel

import (
	"context"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetCall(t *testing.T) {
	t.Parallel()

	foo := func(at int, v Value) Tuple {
		return NewTuple(NewAttr("@", NewNumber(float64(at))), NewAttr("@foo", v))
	}

	set := MustNewSet(
		foo(1, NewNumber(42)),
		foo(1, NewNumber(24)),
	)
	ctx := arraictx.InitRunCtx(context.Background())
	result, err := SetCall(ctx, set, NewNumber(1))
	assert.Error(t, err, "%v", result)
	result, err = SetCall(ctx, set, NewNumber(0))
	assert.Error(t, err, "%v", result)

	set = MustNewSet(
		foo(1, NewNumber(42)),
		foo(2, NewNumber(24)),
	)

	result, err = SetCall(ctx, set, NewNumber(1))
	require.NoError(t, err)
	AssertEqualValues(t, result, NewNumber(42))

	result, err = SetCall(ctx, set, NewNumber(2))
	require.NoError(t, err)
	AssertEqualValues(t, result, NewNumber(24))
}

type Foo struct {
	a int
	b int
}

func TestNewValue(t *testing.T) {
	x := []interface{}{map[string]interface{}{"a": 1, "b": 2}}

	actual, err := NewValue(x)
	require.NoError(t, err)

	expected, err := NewSet(NewTuple(NewIntAttr("a", 1), NewIntAttr("b", 2)))
	require.NoError(t, err)

	y := []*Foo{{1, 2}}
	actual, err = NewValue(y)
	require.NoError(t, err)
	AssertEqualValues(t, expected, actual)

	//AssertEqualValues(t, expected, actual)
}
