package syntax

import "testing"

func TestExprLet(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		`42`,
		`let x = 6; x * 7`,
	)
}
