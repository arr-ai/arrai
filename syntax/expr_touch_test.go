package syntax

import "testing"

// The `->*` ("touch") grammar production has no implementation anywhere in the
// rel package, so compilePostfixAndTouch used to crash with a raw panic("unfinished")
// instead of a compile error.
func TestTouchIsNotImplementedError(t *testing.T) {
	t.Parallel()
	AssertCodeErrors(t,
		"touch (->*) is not implemented",
		`(a: (b: 1)) ->* a ->* b (2)`)
}
