package syntax

import "testing"

func TestOrder(t *testing.T) {
	AssertCodesEvalToSameValue(t, `[-3, -2, -1, 0, 1, 2, 3]`, `{3, 2, 1, 0, -1, -2, -3} order \a \b  a < b`)
	AssertCodesEvalToSameValue(t, `[3, 2, 1, 0, -1, -2, -3]`, `{3, 2, 1, 0, -1, -2, -3} order \a \b  a > b`)
}

func TestOrderBy(t *testing.T) {
	AssertCodesEvalToSameValue(t, `[-3, -2, -1, 0, 1, 2, 3, 4]`, `{3, 1, -1, 4, 0, -3, -2, 2} orderby .`)
	AssertCodesEvalToSameValue(t, `[4, 3, 2, 1, 0, -1, -2, -3]`, `{3, 1, -1, 4, 0, -3, -2, 2} orderby -.`)
	// TODO: Deal with non-deterministic order.
	// AssertCodesEvalToSameValue(t, `[-3, -2, -1, 0, 1, 2, 3, 4]`, `{3, 1, -1, 4, 0, -3, -2, 2} orderby . ^ 2`)
}

func TestRank(t *testing.T) {
	AssertCodesEvalToSameValue(t, `{|x,r| (1,0), (2,1), (3,2)}`, `{|x| (1), (2), (3)} rank (r: .x)`)
	AssertCodesEvalToSameValue(t, `{|x,r| (1,2), (2,1), (3,0)}`, `{|x| (1), (2), (3)} rank (r: -.x)`)
	AssertCodesEvalToSameValue(t, `{|x,r,s| (1,0,2), (2,1,1), (3,2,0)}`, `{|x| (1), (2), (3)} rank (r: .x, s: -.x)`)
}
