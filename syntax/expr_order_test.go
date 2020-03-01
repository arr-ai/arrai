package syntax

import "testing"

func TestOrder(t *testing.T) {
	AssertCodesEvalToSameValue(t, `[-3, -2, -1, 0, 1, 2, 3]`, `{3, 2, 1, 0, -1, -2, -3} order \a \b  a < b`)
	AssertCodesEvalToSameValue(t, `[3, 2, 1, 0, -1, -2, -3]`, `{3, 2, 1, 0, -1, -2, -3} order \a \b  a > b`)
}

func TestOrderBy(t *testing.T) {
	AssertCodesEvalToSameValue(t, `[-3, -2, -1, 0, 1, 2, 3, 4]`, `{3, 1, -1, 4, 0, -3, -2, 2} orderby .`)
	AssertCodesEvalToSameValue(t, `[4, 3, 2, 1, 0, -1, -2, -3]`, `{3, 1, -1, 4, 0, -3, -2, 2} orderby -.`)
	// AssertCodesEvalToSameValue(t, `[-3, -2, -1, 0, 1, 2, 3, 4]`, `{3, 1, -1, 4, 0, -3, -2, 2} orderby . ^ 2`)
}
