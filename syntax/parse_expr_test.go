package syntax

import (
	"fmt"
	"testing"
)

func TestParseTupleExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `(a:42).a`)
	// TODO: Fix
	// AssertCodesEvalToSameValue(t, `{a:42}`, `(\a{a})42`)
}

func TestParseTupleShorthand(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `(a:42)`, `(\a(a:a))(42)`)
}

func TestParseCurry(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `(\x \y x * y)(6)(7)`)
}

func TestParseCurry2(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `(\op \x \y op(x, y))(\a \b a * b)(6)(7)`)
}

func TestParseApply(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `42 -> .`)
	AssertCodesEvalToSameValue(t, `42`, `6 * 7 -> \x x`)
	AssertCodesEvalToSameValue(t, `42`, `7 -> \y (6 -> \x x * y)`)
	AssertCodesEvalToSameValue(t, `23.140692632779263`, `//math -> (//math.pi -> \pi .e^pi)`)
	AssertCodesEvalToSameValue(t, `0`, `//math -> \m (m.sin(0) -> .)`)
}

func TestParseFix(t *testing.T) {
	t.Parallel()

	fix := `(\f f(f))(\f \g \n g(f(f)(g))(n))`
	AssertCodesEvalToSameValue(t, `720`, fix+`(\fact \n cond {n < 2: 1, _: fact(n - 1) * n}) (6)`)
	AssertCodesEvalToSameValue(t, `2`, fix+`(\gcd \a \b cond {b = 0: a, _: gcd(b)(a % b)}) (20)(14)`)
	AssertCodesEvalToSameValue(t, `2`, fix+`(\gcd \a \b cond {b = 0: a, _: gcd(b, a % b)}) (20, 14)`)

	AssertCodesEvalToSameValue(t, `720`, `//fn.fix(\fact \n cond {n < 2: 1, _: fact(n - 1) * n})(6)`)
	AssertCodesEvalToSameValue(t, `2`, `//fn.fix(\gcd \a \b cond {b = 0: a, _: gcd(b)(a % b)})(20)(14)`)
	AssertCodesEvalToSameValue(t, `2`, `//fn.fix(\gcd \a \b cond {b = 0: a, _: gcd(b, a % b)})(20, 14)`)
}

func TestParseFixt(t *testing.T) {
	t.Parallel()
	eo := `(
		even: \t \n n = 0 || t.odd(n - 1),
		odd:  \t \n n != 0 && t.even(n - 1)
	)`
	// AssertCodesEvalToSameValue(t, `true`,
	// 	`(\f f(f))(\f \t t :> \g \n g(f(f)(t))(n)) (`+eo+`)`)
	for _, fixt := range []string{
		`(\f f(f))(\f \t t :> \g \n g(f(f)(t))(n))`,
		`//fn.fixt`,
	} {
		for i := 0; i < 5; i++ {
			for _, p := range []struct {
				m int
				f string
			}{
				{0, "even"},
				{1, "odd"},
			} {
				AssertCodesEvalToSameValue(t,
					fmt.Sprintf("%v", i%2 == p.m),
					fmt.Sprintf(`(%s(`+eo+`)).%s(%d)`, fixt, p.f, i))
			}
		}
	}
}

// func TestParseArrowStarExprAssignExisting(t *testing.T) {
// 	t.Parallel()
// 	AssertCodesEvalToSameValue(t, `(a:42)`, `(a:41)->*a: 42`)
// 	AssertCodesEvalToSameValue(t, `(a:(b:42))`, `(a:(b:41))->*a->*b: 42`)
// 	AssertCodesEvalToSameValue(t,
// 		`(a:(b:(c:42)))`,
// 		`(a:(b:(c:41)))->*a->*b->*c: 42`)
// }

// func TestParseArrowStarExprAssignNew(t *testing.T) {
// 	t.Parallel()
// 	AssertCodesEvalToSameValue(t, `(a:41,b:42)`, `(a:41)->*b: 42`)
// 	AssertCodesEvalToSameValue(t,
// 		`(a:(b:41, c:42))`,
// 		`(a:(b:41))->*a->*c: 42`)
// 	AssertCodesEvalToSameValue(t,
// 		`(a:(b:(c:41,d:42)))`,
// 		`(a:(b:(c:41)))->*a->*b->*d: 42`)
// }

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

func TestParseNestExpr(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{|a,b| (1, {|b| (1), (2)})}`, `{|a,b| (1, 1), (1, 2)} nest |b|b`)
	AssertCodesEvalToSameValue(t,
		`{|a,b| (1, {|b| (1), (2)}), (2, {|b| (3)})}`,
		`{|a,b| (1,1), (1,2), (2,3)} nest |b|b`,
	)
	AssertCodesEvalToSameValue(t,
		`{|a,bc| (1, {|b,c| (1, 1), (2, 1)}), (2, {|b,c| (3, 4)})}`,
		`{|a,b,c| (1,1,1), (2, 3, 4), (1,2,1)} nest |b,c|bc`,
	)
	AssertCodesEvalToSameValue(t,
		`{|a,b,c| (1, 2, {|c| (1), (2)})}`,
		`{|a,b,c| (1,2,1), (1,2,2)} nest |c|c`,
	)
	AssertCodesEvalToSameValue(t, `{}`, `{} nest |a|a`)
	AssertCodesEvalToSameValue(t, `{}`, `{} nest ~|a|a`)
}

func TestParseWhereExpr(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{1, {{}}, (a:3)}`, `{0, 1, {}, {{}}, (a:3), ()} where .`)
	AssertCodesEvalToSameValue(t, `{2}`, `{1, 2, 3} where .%2 = 0`)
	AssertCodesEvalToSameValue(t, `{1, 3}`, `{1, 2, 3} where .%2 = 1`)
}

func TestParseMapExpr(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{0, 1, {}, (a:3)}`, `{0, 1, {}, (a:3)} => .`)
	AssertCodesEvalToSameValue(t, `{true, false}`, `{1, 2, 3} => .%2 = 0`)
	AssertCodesEvalToSameValue(t, `{true, false}`, `{1, 2, 3} => .%2 = 1`)
	AssertCodesEvalToSameValue(t, `{1, 4, 9}`, `{1, 2, 3} => .^2`)
	AssertCodesEvalToSameValue(t, `{2, 4, 8}`, `{1, 2, 3} => 2^.`)
	AssertCodesEvalToSameValue(t, `{0, 1}`, `{1, 2, 3} => .//2`)
}

func TestParseTupleMapExpr(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `(a: 0, b: 1, c: {}, d: (a:3))`, `(a: 0, b: 1, c: {}, d: (a:3)) :> .`)
	AssertCodesEvalToSameValue(t, `(a: false, b: true, c: false)`, `(a: 1, b: 2, c: 3) :> .%2 = 0`)
	AssertCodesEvalToSameValue(t, `(a: true, b: false, c: true)`, `(a: 1, b: 2, c: 3) :> .%2 = 1`)
	AssertCodesEvalToSameValue(t, `(a: 1, b: 4, c: 9)`, `(a: 1, b: 2, c: 3) :> .^2`)
	AssertCodesEvalToSameValue(t, `(a: 2, b: 4, c: 8)`, `(a: 1, b: 2, c: 3) :> 2^.`)
	AssertCodesEvalToSameValue(t, `(a: 0, b: 1, c: 1)`, `(a: 1, b: 2, c: 3) :> .//2`)
}
