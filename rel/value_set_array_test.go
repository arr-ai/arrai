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

func TestArrayWithout(t *testing.T) {
	three := NewArray(NewNumber(10), NewNumber(11), NewNumber(12))

	AssertEqualValues(t,
		NewOffsetArray(1, NewNumber(11), NewNumber(12)),
		three.Without(NewArrayItemTuple(0, NewNumber(10))),
	)
	AssertExprEvalsToType(t,
		Array{},
		three.Without(NewArrayItemTuple(0, NewNumber(10))),
	)

	AssertEqualValues(t,
		NewArray(NewNumber(10), NewNumber(11)),
		three.Without(NewArrayItemTuple(2, NewNumber(12))),
	)
	AssertExprEvalsToType(t,
		Array{},
		three.Without(NewArrayItemTuple(2, NewNumber(12))),
	)

	four := NewArray(NewNumber(10), NewNumber(11), NewNumber(12), NewNumber(13))

	AssertEqualValues(t,
		NewOffsetArray(1, NewNumber(11), NewNumber(12)),
		four.Without(NewArrayItemTuple(3, NewNumber(13))).Without(NewArrayItemTuple(0, NewNumber(10))),
	)
	AssertEqualValues(t,
		NewOffsetArray(1, NewNumber(11), NewNumber(12)),
		four.Without(NewArrayItemTuple(0, NewNumber(10))).Without(NewArrayItemTuple(3, NewNumber(13))),
	)
}

func TestArrayCallAll(t *testing.T) {
	t.Parallel()

	three := NewArray(NewNumber(10), NewNumber(11), NewNumber(12))

	AssertEqualValues(t, NewSet(NewNumber(10)), three.CallAll(NewNumber(0)))
	AssertEqualValues(t, NewSet(NewNumber(11)), three.CallAll(NewNumber(1)))
	AssertEqualValues(t, NewSet(NewNumber(12)), three.CallAll(NewNumber(2)))
	AssertEqualValues(t, None, three.CallAll(NewNumber(5)))
	AssertEqualValues(t, None, three.CallAll(NewNumber(-1)))

	three = NewOffsetArray(-2, NewNumber(10), NewNumber(11), NewNumber(12))
	AssertEqualValues(t, NewSet(NewNumber(10)), three.CallAll(NewNumber(-2)))
	AssertEqualValues(t, NewSet(NewNumber(11)), three.CallAll(NewNumber(-1)))
	AssertEqualValues(t, NewSet(NewNumber(12)), three.CallAll(NewNumber(0)))
	AssertEqualValues(t, None, three.CallAll(NewNumber(1)))
	AssertEqualValues(t, None, three.CallAll(NewNumber(-3)))

	three = NewOffsetArray(2, NewNumber(10), NewNumber(11), NewNumber(12))
	AssertEqualValues(t, NewSet(NewNumber(10)), three.CallAll(NewNumber(2)))
	AssertEqualValues(t, NewSet(NewNumber(11)), three.CallAll(NewNumber(3)))
	AssertEqualValues(t, NewSet(NewNumber(12)), three.CallAll(NewNumber(4)))
	AssertEqualValues(t, None, three.CallAll(NewNumber(1)))
	AssertEqualValues(t, None, three.CallAll(NewNumber(5)))
}
