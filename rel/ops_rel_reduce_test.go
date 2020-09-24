package rel

import (
	"testing"
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
