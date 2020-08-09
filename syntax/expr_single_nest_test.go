package syntax

import "testing"

func TestSingleAttrNest(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`{(x: 1, y: {2, 3}), (x: 2, y: {4})}`,
		`{|x,y| (1, 2), (1, 3), (2, 4)} nest y`,
	)

	AssertCodesEvalToSameValue(t,
		`{(x: 1, y: {2, 3, 4})}`,
		`{|x,y| (1, 2), (1, 3), (1, 4)} nest y`,
	)

	AssertCodesEvalToSameValue(t,
		`{(x: 1, y: {2}), (x: 2, y: {3}), (x: 3, y: {4})}`,
		`{|x,y| (1, 2), (2, 3), (3, 4)} nest y`,
	)

	AssertCodesEvalToSameValue(t, `{}`, `{} nest y`)
}
