package syntax

import "testing"

func TestApplyExpr(t *testing.T) {
	AssertCodesEvalToSameValue(t, `42`, `6 -> 7 * .`)
	AssertCodesEvalToSameValue(t, `42`, `6 -> \x 7 * x`)
}

func TestApplyExprInsideMapExpr(t *testing.T) {
	AssertCodesEvalToSameValue(t, `{2, 4, 8}`, `{1, 2, 3} => (2 -> \y y ^ .)`)
	AssertCodesEvalToSameValue(t, `{2, 4, 8}`, `(\z {1, 2, 3} => (z -> \y y ^ .))(2)`)
}
