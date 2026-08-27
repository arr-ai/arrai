package rel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testNestData = intPairs("a", "b", []intPair{
	{1, 1}, {1, 2}, {1, 3},
	{2, 1}, {2, 2},
}...)
var testNestNames = NewNames("a", "b")

func TestNestA(t *testing.T) {
	t.Parallel()
	AssertEqualValues(
		t,
		MustNewSet(
			NewTuple([]Attr{
				{"b", NewNumber(1)},
				{"g", intRel("a", 1, 2)},
			}...),
			NewTuple([]Attr{
				{"b", NewNumber(2)},
				{"g", intRel("a", 1, 2)},
			}...),
			NewTuple([]Attr{
				{"b", NewNumber(3)},
				{"g", intRel("a", 1)},
			}...),
		),
		Nest(testNestData, testNestNames, NewNames("a"), "g"),
	)
}

func TestNestB(t *testing.T) {
	t.Parallel()
	AssertEqualValues(
		t,
		MustNewSet(
			NewTuple([]Attr{
				{"a", NewNumber(1)},
				{"g", intRel("b", 1, 2, 3)},
			}...),
			NewTuple([]Attr{
				{"a", NewNumber(2)},
				{"g", intRel("b", 1, 2)},
			}...),
		),
		Nest(testNestData, testNestNames, NewNames("b"), "g"),
	)
}

func TestNestAThenB(t *testing.T) {
	t.Parallel()
	AssertEqualValues(
		t,
		MustNewSet(
			NewTuple([]Attr{
				{"g", intRel("a", 1)},
				{"h", intRel("b", 3)},
			}...),
			NewTuple([]Attr{
				{"g", intRel("a", 1, 2)},
				{"h", intRel("b", 1, 2)},
			}...),
		),
		Nest(
			Nest(
				testNestData,
				testNestNames,
				NewNames("a"),
				"g",
			),
			NewNames("b", "g"),
			NewNames("b"),
			"h",
		),
	)
}

func TestNestBThenA(t *testing.T) {
	t.Parallel()
	AssertEqualValues(
		t,
		MustNewSet(
			NewTuple([]Attr{
				{"g", intRel("a", 1)},
				{"h", intRel("b", 1, 2, 3)},
			}...),
			NewTuple([]Attr{
				{"g", intRel("a", 2)},
				{"h", intRel("b", 1, 2)},
			}...),
		),
		Nest(
			Nest(
				testNestData,
				testNestNames,
				NewNames("b"),
				"h",
			),
			NewNames("a", "h"),
			NewNames("a"),
			"g",
		),
	)
}

// TestNestManyRowsPerGroup grows each group well past the point where
// Reduce's internal bucket Set canonicalizes into a Relation (which embeds a
// slice and a map, so its zero value is not comparable with Go's ==).
// Reduce accumulates each group in a frozen.Map[Value, Value]; every row
// added to an existing key replaces that bucket's value, and frozen's
// Map.With compares old vs. new values to skip no-op writes. Using Value
// (rather than a bare `any`) as the bucket's value type lets frozen dispatch
// that comparison to Relation.Equal (Value satisfies frozen.Key[Value]),
// instead of falling back to `==`, which would panic once a group's
// accumulated Set became a Relation. See github.com/arr-ai/frozen#97.
func TestNestManyRowsPerGroup(t *testing.T) {
	t.Parallel()

	const groups, groupSize = 3, 50
	pairs := make([]intPair, 0, groups*groupSize)
	for a := 0; a < groups; a++ {
		for b := 0; b < groupSize; b++ {
			pairs = append(pairs, intPair{a, b})
		}
	}
	data := intPairs("a", "b", pairs...)

	nested := Nest(data, NewNames("a", "b"), NewNames("b"), "g")
	require.Equal(t, groups, nested.Count())
	for e := nested.Enumerator(); e.MoveNext(); {
		tuple := e.Current().(Tuple)
		g, has := tuple.Get("g")
		require.True(t, has)
		require.Equal(t, groupSize, g.(Set).Count())
	}
}
