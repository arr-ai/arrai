package rel

import (
	"testing"
)

func TestTrivialIntersect(t *testing.T) {
	t.Parallel()
	a := intSet()
	b := intSet()
	AssertEqualValues(t, a, Intersect(a, b))
}

func TestOneSidedIntersect(t *testing.T) {
	t.Parallel()
	a := intSet()
	b := intSet(42, 43)
	AssertEqualValues(t, a, Intersect(a, b))
	AssertEqualValues(t, a, Intersect(b, a))
}

func TestEqualIntersect(t *testing.T) {
	t.Parallel()
	a := intSet(42, 43)
	b := intSet(42, 43)
	AssertEqualValues(t, a, Intersect(a, b))
}

func TestMixedIntersect(t *testing.T) {
	t.Parallel()
	a := intSet(42, 43)
	b := intSet(43, NewNumber(44))
	c := intSet(43)
	AssertEqualValues(t, c, Intersect(a, b))
}

func TestTrivialNIntersect(t *testing.T) {
	t.Parallel()
	a := intSet()
	b := intSet()
	AssertEqualValues(t, a, NIntersect(a, b))
}

func TestOneSidedNIntersect(t *testing.T) {
	t.Parallel()
	a := intSet()
	AssertEqualValues(t, a, NIntersect(a, intSet(42, 43), intSet(44, 45)))
}

func TestEqualNIntersect(t *testing.T) {
	t.Parallel()
	a := intSet(42, 43)
	AssertEqualValues(t, a, NIntersect(a, a, a))
}

func TestMixedNIntersect(t *testing.T) {
	t.Parallel()
	a := intSet(42, 43, 44)
	b := intSet(41, 42)
	c := intSet(42, 43)
	d := intSet(42)
	AssertEqualValues(t, intSet(42), NIntersect(a, b, c, d))
}

func TestTrivialUnion(t *testing.T) {
	t.Parallel()
	a := intSet()
	b := intSet()
	AssertEqualValues(t, a, Union(a, b))
}

func TestOneSidedUnion(t *testing.T) {
	t.Parallel()
	a := intSet()
	b := intSet(42, 43)
	AssertEqualValues(t, b, Union(a, b))
	AssertEqualValues(t, b, Union(b, a))
}

func TestEqualUnion(t *testing.T) {
	t.Parallel()
	a := intSet(42, 43)
	b := intSet(42, 43)
	AssertEqualValues(t, a, Union(a, b))
}

func TestMixedUnion(t *testing.T) {
	t.Parallel()
	a := intSet(42, 43)
	b := intSet(43, 44)
	c := intSet(42, 43, 44)
	AssertEqualValues(t, c, Union(a, b))
}

func TestTrivialDifference(t *testing.T) {
	t.Parallel()
	a := intSet()
	b := intSet()
	AssertEqualValues(t, a, Difference(a, b))
}

func TestOneSidedDifference(t *testing.T) {
	t.Parallel()
	a := intSet()
	b := intSet(42, 43)
	AssertEqualValues(t, a, Difference(a, b))
	AssertEqualValues(t, b, Difference(b, a))
}

func TestEqualDifference(t *testing.T) {
	t.Parallel()
	a := intSet(42, 43)
	b := intSet(42, 43)
	AssertEqualValues(t, intSet(), Difference(a, b))
}

func TestMixedDifference(t *testing.T) {
	t.Parallel()
	a := intSet(42, 43)
	b := intSet(43, 44)
	AssertEqualValues(t, intSet(42), Difference(a, b))
	AssertEqualValues(t, intSet(44), Difference(b, a))
}
