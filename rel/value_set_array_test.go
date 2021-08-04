package rel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/arraictx"
)

func TestAsArray(t *testing.T) {
	t.Parallel()
	AssertEqualValues(t,
		NewArray(NewNumber(10), NewNumber(11)),
		MustNewSet(
			NewArrayItemTuple(0, NewNumber(10)),
			NewArrayItemTuple(1, NewNumber(11)),
		),
	)
	AssertEqualValues(t,
		NewOffsetArray(2, NewNumber(10), NewNumber(11)),
		MustNewSet(
			NewArrayItemTuple(2, NewNumber(10)),
			NewArrayItemTuple(3, NewNumber(11)),
		),
	)
}

func TestAsArrayHoles(t *testing.T) {
	t.Parallel()
	AssertEqualValues(t,
		NewArray(NewNumber(1), nil, nil, NewNumber(2)),
		MustNewSet(
			NewArrayItemTuple(0, NewNumber(1)),
			NewArrayItemTuple(3, NewNumber(2)),
		),
	)
	AssertEqualValues(t,
		NewOffsetArray(2, NewNumber(1), nil, nil, NewNumber(2)),
		MustNewSet(
			NewArrayItemTuple(2, NewNumber(1)),
			NewArrayItemTuple(5, NewNumber(2)),
		),
	)
}

func TestArrayWithout(t *testing.T) {
	t.Parallel()

	three := NewArray(NewNumber(10), NewNumber(11), NewNumber(12))

	// Without first item
	AssertEqualValues(t,
		NewOffsetArray(1, NewNumber(11), NewNumber(12)),
		three.Without(NewArrayItemTuple(0, NewNumber(10))),
	)
	AssertExprEvalsToType(t,
		Array{},
		three.Without(NewArrayItemTuple(0, NewNumber(10))),
	)

	// Without middle item
	AssertEqualValues(t,
		NewArray(NewNumber(10), nil, NewNumber(12)),
		three.Without(NewArrayItemTuple(1, NewNumber(11))),
	)
	AssertExprEvalsToType(t,
		Array{},
		three.Without(NewArrayItemTuple(1, NewNumber(11))),
	)

	// Without last item
	AssertEqualValues(t,
		NewArray(NewNumber(10), NewNumber(11)),
		three.Without(NewArrayItemTuple(2, NewNumber(12))),
	)
	AssertExprEvalsToType(t,
		Array{},
		three.Without(NewArrayItemTuple(2, NewNumber(12))),
	)

	four := NewArray(NewNumber(10), NewNumber(11), NewNumber(12), NewNumber(13))

	AssertEqualValues(t,
		NewOffsetArray(1, NewNumber(11), NewNumber(12)),
		four.Without(NewArrayItemTuple(3, NewNumber(13))).Without(NewArrayItemTuple(0, NewNumber(10))),
	)
	AssertEqualValues(t,
		NewOffsetArray(1, NewNumber(11), NewNumber(12)),
		four.Without(NewArrayItemTuple(0, NewNumber(10))).Without(NewArrayItemTuple(3, NewNumber(13))),
	)
}

