package syntax

import "testing"

func TestNest(t *testing.T) {
	AssertCodesEvalToSameValue(t,
		`{{(a: 1, z: {(b: 3, c: 5)})}, {(a: 2, z: {(b: 4, c: 6)})}}`,
		`{{(a: 1, b: 3, c: 5)}, {(a: 2, b: 4, c: 6)}} => ((.) nest ~|a|z)`,
	)
}

func TestInverseNest(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`{(a: {(x: 1, y: 1)}, z: 2), (a: {(x: 1, y: 1)}, z: 3), (a: {(x: 1, y: 2)}, z: 4), (a: {(x: 1, y: 3)}, z: 5)}`,
		`{|x, y, z| (1, 1, 2), (1, 1, 3), (1, 2, 4), (1, 3, 5)} nest ~|z|a`,
	)
	AssertCodesEvalToSameValue(t,
		`{(a: {(z: 2), (z: 3)}, x: 1, y: 1), (a: {(z: 4)}, x: 1, y: 2), (a: {(z: 5)}, x: 1, y: 3)}`,
		`{|x, y, z| (1, 1, 2), (1, 1, 3), (1, 2, 4), (1, 3, 5)} nest ~|x, y|a`,
	)
	AssertCodeErrors(t,
		`nest attrs cannot be on all of relation attrs (|x, y, z|)`,
		`{|x, y, z| (1, 1, 2), (1, 1, 3), (1, 2, 4), (1, 3, 5)} nest ~|x, y, z|a`,
	)
	AssertCodeErrors(t,
		`nest attrs (|b|) not a subset of relation attrs (|x, y, z|)`,
		`{|x, y, z| (1, 1, 2), (1, 1, 3), (1, 2, 4), (1, 3, 5)} nest ~|b|a`,
	)
	AssertCodeErrors(t,
		"not a relation; inconsistent attribute names between tuples",
		`{(x: 1, y: 1), (x: 1, y: 2, z: 3)} nest ~|x|a`)
}
