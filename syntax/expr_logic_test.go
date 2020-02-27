package syntax

import "testing"

func TestExprIfElse(t *testing.T) {
	AssertCodesEvalToSameValue(t, "1", "1 if true else 2")
	AssertCodesEvalToSameValue(t, "2", "1 if false else 2")
}

func TestExprIfNoElse(t *testing.T) {
	AssertCodesEvalToSameValue(t, "1", "1 if true")
	AssertCodesEvalToSameValue(t, "{}", "1 if 0")
}
