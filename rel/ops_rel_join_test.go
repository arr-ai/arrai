package rel

import (
	"testing"
)

var (
	odds   = intRel("a", 1, 3, 5, 7, 9, 11)
	threes = intRel("a", 3, 6, 9, 12)
)

func TestJoinNone(t *testing.T) {
	t.Parallel()
	assertJoin(t, None, None, None)
}

func TestJoinTrue(t *testing.T) {
	t.Parallel()
	assertJoin(t, None, None, True)
	assertJoin(t, None, True, None)
	assertJoin(t, True, True, True)
}

func TestJoinTrue1(t *testing.T) {
	t.Parallel()
	assertJoin(t, odds, odds, True)
	assertJoin(t, threes, threes, True)
}

func TestJoinSelf1(t *testing.T) {
	t.Parallel()
	assertJoin(t, odds, odds, odds)
	assertJoin(t, threes, threes, threes)
}

func TestJoinIntersect1(t *testing.T) {
	t.Parallel()
	assertJoin(t, intRel("a", 3, 9), odds, threes)
}

func TestJoinProduct2(t *testing.T) {
	t.Parallel()
	assertJoin(t,
		intPairs("a", "b", []intPair{
			{1, 1}, {1, 2}, {1, 3},
			{2, 1}, {2, 2}, {2, 3},
		}...),
		intRel("a", 1, 2),
		intRel("b", 1, 2, 3))
}

func TestJoinIntersect2(t *testing.T) {
	t.Parallel()
	assertJoin(t,
		intPairs("a", "b", []intPair{{3, 4}}...),
		intPairs("a", "b", []intPair{{1, 2}, {3, 4}, {5, 6}}...),
		intPairs("a", "b", []intPair{{1, 4}, {3, 4}, {5, 4}}...))
}

func TestJoinSpecialSet(t *testing.T) {
	t.Parallel()
	assertJoin(t,
		NewString([]rune("do")),
		NewString([]rune("dots")),
		NewSet(NewTuple(NewAttr("@", NewNumber(0))), NewTuple(NewAttr("@", NewNumber(1)))),
	)
}

// Helpers

func intRel(name string, values ...int) Set {
	result := None
	for _, value := range values {
		result = result.With(
			NewTuple(Attr{name, NewNumber(float64(value))}))
	}
	return result
}

type intPair [2]int

func intPairs(a, b string, pairs ...intPair) Set { //nolint:unparam
	result := None
	for _, pair := range pairs {
		result = result.With(
			NewTuple(
				Attr{a, NewNumber(float64(pair[0]))},
				Attr{b, NewNumber(float64(pair[1]))},
			),
		)
	}
	return result
}

func assertJoin(t *testing.T, expected, a, b Set) bool { //nolint:unparam
	return AssertEqualValues(t, expected, Join(a, b)) &&
		AssertEqualValues(t, expected, Join(b, a))
}
