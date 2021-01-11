package syntax

import "testing"

func TestError(t *testing.T) {
	AssertCodeErrors(t, "test", "//error('test')")
	AssertCodeErrors(t, "1", "//error(1)")
	AssertCodeErrors(t, "(a: 1)", "//error((a: 1))")
}
