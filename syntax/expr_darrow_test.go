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

// `set => (@: .i, @item: .v)` — the idiom for turning a set into an array —
// must produce a real Array, not a two-column Relation.
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
