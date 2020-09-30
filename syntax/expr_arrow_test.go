package syntax

import "testing"

func TestApplyExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `42`, `6 -> 7 * .`)
	AssertCodesEvalToSameValue(t, `42`, `6 -> \x 7 * x`)

	AssertCodesEvalToSameValue(t, `[1,2,3]`, `[1,2,3] -> (.)`)
	AssertCodesEvalToSameValue(t, `[1,2,3]`, `[1,2,3] -> .`)
}

func TestSeqArrow(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `'HELLO'`, `'hello' >> . - 32`)
	AssertCodesEvalToSameValue(t, `10\'HELLO'`, `10\'hello' >> . - 32`)

	AssertCodesEvalToSameValue(t, `<<'HELLO'>>`, `<<'hello'>> >> . - 32`)
	AssertCodesEvalToSameValue(t, `10\<<'HELLO'>>`, `10\<<'hello'>> >> . - 32`)

	AssertCodesEvalToSameValue(t, `[2, 4, 8]`, `[1, 2, 3] >> 2 ^ .`)
	AssertCodesEvalToSameValue(t, `[2, , 8]`, `[1, , 3] >> 2 ^ .`)
	AssertCodesEvalToSameValue(t, `1\[4, 8]`, `1\[2, 3] >> 2 ^ .`)

	AssertCodesEvalToSameValue(t, `{'a': 42, 'b': 54}`, `{'a': 7, 'b': 9} >> 6 * .`)
}

func TestISeqArrow(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `'ACE'`, `'abc' >>> \i \. . - 32 + i`)
	AssertCodesEvalToSameValue(t, `2\'CEG'`, `2\'abc' >>> \i \. . - 32 + i`)

	AssertCodesEvalToSameValue(t, `<<'ABC'>>`, `<<'ace'>> >>> \i \. . - 32 - i`)
	AssertCodesEvalToSameValue(t, `2\<<'AAA'>>`, `2\<<'cde'>> >>> \i \. . - 32 - i`)

	AssertCodesEvalToSameValue(t, `[2, 8, 32]`, `[1, 2, 3] >>> \i \. 2 ^ (. + i)`)
	AssertCodesEvalToSameValue(t, `[2, , 32]`, `[1, , 3] >>> \i \. 2 ^ (. + i)`)
	AssertCodesEvalToSameValue(t, `1\[8, 32]`, `1\[2, 3] >>> \i \. 2 ^ (. + i)`)

	AssertCodesEvalToSameValue(t, `{1: 42, 2: 54}`, `{1: 6, 2: 7} >>> \i \. 6 * (. + i)`)

	AssertCodesEvalToSameValue(t,
		`{
			3      : ( "key": 3      , "val": (2)      ),
			"ten"  : ( "key": "ten"  , "val": 10       ),
			"stuff": ( "key": "stuff", "val": "random" ),
		}`,
		`{"stuff": "random", "ten": 10, 3: (2)} >>> \i \n ("key": i, "val": n)`,
	)
	AssertCodeErrors(t,
		`>>> lhs must be an indexed type, not set`,
		`{("a": "z"), ("b": "y")} >>> \i \n (i ++ n)`,
	)
}

func TestApplyExprInsideMapExpr(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{2, 4, 8}`, `{1, 2, 3} => (2 -> \y y ^ .)`)
	AssertCodesEvalToSameValue(t, `{2, 4, 8}`, `(\z {1, 2, 3} => (z -> \y y ^ .))(2)`)
}

func TestApplyExprWithPattern(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `2`, `1 -> \x let [(x), y] = [1, 2]; y`)
	AssertCodesEvalToSameValue(t, `3`, `[1, 2] -> \[x, y] x + y`)
	AssertCodesEvalToSameValue(t, `3`, `(m: 1, n: 2) -> \(m: x, n: y) x + y`)
	AssertCodesEvalToSameValue(t, `3`, `{"m": 1, "n": 2} -> \{"m": x, "n": y} x + y`)
	AssertCodesEvalToSameValue(t, `6`, `[1, [2, 3]] -> \[x, [y, z]] x + y + z`)
	AssertCodesEvalToSameValue(t, `[3, [3, 4]]`, `[1, 2, 3, 4] -> \[x, y, ...t] [x + y, t]`)
	AssertCodeErrors(t, "", `1 -> \x let [(x), y] = [2, 2]; y`)
}

func TestUnaryArrows(t *testing.T) {
	t.Parallel()

	// This tests an error when stringifying sets of arrays.
	AssertCodesEvalToSameValue(t, `{{2, 4}, {8, 16}}`, `{{1, 2}, {3, 4}} => => 2 ^ .`)
	AssertCodesEvalToSameValue(t, `{[2, 4], [8, 16]}`, `{[1, 2], [3, 4]} => >> 2 ^ .`)
	AssertCodesEvalToSameValue(t, `{{'a':2, 'b':4}, {'c':8, 'd':16}}`, `{{'a':1, 'b':2}, {'c':3, 'd':4}} => >> 2 ^ .`)
	AssertCodesEvalToSameValue(t, `{(a:2, b:4), (c:8, d:16)}`, `{(a:1, b:2), (c:3, d:4)} => :> 2 ^ .`)

	AssertCodesEvalToSameValue(t, `[{2, 4}, {8, 16}]`, `[{1, 2}, {3, 4}] >> => 2 ^ .`)

	AssertCodesEvalToSameValue(t,
		`(a:{[2, 4], [8]}, b:{[16], [32, 64]})`,
		`(a:{[1, 2], [3]}, b:{[4], [5, 6]}) :> => >> 2 ^ .`)
}
