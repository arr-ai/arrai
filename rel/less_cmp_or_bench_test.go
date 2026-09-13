package rel

import (
	"cmp"
	"testing"
)

// cmp.Or evaluation for 🎯T14 / #741. The two-field lexicographic Less used
// by StringCharTuple (and FilePos) is the best-case for a comparator chain:
// two ints, no interface dispatch. Measured on Apple M4 Max / go1.25:
//
//	hand-rolled  4.97 ns/op  0 B  0 allocs
//	cmp.Or       8.90 ns/op  0 B  0 allocs   (~1.8× slower)
//
// Rejected. Kind-prefix is identical on both sides and omitted so the
// comparison is the lexicographic core only.

func lessHandrolled(t, u StringCharTuple) bool {
	if t.at != u.at {
		return t.at < u.at
	}
	return t.char < u.char
}

func lessCmpOr(t, u StringCharTuple) bool {
	return cmp.Or(
		cmp.Compare(t.at, u.at),
		cmp.Compare(t.char, u.char),
	) < 0
}

var lessSink bool //nolint:gochecknoglobals

func benchTwoFieldLess(b *testing.B, fn func(StringCharTuple, StringCharTuple) bool) {
	b.Helper()
	// Mix of first-field differs (common) and first-field ties (falls through).
	pairs := [][2]StringCharTuple{
		{NewStringCharTuple(1, 'a'), NewStringCharTuple(2, 'a')},
		{NewStringCharTuple(3, 'x'), NewStringCharTuple(3, 'y')},
		{NewStringCharTuple(0, 'z'), NewStringCharTuple(0, 'z')},
		{NewStringCharTuple(9, 'm'), NewStringCharTuple(8, 'm')},
	}
	b.ReportAllocs()
	for b.Loop() {
		for _, p := range pairs {
			lessSink = fn(p[0], p[1])
		}
	}
}

func BenchmarkLessHandrolled(b *testing.B) { benchTwoFieldLess(b, lessHandrolled) }
func BenchmarkLessCmpOr(b *testing.B)      { benchTwoFieldLess(b, lessCmpOr) }
