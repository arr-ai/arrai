package syntax

import (
	"testing"
)

// These run `where` and `=>` over sets big enough to cross the parallel
// evaluation threshold (frozen.MinParallelChunk × 2, default 512), so the
// closure bodies evaluate on multiple goroutines. An array is a relation of
// (@, @item) tuples, so //seq.repeat provides a large arena-backed relation
// with indices. Under the slowpath tag the same tests run fully
// sequentially, giving the differential oracle a parallel-vs-sequential
// axis for free.

func TestParallelWhereRelation(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`100`,
		`(//seq.repeat(5000, [0]) => (a: .@, b: .@ % 7) where .a < 100) count`,
	)
	AssertCodesEvalToSameValue(t,
		`715`,
		`(//seq.repeat(5000, [0]) => (a: .@, b: .@ % 7) where .b = 0) count`,
	)
	// A predicate error must surface from a worker, and deterministically:
	// the first failing element in enumeration order, here @ = 2.
	AssertCodeErrors(t, "Call: no return values for input 2",
		`//seq.repeat(5000, [0]) => (a: .@) where {0: 1, 1: 1}(.a) < 2`,
	)
}

func TestParallelDArrow(t *testing.T) {
	t.Parallel()
	AssertCodesEvalToSameValue(t,
		`5000`,
		`(//seq.repeat(5000, [0]) => .@ * 2) count`,
	)
	// Nested transform: outer bodies run in parallel, inner sets are small.
	AssertCodesEvalToSameValue(t,
		`5000`,
		`(//seq.repeat(5000, [0]) => (x: .@, y: {1, 2, 3} => . * 2)) count`,
	)
	AssertCodeErrors(t, "Call: no return values for input 2",
		`//seq.repeat(5000, [0]) => {0: 1, 1: 1}(.@)`,
	)
}
