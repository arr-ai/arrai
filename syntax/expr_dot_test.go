package syntax

import "testing"

func TestExprDotOnTupleSet(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{(foo: 42)}.foo`, `42`)
}

func TestExprDotErrorOnEmptySet(t *testing.T) {
	t.Parallel()
	AssertCodeErrors(t, `Cannot get attr "foo" from empty set`, `{}.foo`)
}

func TestExprDotErrorOnMultipletSet(t *testing.T) {
	t.Parallel()
	AssertCodeErrors(t, `Too many elts to get attr "foo" from set`, `{1, 2, 3}.foo`)
}

func TestExprDotErrorOnNonTupleSet(t *testing.T) {
	t.Parallel()
	AssertCodeErrors(t, `Cannot get attr "foo" from non-tuple set elt`, `{1}.foo`)
}
