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
