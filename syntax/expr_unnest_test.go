package syntax

import "testing"

func TestUnnest(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`{|x, y, z| (1, 1, 2), (1, 1, 3), (1, 2, 4), (1, 3, 5)}`,
		`{
			(a: {(x: 1, y: 1)}, z: 2),
			(a: {(x: 1, y: 1)}, z: 3),
			(a: {(x: 1, y: 2)}, z: 4),
			(a: {(x: 1, y: 3)}, z: 5)
		} unnest a`,
	)
}
