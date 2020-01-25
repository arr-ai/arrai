package syntax

import "testing"

func TestParseIfElseExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `42 if true else 43`)
	AssertCodesEvalToSameValue(t, `43`, `42 if (false) else 43`)
}

func TestParseTupleExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `(a:42).a`)
	// TODO: Fix
	// AssertCodesEvalToSameValue(t, `{a:42}`, `(\a{a})42`)
}

func TestParseTupleShorthand(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `(a:42)`, `(\a(a:a))(42)`)
	// TODO: Fix
	// AssertCodesEvalToSameValue(t, `{a:42}`, `(\a{a})42`)
}

func TestParseCurry(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `(\x \y x * y)(6)(7)`)
	// TODO: Fix
	// AssertCodesEvalToSameValue(t, `{a:42}`, `(\a{a})42`)
}

func TestParseCurry2(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `(\op \x \y op(x, y))(\a \b a * b)(6)(7)`)
	// TODO: Fix
	// AssertCodesEvalToSameValue(t, `{a:42}`, `(\a{a})42`)
}

func TestParseFix(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `720`,
		`(\f f(f))(\f \g \n g(f(f)(g))(n))(\fact \n 1 if n < 2 else fact(n - 1) * n)(6)`)
	AssertCodesEvalToSameValue(t, `2`,
		`(\f f(f))(\f \g \n g(f(f)(g))(n))(\gcd \a \b a if b = 0 else gcd(b)(a % b))(20)(14)`)
	AssertCodesEvalToSameValue(t, `2`,
		`(\f f(f))(\f \g \n g(f(f)(g))(n))(\gcd \a \b a if b = 0 else gcd(b, a % b))(20, 14)`)
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
	t.Parallel()

	// TODO: Fix
	// AssertCodesEvalToSameValue(t,
	// 	`{a:{b:41,c:42}}`,
	// 	`({a:{b:41}}->*a with {c:42})`)
	// AssertCodesEvalToSameValue(t, `{a:{b:42}}`, `{a:{b:41}}->*a->*b: 42`)
	// AssertCodesEvalToSameValue(t,
	// 	`{a:{b:{c:42}}}`,
	// 	`{a:{b:{c:41}}}->*a->*b->*c: 42`)
}
