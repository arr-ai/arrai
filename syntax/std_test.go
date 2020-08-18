package syntax

import "testing"

func TestStdTuple(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `()`, `//tuple({})`)
	AssertCodesEvalToSameValue(t, `('':1)`, `//tuple({"":1})`)
	AssertCodesEvalToSameValue(t, `('':1)`, `//tuple({{}:1})`)
	AssertCodesEvalToSameValue(t, `('':1)`, `//tuple({[]:1})`)
	AssertCodesEvalToSameValue(t, `('1':2)`, `//tuple({1:2})`)
	AssertCodesEvalToSameValue(t, `('[1]':2)`, `//tuple({[1]:2})`)
	AssertCodesEvalToSameValue(t, `(a:1)`, `//tuple({"a":1})`)
	AssertCodesEvalToSameValue(t, `(a:1, b:2)`, `//tuple({"a":1, "b":2})`)

	AssertCodeErrors(t, "", `//tuple((a:1))`)
	AssertCodeErrors(t, "", `//tuple(42)`)
}

func TestStdDict(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{}`, `//dict(())`)
	AssertCodesEvalToSameValue(t, `{"a":1}`, `//dict((a:1))`)
	AssertCodesEvalToSameValue(t, `{"a":1, "b":2}`, `//dict((a:1, b:2))`)

	AssertCodeErrors(t, "", `//dict({42:43})`)
	AssertCodeErrors(t, "", `//dict(42)`)
}
