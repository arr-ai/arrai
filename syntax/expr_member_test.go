package syntax

import "testing"

func TestExpr(t *testing.T) {
	AssertCodesEvalToSameValue(t, `true`, `2 <: {1, 2, 3}`)
	AssertCodesEvalToSameValue(t, `false`, `42 <: {1, 2, 3}`)
}
