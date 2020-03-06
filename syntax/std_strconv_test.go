package syntax

import "testing"

func TestStrconvEval(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `123             `, `//.strconv.eval("123")             `)
	AssertCodesEvalToSameValue(t, `true            `, `//.strconv.eval("true")            `)
	AssertCodesEvalToSameValue(t, `123.321         `, `//.strconv.eval("123.321")         `)
	AssertCodesEvalToSameValue(t, `"this is a test"`, `//.strconv.eval("'this is a test'")`)
	AssertCodesEvalToSameValue(t,
		`(str: "stuff", num:123, array: [1,2,3])`,
		`//.strconv.eval("(str: 'stuff', num:123, array: [1,2,3])")`)
	assertExprPanics(t, `//.strconv.eval(123)`)
}

func TestStrconvUnsafeEval(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `123             `, `//.strconv.unsafe_eval("123")             `)
	AssertCodesEvalToSameValue(t, `true            `, `//.strconv.unsafe_eval("true")            `)
	AssertCodesEvalToSameValue(t, `123.321         `, `//.strconv.unsafe_eval("123.321")         `)
	AssertCodesEvalToSameValue(t, `"this is a test"`, `//.strconv.unsafe_eval("'this is a test'")`)
	AssertCodesEvalToSameValue(t, `12300           `, `//.strconv.unsafe_eval("123*100")         `)
	AssertCodesEvalToSameValue(t, `6               `, `//.strconv.unsafe_eval("(\\x x + 1)(5)")  `)
	AssertCodesEvalToSameValue(t, `6               `, `//.strconv.unsafe_eval("[2, 4, 6, 8](2)") `)
	AssertCodesEvalToSameValue(t,
		`(str: "stuff", num:123, array: [1,2,3])`,
		`//.strconv.unsafe_eval("(str: 'stuff', num:123, array: [1,2,3])")`)
	assertExprPanics(t, `//.strconv.unsafe_eval(123*345)`)
}
