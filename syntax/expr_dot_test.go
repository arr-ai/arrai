package syntax

import "testing"

func TestExprDotOnTupleSet(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{(foo: 42)}.foo`, `42`)
}

func TestExprDotErrorOnEmptySet(t *testing.T) {
	t.Parallel()
	AssertCodeErrors(t, `{}.foo`, `Cannot get attr "foo" from empty set`)
}

func TestExprDotErrorOnMultipletSet(t *testing.T) {
	t.Parallel()
	AssertCodeErrors(t, `{1, 2, 3}.foo`, `Too many elts to get attr "foo" from set`)
}

func TestExprDotErrorOnNonTupleSet(t *testing.T) {
	t.Parallel()
	AssertCodeErrors(t, `{1}.foo`, `Cannot get attr "foo" from non-tuple set elt`)
}
