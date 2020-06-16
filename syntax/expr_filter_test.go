package syntax

import "testing"

func TestFilter(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{5, 11}`, `{1, [2, 3], 4, [5, 6]} filter . {[a, b]: a + b}`)
}
