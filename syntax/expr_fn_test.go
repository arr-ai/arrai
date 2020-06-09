package syntax

import "testing"

func TestExprLet(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `3`, `(\[x, y] x + y)([1, 2])`)
	AssertCodesEvalToSameValue(t, `3`, `(\z \[x, y] z/(x + y))(9, [1, 2])`)
}
