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

	set := NewSet(
		foo(1, NewNumber(42)),
		foo(1, NewNumber(24)),
	)
	ctx := arraictx.InitRunCtx(context.Background())
	result, err := SetCall(ctx, set, NewNumber(1))
	assert.Error(t, err, "%v", result)
	result, err = SetCall(ctx, set, NewNumber(0))
	assert.Error(t, err, "%v", result)

	set = NewSet(
		foo(1, NewNumber(42)),
		foo(2, NewNumber(24)),
	)

	result, err = SetCall(ctx, set, NewNumber(1))
	require.NoError(t, err)
	assert.True(t, result.Equal(NewNumber(42)))
	result, err = SetCall(ctx, set, NewNumber(2))
	require.NoError(t, err)
	assert.True(t, result.Equal(NewNumber(24)))
}
