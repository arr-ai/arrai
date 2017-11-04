package tests

import "testing"

// TestParseIfElseExpr tests Parse recognising `a if b else c`.
func TestParseIfElseExpr(t *testing.T) {
	assertCodesEvalToSameValue(t, `42`, `42 if true else 43`)
	assertCodesEvalToSameValue(t, `43`, `42 if false else 43`)
}

func TestParseTupleShorthand(t *testing.T) {
	assertCodesEvalToSameValue(t, `{a:42}`, `(\a{a:a})42`)
	// TODO: Fix
	// assertCodesEvalToSameValue(t, `{a:42}`, `(\a{a})42`)
}

func TestParseArrowStarExprAssignExisting(t *testing.T) {
	assertCodesEvalToSameValue(t, `{a:42}`, `{a:41}->*a: 42`)
	assertCodesEvalToSameValue(t, `{a:{b:42}}`, `{a:{b:41}}->*a->*b: 42`)
	assertCodesEvalToSameValue(t,
		`{a:{b:{c:42}}}`,
		`{a:{b:{c:41}}}->*a->*b->*c: 42`)
}

func TestParseArrowStarExprAssignNew(t *testing.T) {
	assertCodesEvalToSameValue(t, `{a:41,b:42}`, `{a:41}->*b: 42`)
	assertCodesEvalToSameValue(t,
		`{a:{b:41, c:42}}`,
		`{a:{b:41}}->*a->*c: 42`)
	assertCodesEvalToSameValue(t,
		`{a:{b:{c:41,d:42}}}`,
		`{a:{b:{c:41}}}->*a->*b->*d: 42`)
}

func TestParseArrowStarExprAlterExisting(t *testing.T) {
	// TODO: Fix
	// assertCodesEvalToSameValue(t,
	// 	`{a:{b:41,c:42}}`,
	// 	`({a:{b:41}}->*a with {c:42})`)
	// assertCodesEvalToSameValue(t, `{a:{b:42}}`, `{a:{b:41}}->*a->*b: 42`)
	// assertCodesEvalToSameValue(t,
	// 	`{a:{b:{c:42}}}`,
	// 	`{a:{b:{c:41}}}->*a->*b->*c: 42`)
}
