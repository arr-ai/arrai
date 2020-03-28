package syntax

import "testing"

func TestApplyExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `6 -> 7 * .`)
	AssertCodesEvalToSameValue(t, `42`, `6 -> . * 7`)
	AssertCodesEvalToSameValue(t, `42`, `6 -> \x 7 * x`)
}

func TestApplyExprInsideMapExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{2, 4, 8}`, `{1, 2, 3} => (2 -> \y y ^ .)`)
	AssertCodesEvalToSameValue(t, `{2, 4, 8}`, `(\z {1, 2, 3} => (z -> \y y ^ .))(2)`)
}

func TestTransformExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `[2, 4, 6, 8] `, `[1, 2, 3, 4] >> . * 2`)
	AssertCodesEvalToSameValue(t, `[1, 4, 9, 16]`, `[1, 2, 3, 4] >> . * .`)
}
