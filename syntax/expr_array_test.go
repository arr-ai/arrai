package syntax

import "testing"

func TestArrayExprStringEmpty(t *testing.T) {
	AssertEvalExprString(t, `{}`, `[]`)
}

func TestArrayExprStringHoles(t *testing.T) {
	AssertEvalExprString(t, `[1,,]`, `[1,,]`)
	AssertEvalExprString(t, `[1,,2]`, `[1,,2]`)
	AssertEvalExprString(t, `[1,,,2]`, `[1,,,2]`)
}
