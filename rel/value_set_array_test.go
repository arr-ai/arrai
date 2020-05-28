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

func TestArrayWhere(t *testing.T) {
	t.Parallel()

	three := NewArray(NewNumber(10), NewNumber(11), NewNumber(12))

	atBetween := func(a, b int) func(v Value) bool {
		return func(v Value) bool {
			i := int(v.(ArrayItemTuple).MustGet("@").(Number).Float64())
			return a <= i && i <= b
		}
	}

	AssertEqualValues(t, three, three.Where(atBetween(0, 2)))
	AssertEqualValues(t, NewArray(NewNumber(10), NewNumber(11)), three.Where(atBetween(0, 1)))
	AssertEqualValues(t, NewArray(NewNumber(10)), three.Where(atBetween(0, 0)))
	AssertEqualValues(t, None, three.Where(atBetween(-1, -1)))

	AssertEqualValues(t, None, three.Where(atBetween(3, 3)))
	AssertEqualValues(t, NewOffsetArray(2, NewNumber(12)), three.Where(atBetween(2, 3)))
	AssertEqualValues(t, NewOffsetArray(1, NewNumber(11), NewNumber(12)), three.Where(atBetween(1, 3)))
	AssertEqualValues(t, three, three.Where(atBetween(0, 3)))

	offsetThree := NewOffsetArray(-2, NewNumber(10), NewNumber(11), NewNumber(12))

	AssertEqualValues(t, offsetThree, offsetThree.Where(atBetween(-2, 0)))
	AssertEqualValues(t, NewOffsetArray(-2, NewNumber(10), NewNumber(11)), offsetThree.Where(atBetween(-2, -1)))
	AssertEqualValues(t, NewOffsetArray(-2, NewNumber(10)), offsetThree.Where(atBetween(-2, -2)))
	AssertEqualValues(t, None, offsetThree.Where(atBetween(-3, -3)))

	AssertEqualValues(t, None, offsetThree.Where(atBetween(1, 1)))
	AssertEqualValues(t, NewArray(NewNumber(12)), offsetThree.Where(atBetween(0, 1)))
	AssertEqualValues(t, NewOffsetArray(-1, NewNumber(11), NewNumber(12)), offsetThree.Where(atBetween(-1, 1)))
	AssertEqualValues(t, offsetThree, offsetThree.Where(atBetween(-2, 1)))
}
