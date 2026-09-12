//nolint:dupl
package rel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/arraictx"
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

	b := NewSetBuilder()
	err := abc.CallAll(ctx, NewString([]rune("0")), b)
	if assert.NoError(t, err) {
		set, err := b.Finish()
		require.NoError(t, err)
		assert.False(t, set.IsTrue())
	}
}

func TestStringUTF8Backing(t *testing.T) {
	t.Parallel()

	const cafe = "café"
	fromRunes := NewString([]rune(cafe)).(String)
	fromGo := NewGoString(cafe).(String)
	require.NotNil(t, fromRunes.utf8)
	require.NotNil(t, fromGo.utf8)
	assert.Nil(t, fromRunes.ascii)
	assert.Nil(t, fromRunes.s)
	assert.Equal(t, 4, fromRunes.Count())
	assert.True(t, fromRunes.Equal(fromGo))
	assert.Equal(t, fromRunes.Hash128(), fromGo.Hash128())
	assert.Equal(t, cafe, fromRunes.goString())

	ascii := NewGoString("cafe").(String)
	require.NotNil(t, ascii.ascii)
	assert.False(t, ascii.Equal(fromGo))

	ctx := arraictx.InitRunCtx(context.Background())
	AssertEqualValues(t, MustNewSet(NewNumber(float64('é'))), mustCallAll(ctx, fromGo, NewNumber(3)))
	AssertEqualValues(t, MustNewSet(NewNumber(float64('c'))), mustCallAll(ctx, fromGo, NewNumber(0)))

	joined := concatStrings(fromGo, NewGoString("!").(String))
	assert.Equal(t, "café!", joined.goString())
	assert.Equal(t, 5, joined.Count())
	assert.NotNil(t, joined.utf8)
}
