package rel

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// frozen runs a Where predicate on several goroutines once the set is large
// enough for its parallelism to engage, which needs 2^17 elements at the
// default concurrency setting. A predicate that fails must therefore report
// its error without racing — this used to be a plain captured variable, and
// the race detector flags it at this size.
//
// The size is the point of the test, so it is not reduced: a smaller
// relation exercises only the sequential path, which never had the bug.
func TestWhereReportsErrorsWithoutRacing(t *testing.T) {
	if testing.Short() {
		t.Skip("builds a relation large enough to engage frozen's parallelism")
	}
	t.Parallel()

	const rows = 1 << 17 // frozen's parallelism threshold
	boom := errors.New("boom")

	sb := NewSetBuilder()
	for i := 0; i < rows; i++ {
		sb.Add(NewTuple(NewAttr("a", NewNumber(float64(i)))))
	}
	s, err := sb.Finish()
	require.NoError(t, err)
	r, is := s.(Relation)
	require.True(t, is, "a set of same-shaped tuples is a Relation")

	// Every row fails, so many goroutines report an error at once.
	_, err = r.Where(func(Value) (bool, error) { return false, boom })
	assert.ErrorIs(t, err, boom, "a failing predicate must surface its error")

	// One row in the middle fails; the error must still escape.
	_, err = r.Where(func(v Value) (bool, error) {
		if v.(Tuple).MustGet("a").(Number).Float64() == rows/2 {
			return false, boom
		}
		return true, nil
	})
	assert.ErrorIs(t, err, boom, "an error from any row must surface")

	// No row fails: the filter still works at this size.
	got, err := r.Where(func(v Value) (bool, error) {
		return v.(Tuple).MustGet("a").(Number).Float64() < 100, nil
	})
	require.NoError(t, err)
	assert.Equal(t, 100, got.Count())

	// The same, through a GenericSet rather than a Relation.
	gsb := NewSetBuilder()
	for i := 0; i < rows; i++ {
		gsb.Add(NewNumber(float64(i)))
	}
	gs, err := gsb.Finish()
	require.NoError(t, err)
	_, err = gs.Where(func(Value) (bool, error) { return false, boom })
	assert.ErrorIs(t, err, boom, "GenericSet.Where must surface its error too")
}
