package syntax

import (
	"testing"
)

func TestEvalValue(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `123             `, `//eval.value("123")             `)
	AssertCodesEvalToSameValue(t, `true            `, `//eval.value("true")            `)
	AssertCodesEvalToSameValue(t, `123.321         `, `//eval.value("123.321")         `)
	AssertCodesEvalToSameValue(t, `"this is a test"`, `//eval.value("'this is a test'")`)
	AssertCodesEvalToSameValue(t,
		`(str: "stuff", num:123, array: [1,2,3])`,
		`//eval.value("(str: 'stuff', num:123, array: [1,2,3])")`)
	AssertCodesEvalToSameValue(t, `123             `, `//eval.value(<<"123">>)             `)
	AssertCodesEvalToSameValue(t, `true            `, `//eval.value(<<"true">>)            `)
	AssertCodesEvalToSameValue(t, `123.321         `, `//eval.value(<<"123.321">>)         `)
	AssertCodesEvalToSameValue(t, `"this is a test"`, `//eval.value(<<"'this is a test'">>)`)
	AssertCodesEvalToSameValue(t,
		`(str: "stuff", num:123, array: [1,2,3])`,
		`//eval.value(<<"(str: 'stuff', num:123, array: [1,2,3])">>)`)
	AssertCodeErrors(t, "", `//eval.value(123)`)
}
func TestEvalEval(t *testing.T) {
	t.Parallel()
	AssertCodeErrors(t, "", `//eval.eval("//os.file('Makefile')")`)
	AssertCodeErrors(t,
		`Missing attr "upper" (available: |lower|)`,
		`//eval.evaluator((stdlib: (str: (lower: //str.lower)))).eval("//str.upper('cat')")`)
	AssertCodesEvalToSameValue(t,
		`"cat"`,
		`//eval.evaluator((stdlib: (str: (lower: //str.lower)))).eval("//str.lower('CAT')")`)
	AssertCodesEvalToSameValue(t,
		`"cat"`,
		`//eval.evaluator(()).eval("//str.lower('CAT')")`)
	AssertCodesEvalToSameValue(t, `123`, `//eval.eval("123")`)
	AssertCodesEvalToSameValue(t, `"cat"`, `//eval.eval("//str.lower('CAT')")`)
	AssertCodesEvalToSameValue(t, `3`, `let x = 1; let y = 2; //eval.evaluator((scope: (:x, :y))).eval("x + y") `)
	AssertCodesEvalToSameValue(t, `6`, `let double = \d d * 2; //eval.evaluator((scope: (:double))).eval("double(1 + 2)")`)
}
