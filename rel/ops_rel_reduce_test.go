package rel

import (
	"testing"
)

var nestData = intPairs("a", "b", []intPair{
	{1, 1}, {1, 2}, {1, 3},
	{2, 1}, {2, 2},
}...)

// TestNestA tests nesting attr a of the test data set.
func TestNestA(t *testing.T) {
	AssertEqualValues(
		t,
		NewSet(
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
		Nest(nestData, NewNames("a"), "g"),
	)
}

// TestNestB tests nesting attr b of the test data set.
func TestNestB(t *testing.T) {
	AssertEqualValues(
		t,
		NewSet(
			NewTuple([]Attr{
				{"a", NewNumber(1)},
				{"g", intRel("b", 1, 2, 3)},
			}...),
			NewTuple([]Attr{
				{"a", NewNumber(2)},
				{"g", intRel("b", 1, 2)},
			}...),
		),
		Nest(nestData, NewNames("b"), "g"),
	)
}

// TestNestAThenB tests nesting attr b then a of the test data set.
func TestNestAThenB(t *testing.T) {
	AssertEqualValues(
		t,
		NewSet(
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
				nestData,
				NewNames("a"),
				"g"),
			NewNames("b"),
			"h"),
	)
}

// TestNestBThenA tests nesting attr b then a of the test data set.
func TestNestBThenA(t *testing.T) {
	AssertEqualValues(
		t,
		NewSet(
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
				nestData,
				NewNames("b"),
				"h"),
			NewNames("a"),
			"g"),
	)
}
