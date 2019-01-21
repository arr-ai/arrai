package rel

import (
	"testing"
)

// TestTrivialIntersect test Intersect({}, {}).
func TestTrivialIntersect(t *testing.T) {
	a := intSet()
	b := intSet()
	AssertEqualValues(t, a, Intersect(a, b))
}

// TestOneSidedIntersect tests Intersect({}, {a:42, b:43}) and vice-versa.
func TestOneSidedIntersect(t *testing.T) {
	a := intSet()
	b := intSet(42, 43)
	AssertEqualValues(t, a, Intersect(a, b))
	AssertEqualValues(t, a, Intersect(b, a))
}

// TestEqualIntersect tests Intersect({a:42, b:43}, {a:42, b:43}).
func TestEqualIntersect(t *testing.T) {
	a := intSet(42, 43)
	b := intSet(42, 43)
	AssertEqualValues(t, a, Intersect(a, b))
}

// TestMixedIntersect tests Intersect({a:42, b:43}, {b:43, c:44}).
func TestMixedIntersect(t *testing.T) {
	a := intSet(42, 43)
	b := intSet(43, NewNumber(44))
	c := intSet(43)
	AssertEqualValues(t, c, Intersect(a, b))
}

// TestTrivialUnion test Union({}, {}).
func TestTrivialUnion(t *testing.T) {
	a := intSet()
	b := intSet()
	AssertEqualValues(t, a, Union(a, b))
}

// TestOneSidedUnion tests Union({}, {a:42, b:43}) and vice-versa.
func TestOneSidedUnion(t *testing.T) {
	a := intSet()
	b := intSet(42, 43)
	AssertEqualValues(t, b, Union(a, b))
	AssertEqualValues(t, b, Union(b, a))
}

// TestEqualUnion tests Union({a:42, b:43}, {a:42, b:43}).
func TestEqualUnion(t *testing.T) {
	a := intSet(42, 43)
	b := intSet(42, 43)
	AssertEqualValues(t, a, Union(a, b))
}

// TestMixedUnion tests Union({a:42, b:43}, {b:43, c:44}).
func TestMixedUnion(t *testing.T) {
	a := intSet(42, 43)
	b := intSet(43, 44)
	c := intSet(42, 43, 44)
	AssertEqualValues(t, c, Union(a, b))
}

// TestTrivialDifference test Difference({}, {}).
func TestTrivialDifference(t *testing.T) {
	a := intSet()
	b := intSet()
	AssertEqualValues(t, a, Difference(a, b))
}

// TestOneSidedDifference tests Difference({}, {a:42, b:43}) and vice-versa.
func TestOneSidedDifference(t *testing.T) {
	a := intSet()
	b := intSet(42, 43)
	AssertEqualValues(t, a, Difference(a, b))
	AssertEqualValues(t, b, Difference(b, a))
}

// TestEqualDifference tests Difference({a:42, b:43}, {a:42, b:43}).
func TestEqualDifference(t *testing.T) {
	a := intSet(42, 43)
	b := intSet(42, 43)
	AssertEqualValues(t, intSet(), Difference(a, b))
}

// TestMixedDifference tests Difference({a:42, b:43}, {b:43, c:44}).
func TestMixedDifference(t *testing.T) {
	a := intSet(42, 43)
	b := intSet(43, 44)
	AssertEqualValues(t, intSet(42), Difference(a, b))
	AssertEqualValues(t, intSet(44), Difference(b, a))
}
