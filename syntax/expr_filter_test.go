package syntax

import "testing"

func TestFilter(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t, `{42}`, `{1, [2, 3], 4, [5, 6]} filter . {[a, b]: 42}`)
	AssertCodesEvalToSameValue(t, `{42}`, `{1, [2, 3], 4, [5, 6]} filter . {[_, _]: 42}`)
	AssertCodesEvalToSameValue(t, `{5, 11}`, `{1, [2, 3], 4, [5, 6]} filter . {[a, b]: a + b}`)
	AssertCodesEvalToSameValue(t, `{1, 4, 5, 11}`, `{1, [2, 3], 4, [5, 6]} filter . {[a, b]: a + b, k: k}`)
	AssertCodesEvalToSameValue(t, `{5, 11, 16}`, `{1, [2, 3], 4, [5, 6], [7, 8, 9]} filter . {[a, ..., b]: a + b}`)
	AssertCodesEvalToSameValue(t, `{1}`, `let x = 1; {|a, b| (1, 1), (2, 2), (3, 3)} filter . {(a: (x), :b): b}`)

	AssertCodesEvalToSameValue(t,
		`{1, 2, 3, 4}`,
		`//rel.union({[1, 2], {3, 4}} filter . {[...a]: a => .@item, {...a}: a})`)
}
