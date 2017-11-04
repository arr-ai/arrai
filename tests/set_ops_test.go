package tests

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
)

// TestTrivialIntersect test rel.Intersect({}, {}).
func TestTrivialIntersect(t *testing.T) {
	a := intSet()
	b := intSet()
	assertEqualValues(t, a, rel.Intersect(a, b))
}

// TestOneSidedIntersect tests rel.Intersect({}, {a:42, b:43}) and vice-versa.
func TestOneSidedIntersect(t *testing.T) {
	a := intSet()
	b := intSet(42, 43)
	assertEqualValues(t, a, rel.Intersect(a, b))
	assertEqualValues(t, a, rel.Intersect(b, a))
}

// TestEqualIntersect tests rel.Intersect({a:42, b:43}, {a:42, b:43}).
func TestEqualIntersect(t *testing.T) {
	a := intSet(42, 43)
	b := intSet(42, 43)
	assertEqualValues(t, a, rel.Intersect(a, b))
}

// TestMixedIntersect tests rel.Intersect({a:42, b:43}, {b:43, c:44}).
func TestMixedIntersect(t *testing.T) {
	a := intSet(42, 43)
	b := intSet(43, rel.NewNumber(44))
	c := intSet(43)
	assertEqualValues(t, c, rel.Intersect(a, b))
}

// TestTrivialUnion test rel.Union({}, {}).
func TestTrivialUnion(t *testing.T) {
	a := intSet()
	b := intSet()
	assertEqualValues(t, a, rel.Union(a, b))
}

// TestOneSidedUnion tests rel.Union({}, {a:42, b:43}) and vice-versa.
func TestOneSidedUnion(t *testing.T) {
	a := intSet()
	b := intSet(42, 43)
	assertEqualValues(t, b, rel.Union(a, b))
	assertEqualValues(t, b, rel.Union(b, a))
}

// TestEqualUnion tests rel.Union({a:42, b:43}, {a:42, b:43}).
func TestEqualUnion(t *testing.T) {
	a := intSet(42, 43)
	b := intSet(42, 43)
	assertEqualValues(t, a, rel.Union(a, b))
}

// TestMixedUnion tests rel.Union({a:42, b:43}, {b:43, c:44}).
func TestMixedUnion(t *testing.T) {
	a := intSet(42, 43)
	b := intSet(43, 44)
	c := intSet(42, 43, 44)
	assertEqualValues(t, c, rel.Union(a, b))
}

// TestTrivialDifference test rel.Difference({}, {}).
func TestTrivialDifference(t *testing.T) {
	a := intSet()
	b := intSet()
	assertEqualValues(t, a, rel.Difference(a, b))
}

// TestOneSidedDifference tests rel.Difference({}, {a:42, b:43}) and vice-versa.
func TestOneSidedDifference(t *testing.T) {
	a := intSet()
	b := intSet(42, 43)
	assertEqualValues(t, a, rel.Difference(a, b))
	assertEqualValues(t, b, rel.Difference(b, a))
}

// TestEqualDifference tests rel.Difference({a:42, b:43}, {a:42, b:43}).
func TestEqualDifference(t *testing.T) {
	a := intSet(42, 43)
	b := intSet(42, 43)
	assertEqualValues(t, intSet(), rel.Difference(a, b))
}

// TestMixedDifference tests rel.Difference({a:42, b:43}, {b:43, c:44}).
func TestMixedDifference(t *testing.T) {
	a := intSet(42, 43)
	b := intSet(43, 44)
	assertEqualValues(t, intSet(42), rel.Difference(a, b))
	assertEqualValues(t, intSet(44), rel.Difference(b, a))
}