func TestArrayCallAll(t *testing.T) {
	t.Parallel()

	three := NewArray(NewNumber(10), NewNumber(11), NewNumber(12))
	ctx := arraictx.InitRunCtx(context.Background())
	AssertEqualValues(t, MustNewSet(NewNumber(10)), mustCallAll(ctx, three, NewNumber(0)))
	AssertEqualValues(t, MustNewSet(NewNumber(11)), mustCallAll(ctx, three, NewNumber(1)))
	AssertEqualValues(t, MustNewSet(NewNumber(12)), mustCallAll(ctx, three, NewNumber(2)))
	AssertEqualValues(t, None, mustCallAll(ctx, three, NewNumber(5)))
	AssertEqualValues(t, None, mustCallAll(ctx, three, NewNumber(-1)))

	three = NewOffsetArray(-2, NewNumber(10), NewNumber(11), NewNumber(12))
	AssertEqualValues(t, MustNewSet(NewNumber(10)), mustCallAll(ctx, three, NewNumber(-2)))
	AssertEqualValues(t, MustNewSet(NewNumber(11)), mustCallAll(ctx, three, NewNumber(-1)))
	AssertEqualValues(t, MustNewSet(NewNumber(12)), mustCallAll(ctx, three, NewNumber(0)))
	AssertEqualValues(t, None, mustCallAll(ctx, three, NewNumber(1)))
	AssertEqualValues(t, None, mustCallAll(ctx, three, NewNumber(-3)))

	three = NewOffsetArray(2, NewNumber(10), NewNumber(11), NewNumber(12))
	AssertEqualValues(t, MustNewSet(NewNumber(10)), mustCallAll(ctx, three, NewNumber(2)))
	AssertEqualValues(t, MustNewSet(NewNumber(11)), mustCallAll(ctx, three, NewNumber(3)))
	AssertEqualValues(t, MustNewSet(NewNumber(12)), mustCallAll(ctx, three, NewNumber(4)))
	AssertEqualValues(t, None, mustCallAll(ctx, three, NewNumber(1)))
	AssertEqualValues(t, None, mustCallAll(ctx, three, NewNumber(5)))

	b := NewSetBuilder()
	err := three.CallAll(ctx, NewString([]rune("0")), b)
	if assert.NoError(t, err) {
		set, err := b.Finish()
		require.NoError(t, err)
		assert.False(t, set.IsTrue())
	}
}

func TestArrayWhere(t *testing.T) {
	t.Parallel()

	three := NewArray(NewNumber(10), NewNumber(11), NewNumber(12))

	atBetween := func(a, b int) func(v Value) bool {
		return func(v Value) bool {
			i := int(v.(ArrayItemTuple).MustGet("@").(Number).Float64())
			return a <= i && i <= b
		}
	}

	where := func(s Set, p func(v Value) bool) Set {
		result, err := s.Where(func(v Value) (bool, error) { return p(v), nil })
		require.NoError(t, err)
		return result
	}

	AssertEqualValues(t, three, where(three, atBetween(0, 2)))
	AssertEqualValues(t, NewArray(NewNumber(10), NewNumber(11)), where(three, atBetween(0, 1)))
	AssertEqualValues(t, NewArray(NewNumber(10)), where(three, atBetween(0, 0)))
	AssertEqualValues(t, None, where(three, atBetween(-1, -1)))

	AssertEqualValues(t, None, where(three, atBetween(3, 3)))
	AssertEqualValues(t, NewOffsetArray(2, NewNumber(12)), where(three, atBetween(2, 3)))
	AssertEqualValues(t, NewOffsetArray(1, NewNumber(11), NewNumber(12)), where(three, atBetween(1, 3)))
	AssertEqualValues(t, three, where(three, atBetween(0, 3)))

	offsetThree := NewOffsetArray(-2, NewNumber(10), NewNumber(11), NewNumber(12))

	AssertEqualValues(t, offsetThree, where(offsetThree, atBetween(-2, 0)))
	AssertEqualValues(t, NewOffsetArray(-2, NewNumber(10), NewNumber(11)), where(offsetThree, atBetween(-2, -1)))
	AssertEqualValues(t, NewOffsetArray(-2, NewNumber(10)), where(offsetThree, atBetween(-2, -2)))
	AssertEqualValues(t, None, where(offsetThree, atBetween(-3, -3)))

	AssertEqualValues(t, None, where(offsetThree, atBetween(1, 1)))
	AssertEqualValues(t, NewArray(NewNumber(12)), where(offsetThree, atBetween(0, 1)))
	AssertEqualValues(t, NewOffsetArray(-1, NewNumber(11), NewNumber(12)), where(offsetThree, atBetween(-1, 1)))
	AssertEqualValues(t, offsetThree, where(offsetThree, atBetween(-2, 1)))
}
