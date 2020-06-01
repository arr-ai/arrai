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

func TestApplyExprWithPattern(t *testing.T) {
	AssertCodesEvalToSameValue(t, `2`, `1 -> \x let [(x), y] = [1, 2]; y`)
	AssertCodesEvalToSameValue(t, `3`, `[1, 2] -> \[x, y] x + y`)
	AssertCodesEvalToSameValue(t, `3`, `(m: 1, n: 2) -> \(m: x, n: y) x + y`)
	AssertCodesEvalToSameValue(t, `3`, `{"m": 1, "n": 2} -> \{"m": x, "n": y} x + y`)
	AssertCodesEvalToSameValue(t, `6`, `[1, [2, 3]] -> \[x, [y, z]] x + y + z`)
	AssertCodesEvalToSameValue(t, `[3, [3, 4]]`, `[1, 2, 3, 4] -> \[x, y, ...t] [x + y, t]`)
	AssertCodeErrors(t, `1 -> \x let [(x), y] = [2, 2]; y`, "")
}
