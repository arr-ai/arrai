package syntax

import "testing"

func TestDArrowExprEmpty(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{}`, `{} => .`)
}

func TestDArrowExprIdent(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{1, 2, 3}`, `{1,2,3} => .`)
}

func TestDArrowExprDouble(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{2, 4, 6}`, `{1,2,3} => \i i * 2`)
}

func TestDArrowExprIdentHoles(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{1, , , 2}`, `{1,,,2} => .`)
}
