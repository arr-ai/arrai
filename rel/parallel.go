package rel

import (
	"runtime"
	"sync"

	"github.com/arr-ai/frozen"
)

// Element-wise parallelism for evaluator operations that run an arr.ai
// closure per element: `where` predicates and `=>` bodies. Closures are
// pure and every evaluator cache is concurrency-safe, so elements can be
// evaluated on any goroutine; the machinery here only decides when the work
// is big enough to pay for the fan-out and keeps error reporting
// deterministic.
//
// The threshold is frozen's: one knob — the FROZEN_CONCURRENCY environment
// variable, frozen.SetMinParallelChunk and frozen.DisableParallel — governs
// both frozen's internal fan-out and arr.ai's.

// parallelRanges splits n element-wise tasks into contiguous ranges, one
// per worker. It returns nil when the work should stay sequential:
// parallelism disabled, or n too small to give two workers a full chunk.
func parallelRanges(n int) [][2]int {
	if !frozen.ParallelEnabled() {
		return nil
	}
	w := n / frozen.MinParallelChunk()
	if maxW := runtime.GOMAXPROCS(0); w > maxW {
		w = maxW
	}
	if w < 2 {
		return nil
	}
	ranges := make([][2]int, w)
	for i := range ranges {
		ranges[i] = [2]int{i * n / w, (i + 1) * n / w}
	}
	return ranges
}

// runRanges runs f for each range on its own goroutine and waits for all of
// them. A panic in any worker is re-raised on the caller's goroutine.
//
// Determinism contract: a worker that fails records its error and stops its
// own range, but no worker is interrupted by another's failure. Every range
// below the lowest failing one therefore runs to completion, so selecting
// the lowest worker's error reports the same error a sequential scan would:
// the first failing element's. Errors are exceptional; burning the rest of
// the scan on that path buys reproducibility cheaply.
func runRanges(ranges [][2]int, f func(worker, lo, hi int)) {
	var wg sync.WaitGroup
	panics := make([]any, len(ranges))
	for i, r := range ranges {
		wg.Add(1)
		go func(i, lo, hi int) {
			defer wg.Done()
			defer func() {
				if p := recover(); p != nil {
					panics[i] = p
				}
			}()
			f(i, lo, hi)
		}(i, r[0], r[1])
	}
	wg.Wait()
	for _, p := range panics {
		if p != nil {
			panic(p)
		}
	}
}

// firstErr returns the lowest-indexed error: the deterministic counterpart
// of a sequential scan's first error.
func firstErr(errs []error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
