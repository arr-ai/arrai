package syntax

import "testing"

func TestExprLet(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		`42`,
		`let x = 6; x * 7`,
	)
}

func TestExprLetPattern(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		`42`,
		`let 42 = 42; 42`,
	)
	AssertCodePanics(t,
		`let 42 = 1; 42`,
	)
}
