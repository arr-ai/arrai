package rel

import "testing"

func TestAsArray(t *testing.T) {
	AssertEqualValues(t,
		NewArray(NewNumber(10), NewNumber(11)),
		NewSet(
			NewArrayItemTuple(0, NewNumber(10)),
			NewArrayItemTuple(1, NewNumber(11)),
		),
	)
}
