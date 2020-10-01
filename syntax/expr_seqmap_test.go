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
	AssertCodesEvalToSameValue(t, `{'a': 2, 'b': 4, 'c': 6}`, `{'a':1,'b':2,'c':3} >> \v v * 2`)
}

func TestExprSequenceMapIdentHoles(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `[1, , , 2]`, `[1,,,2] >> .`)
}

func TestExprSequenceMapNonSeq(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, `>> lhs must be an indexed type, not set`, `{1,2,3} >> .`)
}
