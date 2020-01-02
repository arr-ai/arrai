package tests

import "testing"

func TestParseIfElseExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `42 if true else 43`)
	AssertCodesEvalToSameValue(t, `43`, `42 if false else 43`)
}

func TestParseTupleShorthand(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `(a:42)`, `(\a(a:a))42`)
	// TODO: Fix
	// AssertCodesEvalToSameValue(t, `{a:42}`, `(\a{a})42`)
}

func TestParseArrowStarExprAssignExisting(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `(a:42)`, `(a:41)->*a: 42`)
	AssertCodesEvalToSameValue(t, `(a:(b:42))`, `(a:(b:41))->*a->*b: 42`)
	AssertCodesEvalToSameValue(t,
		`(a:(b:(c:42)))`,
		`(a:(b:(c:41)))->*a->*b->*c: 42`)
}

func TestParseArrowStarExprAssignNew(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `(a:41,b:42)`, `(a:41)->*b: 42`)
	AssertCodesEvalToSameValue(t,
		`(a:(b:41, c:42))`,
		`(a:(b:41))->*a->*c: 42`)
	AssertCodesEvalToSameValue(t,
		`(a:(b:(c:41,d:42)))`,
		`(a:(b:(c:41)))->*a->*b->*d: 42`)
}

func TestParseArrowStarExprAlterExisting(t *testing.T) {
	// t.Parallel()
	// TODO: Fix
	// AssertCodesEvalToSameValue(t,
	// 	`{a:{b:41,c:42}}`,
	// 	`({a:{b:41}}->*a with {c:42})`)
	// AssertCodesEvalToSameValue(t, `{a:{b:42}}`, `{a:{b:41}}->*a->*b: 42`)
	// AssertCodesEvalToSameValue(t,
	// 	`{a:{b:{c:42}}}`,
	// 	`{a:{b:{c:41}}}->*a->*b->*c: 42`)
}
