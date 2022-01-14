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
		"not a relation",
		`{(x: 1, y: 1), (x: 1, y: 2, z: 3)} nest ~|x|a`)
}

func TestDeterministicComplementNest(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`{{(a: 1, z: {(b: 5, c: 10)})}, {(a: 2, z: {(b: 7, c: 11)})}}`,
		`{{(a: 1, b: 5, c: 10)}, {(a: 2, b: 7, c: 11)}} => ((.) nest ~|a|z)`,
	)
	AssertCodesEvalToSameValue(t,
		`
			{
				|a , nona|
				(11, {(b: 12, nonb: {(c: 13)})}),
				(21, {(b: 22, nonb: {(c: 23)})}),
				(31, {(b: 22, nonb: {(c: 33)})}),
				(41, {(b: 22, nonb: {(c: 43)})}),
			}
		`,
		`
			let db = {
				(a: 11, b: 12, c: 13),
				(a: 21, b: 22, c: 23),
				(a: 31, b: 22, c: 33),
				(a: 41, b: 22, c: 43),
			};
			db nest ~|a|nona => . +> (nona: .nona nest ~|b|nonb)
		`,
	)
}
