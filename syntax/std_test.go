package syntax

import "testing"

func TestStdTuple(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `()`, `//.tuple({})`)
	AssertCodesEvalToSameValue(t, `(a:1)`, `//.tuple({"a":1})`)
	AssertCodesEvalToSameValue(t, `(a:1, b:2)`, `//.tuple({"a":1, "b":2})`)
}

func TestStdDict(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{}`, `//.dict(())`)
	AssertCodesEvalToSameValue(t, `{"a":1}`, `//.dict((a:1))`)
	AssertCodesEvalToSameValue(t, `{"a":1, "b":2}`, `//.dict((a:1, b:2))`)
}
