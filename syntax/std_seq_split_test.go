package syntax

import "testing"

func TestStrSplit(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`["t", "h", "i", "s", " ", "i", "s", " ", "a", " ", "t", "e", "s", "t"]`,
		`//seq.split("","this is a test")`)
	AssertCodesEvalToSameValue(t, `["this", "is", "a", "test"]`, `//seq.split(" ","this is a test") `)
	AssertCodesEvalToSameValue(t, `["this is a test"]         `, `//seq.split(",","this is a test") `)
	AssertCodesEvalToSameValue(t, `["th", " ", " a test"]     `, `//seq.split("is","this is a test")`)
	assertExprPanics(t, `//seq.split(1, "this is a test")`)
}

func TestArraySplit(t *testing.T) {
	t.Parallel()
	// TODO
	// AssertCodesEvalToSameValue(t,
	// 	`[['B'],['C', 'D', 'E']]`,
	// 	`//seq.sub(['A', 'B', 'A', 'C', 'D', 'E'], 'A')`)
}

func TestBytesSplit(t *testing.T) {
	t.Parallel()
}
