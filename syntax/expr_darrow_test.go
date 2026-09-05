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

// 🎯T25 regression: the injective projectDots fast path renamed a Relation's
// columns to @/@item without canonicalising the result, so `set => (@: .i,
// @item: .v)` stayed a Relation instead of becoming an Array.
func TestDArrowExprProjectToArrayItemShape(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`['a']`,
		`{(index: 0, val: 'a')} => (@: .index, @item: .val)`,
	)
	AssertCodesEvalToSameValue(t,
		`['a', 'b']`,
		`{(index: 0, val: 'a'), (index: 1, val: 'b')} => (@: .index, @item: .val)`,
	)
}
