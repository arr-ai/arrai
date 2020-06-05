package syntax

import "testing"

func TestExprDotOnTupleSet(t *testing.T) {
	AssertCodesEvalToSameValue(t, `{(foo: 42)}.foo`, `42`)
}

func TestExprDotErrorOnEmptySet(t *testing.T) {
	AssertCodeErrors(t, `{}.foo`, `Cannot get attr "foo" from empty set`)
}

func TestExprDotErrorOnMultipletSet(t *testing.T) {
	AssertCodeErrors(t, `{1, 2, 3}.foo`, `Too many elts to get attr "foo" from set`)
}

func TestExprDotErrorOnNonTupleSet(t *testing.T) {
	AssertCodeErrors(t, `{1}.foo`, `Cannot get attr "foo" from non-tuple set elt`)
}
