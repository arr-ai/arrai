package syntax

import "testing"

func TestArrayExprStringEmpty(t *testing.T) {
	t.Parallel()
	AssertEvalExprString(t, `{}`, `[]`)
}

func TestArrayExprStringHoles(t *testing.T) {
	t.Parallel()
	AssertEvalExprString(t, `[1, , ]`, `[1,,]`)
	AssertEvalExprString(t, `[1, , 2]`, `[1,,2]`)
	AssertEvalExprString(t, `[1, , , 2]`, `[1,,,2]`)
}
