package syntax

import "testing"

func TestFilter(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{42}`, `{1, [2, 3], 4, [5, 6]} filter . {[a, b]: 42}`)
	AssertCodesEvalToSameValue(t, `{5, 11}`, `{1, [2, 3], 4, [5, 6]} filter . {[a, b]: a + b}`)
	AssertCodesEvalToSameValue(t, `{5, 11, 16}`, `{1, [2, 3], 4, [5, 6], [7, 8, 9]} filter . {[a, ..., b]: a + b}`)
}
