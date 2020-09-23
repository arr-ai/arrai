//nolint:dupl
package rel

import (
	"context"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/stretchr/testify/assert"
)

const hello = "hello"

func TestIsStringTuple(t *testing.T) {
	t.Parallel()
	for e := NewString([]rune(hello)).Enumerator(); e.MoveNext(); {
		tuple, is := e.Current().(StringCharTuple)
		if assert.True(t, is) {
			assert.Equal(t, rune(hello[tuple.at]), tuple.char)
		}
	}
}

func TestStringCallAll(t *testing.T) {
	t.Parallel()

	abc := NewString([]rune("abc"))
	ctx := arraictx.InitRunCtx(context.Background())
	AssertEqualValues(t, MustNewSet(NewNumber(float64('a'))), mustCallAll(ctx, abc, NewNumber(0)))
	AssertEqualValues(t, MustNewSet(NewNumber(float64('b'))), mustCallAll(ctx, abc, NewNumber(1)))
	AssertEqualValues(t, MustNewSet(NewNumber(float64('c'))), mustCallAll(ctx, abc, NewNumber(2)))
	AssertEqualValues(t, None, mustCallAll(ctx, abc, NewNumber(5)))
	AssertEqualValues(t, None, mustCallAll(ctx, abc, NewNumber(-1)))

	abc = NewOffsetString([]rune("abc"), -2)
	AssertEqualValues(t, MustNewSet(NewNumber(float64('a'))), mustCallAll(ctx, abc, NewNumber(-2)))
	AssertEqualValues(t, MustNewSet(NewNumber(float64('b'))), mustCallAll(ctx, abc, NewNumber(-1)))
	AssertEqualValues(t, MustNewSet(NewNumber(float64('c'))), mustCallAll(ctx, abc, NewNumber(0)))
	AssertEqualValues(t, None, mustCallAll(ctx, abc, NewNumber(1)))
	AssertEqualValues(t, None, mustCallAll(ctx, abc, NewNumber(-3)))

	abc = NewOffsetString([]rune("abc"), 2)
	AssertEqualValues(t, MustNewSet(NewNumber(float64('a'))), mustCallAll(ctx, abc, NewNumber(2)))
	AssertEqualValues(t, MustNewSet(NewNumber(float64('b'))), mustCallAll(ctx, abc, NewNumber(3)))
	AssertEqualValues(t, MustNewSet(NewNumber(float64('c'))), mustCallAll(ctx, abc, NewNumber(4)))
	AssertEqualValues(t, None, mustCallAll(ctx, abc, NewNumber(1)))
	AssertEqualValues(t, None, mustCallAll(ctx, abc, NewNumber(5)))
}
