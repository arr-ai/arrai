package syntax

import "testing"

func TestStdTuple(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `()`, `//tuple({})`)
	AssertCodesEvalToSameValue(t, `(a:1)`, `//tuple({"a":1})`)
	AssertCodesEvalToSameValue(t, `(a:1, b:2)`, `//tuple({"a":1, "b":2})`)

	AssertCodePanics(t, `//tuple((a:1))`)
	AssertCodePanics(t, `//tuple(42)`)
}

func TestStdDict(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{}`, `//dict(())`)
	AssertCodesEvalToSameValue(t, `{"a":1}`, `//dict((a:1))`)
	AssertCodesEvalToSameValue(t, `{"a":1, "b":2}`, `//dict((a:1, b:2))`)

	AssertCodePanics(t, `//dict({42:43})`)
	AssertCodePanics(t, `//dict(42)`)
}
