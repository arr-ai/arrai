package syntax

import "testing"

func TestExprMemberOf(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `true`, `2 <: {1, 2, 3}`)
	AssertCodesEvalToSameValue(t, `false`, `42 <: {1, 2, 3}`)
}

func TestExprNotMemberOf(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `false`, `2 !<: {1, 2, 3}`)
	AssertCodesEvalToSameValue(t, `true`, `42 !<: {1, 2, 3}`)
}
