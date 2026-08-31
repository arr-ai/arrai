package rel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/arraictx"
)

// TestArrayHash128DistinguishesRepeatedElementValue guards against a
// cancellation bug: Array.Hash128 combined per-index item hashes with xor,
// but each item hash already xors together an index-hash and a value-hash
// via hashTuple2. When the same value occurs at two different indices (e.g.
// ["x", "x"] vs ["y", "y"]), the value-dependent part of those two item
// hashes is identical and cancels out under xor, making the array's hash
// completely independent of the repeated value -- so, for example,
// ["AppA", "AppA"] and ["AppB", "AppB"] hashed identically despite being
// unequal, silently corrupting any hash-based Set/Map containing such
// arrays as elements or as part of a composite key (found via a real
// fully-qualified name convention that repeats a path segment).
//
// Covers repeats at adjacent (0,1) and non-adjacent (0,2) indices, an
// all-elements-equal array with an even element count (an odd count doesn't
// fully cancel — xor-ing the same value-dependent term in three times
// leaves one copy, so a 3-element all-equal array would pass even against
// the old, buggy formula), and a repeat one level down inside a nested
// array.
func TestArrayHash128DistinguishesRepeatedElementValue(t *testing.T) {
	t.Parallel()

	str := func(s string) Value { return NewString([]rune(s)) }

	cases := []struct {
		name string
		a, b Set
	}{
		{
			name: "adjacent repeat",
			a:    NewArray(str("x"), str("x")),
			b:    NewArray(str("y"), str("y")),
		},
		{
			name: "non-adjacent repeat",
			a:    NewArray(str("x"), str("mid"), str("x")),
			b:    NewArray(str("y"), str("mid"), str("y")),
		},
		{
			name: "all elements equal, even count",
			a:    NewArray(str("x"), str("x"), str("x"), str("x")),
			b:    NewArray(str("y"), str("y"), str("y"), str("y")),
		},
		{
			name: "repeat nested one level down",
			a:    NewArray(NewArray(str("x"), str("x")), str("tail")),
			b:    NewArray(NewArray(str("y"), str("y")), str("tail")),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			assert.False(t, c.a.Equal(c.b), "arrays with different repeated values must not be equal")
			assert.NotEqual(t, c.a.(Array).Hash128(), c.b.(Array).Hash128(),
				"arrays with different repeated values must not hash the same")
		})
	}
}

// FuzzArrayHash128RepeatedValue generalizes
// TestArrayHash128DistinguishesRepeatedElementValue's fixed cases across
// arbitrary element values and repeat counts, so it can catch a
// reintroduction of the xor-cancellation bug (or a similar one) for
// combinations the hand-written cases don't happen to cover. Run with
// `go test -fuzz=FuzzArrayHash128RepeatedValue ./rel/`.
func FuzzArrayHash128RepeatedValue(f *testing.F) {
	f.Add("x", "y", uint8(2))
	f.Add("AppA", "AppB", uint8(2))
	f.Add("x", "y", uint8(3))
	f.Add("x", "y", uint8(4))
	f.Fuzz(func(t *testing.T, x, y string, rawCount uint8) {
		// Invalid UTF-8 bytes collapse to the same replacement rune, so
		// compare post-[]rune-round-trip -- that's what NewString actually
		// sees -- rather than the raw fuzzer-provided strings.
		xr, yr := string([]rune(x)), string([]rune(y))
		if xr == yr {
			return
		}
		count := int(rawCount%6) + 2 // 2..7 repeats of the same value.
		xs := make([]Value, count)
		ys := make([]Value, count)
		for i := range xs {
			xs[i] = NewString([]rune(xr))
			ys[i] = NewString([]rune(yr))
		}
		a, b := NewArray(xs...).(Array), NewArray(ys...).(Array)

		if a.Equal(b) {
			t.Fatalf("arrays of %d copies of %q and %q must not be equal", count, xr, yr)
		}
		if a.Hash128() == b.Hash128() {
			t.Fatalf("arrays of %d copies of %q and %q hashed the same", count, xr, yr)
		}
	})
}

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
