package rel

import (
	"fmt"
	"testing"

	"github.com/arr-ai/frozen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// bigPosRel builds a relation big enough to cross the parallel threshold.
func bigPosRel(tb testing.TB, n int) *positionalRelation {
	tb.Helper()
	b := newStoreBuilder(2, n, false)
	for i := 0; i < n; i++ {
		b.add(row(i, i%97))
	}
	return b.finish()
}

func TestWhereParallelMatchesSequential(t *testing.T) {
	t.Parallel()

	n := 4 * frozen.MinParallelChunk()
	r := bigPosRel(t, n)
	require.NotNil(t, parallelRanges(n), "n must be over the parallel threshold")

	pred := func(v Values) (bool, error) {
		return int(v[1].(Number).Float64()) < 13, nil
	}
	got, err := r.Where(pred)
	require.NoError(t, err)

	// Sequential reference, bypassing the parallel gate.
	sel := make([]uint32, 0, n)
	for i := 0; i < r.n; i++ {
		match, err := pred(r.rowAt(i))
		require.NoError(t, err)
		if match {
			sel = append(sel, r.arenaID(i))
		}
	}
	want := r.selView(sel)
	assert.Equal(t, want.Count(), got.Count())
	assert.True(t, got.EqualPositionalRelation(want))

	// All-pass returns the same view, as in the sequential path.
	all, err := r.Where(func(Values) (bool, error) { return true, nil })
	require.NoError(t, err)
	assert.Same(t, r, all)
}

// A parallel where must report the error of the first failing row in view
// order, exactly as a sequential scan would, no matter which worker fails
// first on the clock.
func TestWhereParallelFirstError(t *testing.T) {
	t.Parallel()

	n := 4 * frozen.MinParallelChunk()
	r := bigPosRel(t, n)
	first := int(r.rowAt(300)[0].(Number).Float64())
	for i := 0; i < 20; i++ {
		_, err := r.Where(func(v Values) (bool, error) {
			if a := int(v[0].(Number).Float64()); a >= 300 {
				return false, fmt.Errorf("boom at %d", a)
			}
			return true, nil
		})
		require.EqualError(t, err, fmt.Sprintf("boom at %d", first))
	}
}

func TestParallelRanges(t *testing.T) {
	// Not parallel: mutates the process-wide parallelism setting.
	restore := frozen.MinParallelChunk()
	defer frozen.SetMinParallelChunk(restore)

	assert.Nil(t, parallelRanges(frozen.MinParallelChunk()*2-1))
	ranges := parallelRanges(frozen.MinParallelChunk() * 2)
	require.NotNil(t, ranges)
	// Ranges tile [0, n) contiguously in order.
	at := 0
	for _, r := range ranges {
		assert.Equal(t, at, r[0])
		at = r[1]
	}
	assert.Equal(t, frozen.MinParallelChunk()*2, at)

	frozen.DisableParallel()
	assert.Nil(t, parallelRanges(1<<20), "disabled parallelism must force sequential")
	frozen.SetMinParallelChunk(restore)
	assert.NotNil(t, parallelRanges(1<<20))
}

// A predicate panic on a worker goroutine must surface on the caller, not
// crash the process.
func TestWhereParallelPanic(t *testing.T) {
	t.Parallel()

	r := bigPosRel(t, 4*frozen.MinParallelChunk())
	assert.PanicsWithValue(t, "kaboom", func() {
		_, _ = r.Where(func(v Values) (bool, error) { //nolint:errcheck // the panic preempts the return
			if int(v[0].(Number).Float64()) == 500 {
				panic("kaboom")
			}
			return true, nil
		})
	})
}
