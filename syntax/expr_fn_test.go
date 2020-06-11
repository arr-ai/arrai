package syntax

import "testing"

func TestExprFn(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `3`, `(\[x, y] x + y)([1, 2])`)
	AssertCodesEvalToSameValue(t, `3`, `(\z \[x, y] z/(x + y))(9, [1, 2])`)
	AssertCodeErrors(t, `(\[x, y] 42)([1, 2, 3])`, "")
	AssertCodeErrors(t, `(\[x, y] x + y)([1, 2, 3])`, "")
}
