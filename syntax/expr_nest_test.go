package syntax

import "testing"

func TestInverseNest(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`{(a: {(z: 2), (z: 3)}, x: 1, y: 1), (a: {(z: 4)}, x: 1, y: 2), (a: {(z: 5)}, x: 1, y: 3)}`,
		`{|x, y, z| (1, 1, 2), (1, 1, 3), (1, 2, 4), (1, 3, 5)} nest ~|z|a`,
	)
	AssertCodesEvalToSameValue(t,
		`{(a: {(x: 1, y: 1)}, z: 2), (a: {(x: 1, y: 1)}, z: 3), (a: {(x: 1, y: 2)}, z: 4), (a: {(x: 1, y: 3)}, z: 5)}`,
		`{|x, y, z| (1, 1, 2), (1, 1, 3), (1, 2, 4), (1, 3, 5)} nest ~|x, y|a`,
	)
	AssertCodesEvalToSameValue(t,
		`{(a: {|x, y, z| (1, 1, 2), (1, 1, 3), (1, 2, 4), (1, 3, 5)})}`,
		`{|x, y, z| (1, 1, 2), (1, 1, 3), (1, 2, 4), (1, 3, 5)} nest ~|x, y, z|a`,
	)

	AssertCodePanics(t, `{|x, y, z| (1, 1, 2), (1, 1, 3), (1, 2, 4), (1, 3, 5)} nest ~|b|a`)
	AssertCodePanics(t, `{(x: 1, y: 1), (x: 1, y: 2, z: 3)} nest ~|x|a                    `)
}
