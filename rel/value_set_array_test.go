package rel

import "testing"

func TestAsArray(t *testing.T) {
	AssertEqualValues(t,
		NewArray(NewNumber(10), NewNumber(11)),
		NewSet(
			newArrayTuple(0, NewNumber(10)),
			newArrayTuple(1, NewNumber(11)),
		),
	)
}
