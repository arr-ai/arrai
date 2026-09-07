package syntax

import "testing"

// The dict grammar's `pairs` rule shares its "extra" (...) alternative with dict
// patterns, so it parses in a dict expression too, but compileDictEntryExprs never
// handled it there and crashed with a raw type-assertion panic instead of a compile
// error.
func TestDictExtraElementIsExpressionError(t *testing.T) {
	t.Parallel()
	AssertCodeErrors(t,
		"extra element (...) is only valid in a dict pattern, not a dict expression",
		`let d = {'a': 1}; {...d, 'b': 2}`)
}
