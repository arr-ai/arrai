package tests

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
)

var nestData = intPairs("a", "b", []intPair{
	{1, 1}, {1, 2}, {1, 3},
	{2, 1}, {2, 2},
}...)

// TestNestA tests nesting attr a of the test data set.
func TestNestA(t *testing.T) {
	assertEqualValues(
		t,
		rel.NewSet(
			rel.NewTuple([]rel.Attr{
				{"b", rel.NewNumber(1)},
				{"g", intRel("a", 1, 2)},
			}...),
			rel.NewTuple([]rel.Attr{
				{"b", rel.NewNumber(2)},
				{"g", intRel("a", 1, 2)},
			}...),
			rel.NewTuple([]rel.Attr{
				{"b", rel.NewNumber(3)},
				{"g", intRel("a", 1)},
			}...),
		),
		rel.Nest(nestData, rel.NewNames("a"), "g"),
	)
}

// TestNestB tests nesting attr b of the test data set.
func TestNestB(t *testing.T) {
	assertEqualValues(
		t,
		rel.NewSet(
			rel.NewTuple([]rel.Attr{
				{"a", rel.NewNumber(1)},
				{"g", intRel("b", 1, 2, 3)},
			}...),
			rel.NewTuple([]rel.Attr{
				{"a", rel.NewNumber(2)},
				{"g", intRel("b", 1, 2)},
			}...),
		),
		rel.Nest(nestData, rel.NewNames("b"), "g"),
	)
}

// TestNestAThenB tests nesting attr b then a of the test data set.
func TestNestAThenB(t *testing.T) {
	assertEqualValues(
		t,
		rel.NewSet(
			rel.NewTuple([]rel.Attr{
				{"g", intRel("a", 1)},
				{"h", intRel("b", 3)},
			}...),
			rel.NewTuple([]rel.Attr{
				{"g", intRel("a", 1, 2)},
				{"h", intRel("b", 1, 2)},
			}...),
		),
		rel.Nest(
			rel.Nest(
				nestData,
				rel.NewNames("a"),
				"g"),
			rel.NewNames("b"),
			"h"),
	)
}

// TestNestBThenA tests nesting attr b then a of the test data set.
func TestNestBThenA(t *testing.T) {
	assertEqualValues(
		t,
		rel.NewSet(
			rel.NewTuple([]rel.Attr{
				{"g", intRel("a", 1)},
				{"h", intRel("b", 1, 2, 3)},
			}...),
			rel.NewTuple([]rel.Attr{
				{"g", intRel("a", 2)},
				{"h", intRel("b", 1, 2)},
			}...),
		),
		rel.Nest(
			rel.Nest(
				nestData,
				rel.NewNames("b"),
				"h"),
			rel.NewNames("a"),
			"g"),
	)
}
