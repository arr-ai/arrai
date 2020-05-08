package syntax

import "testing"

func TestExprLetIdentPattern(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		`7`,
		`let x = 6; 7`,
	)
	AssertCodesEvalToSameValue(t,
		`42`,
		`let x = 6; x * 7`,
	)
	AssertCodesEvalToSameValue(t,
		`[1, 2]`,
		`let x = 1; [x, 2]`,
	)
}

func TestExprLetValuePattern(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		`42`,
		`let 42 = 42; 42`,
	)
	AssertCodesEvalToSameValue(t,
		`1`,
		`let 42 = 42; 1`,
	)
	AssertCodePanics(t,
		`let 42 = 1; 42`,
	)
	AssertCodePanics(t,
		`let 42 = 1; 1`,
	)
}

func TestExprLetArrayPattern(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		`9`,
		`let [a, b, c] = [1, 2, 3]; 9`,
	)
	AssertCodesEvalToSameValue(t,
		`[1, 2, 3]`,
		`let [a, b, c] = [1, 2, 3]; [a, b, c]`,
	)
	AssertCodesEvalToSameValue(t,
		`2`,
		`let [a, b, c] = [1, 2, 3]; [a, b, c](1)`,
	)
	AssertCodesEvalToSameValue(t,
		`2`,
		`let [a, b, c] = [1, 2, 3]; b`,
	)
	AssertCodesEvalToSameValue(t,
		`2`,
		`let [a, b, c] = [1, 2, 3]; [c, b](1)`,
	)
	AssertCodesEvalToSameValue(t,
		`2`,
		`let arr = [1, 2]; let [a, b] = arr; b`,
	)
}
