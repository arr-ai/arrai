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

func TestExprLetNumPattern(t *testing.T) {
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
