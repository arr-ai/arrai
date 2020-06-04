package syntax

import "testing"

func TestExprSequenceMapEmpty(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `[]`, `[] >> .`)
}

func TestExprSequenceMapIdent(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `[1, 2, 3]`, `[1,2,3] >> .`)
}

func TestExprSequenceMapDouble(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `[2, 4, 6]`, `[1,2,3] >> \i i * 2`)
}

func TestExprSequenceMapIdentHoles(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `[1, , , 2]`, `[1,,,2] >> .`)
}
