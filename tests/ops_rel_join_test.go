package tests

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
)

var (
	odds   = intRel("a", 1, 3, 5, 7, 9, 11)
	threes = intRel("a", 3, 6, 9, 12)
)

func TestJoinNone(t *testing.T) {
	assertJoin(t, rel.None, rel.None, rel.None)
}

func TestJoinTrue(t *testing.T) {
	assertJoin(t, rel.None, rel.None, rel.True)
	assertJoin(t, rel.None, rel.True, rel.None)
	assertJoin(t, rel.True, rel.True, rel.True)
}

func TestJoinTrue1(t *testing.T) {
	assertJoin(t, odds, odds, rel.True)
	assertJoin(t, threes, threes, rel.True)
}

func TestJoinSelf1(t *testing.T) {
	assertJoin(t, odds, odds, odds)
	assertJoin(t, threes, threes, threes)
}

func TestJoinIntersect1(t *testing.T) {
	assertJoin(t, intRel("a", 3, 9), odds, threes)
}

func TestJoinProduct2(t *testing.T) {
	assertJoin(t,
		intPairs("a", "b", []intPair{
			{1, 1}, {1, 2}, {1, 3},
			{2, 1}, {2, 2}, {2, 3},
		}...),
		intRel("a", 1, 2),
		intRel("b", 1, 2, 3))
}

func TestJoinIntersect2(t *testing.T) {
	assertJoin(t,
		intPairs("a", "b", []intPair{{3, 4}}...),
		intPairs("a", "b", []intPair{{1, 2}, {3, 4}, {5, 6}}...),
		intPairs("a", "b", []intPair{{1, 4}, {3, 4}, {5, 4}}...))
}

// Helpers

func intRel(name string, values ...int) rel.Set {
	result := rel.None
	for _, value := range values {
		result = result.With(
			rel.NewTuple(rel.Attr{name, rel.NewNumber(float64(value))}))
	}
	return result
}

type intPair [2]int

func intPairs(a, b string, pairs ...intPair) rel.Set {
	result := rel.None
	for _, pair := range pairs {
		result = result.With(
			rel.NewTuple(
				rel.Attr{a, rel.NewNumber(float64(pair[0]))},
				rel.Attr{b, rel.NewNumber(float64(pair[1]))},
			),
		)
	}
	return result
}

func intFunc(a, b string, xes []int, f func(int) int) rel.Set {
	pairs := make([]intPair, len(xes))
	for i, x := range xes {
		pairs[i] = intPair{x, f(x)}
	}
	return intPairs(a, b, pairs...)
}

func assertJoin(t *testing.T, expected, a, b rel.Set) bool {
	return assertEqualValues(t, expected, rel.Join(a, b)) &&
		assertEqualValues(t, expected, rel.Join(b, a))
}
