package rel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func mustEncodePlan(t *testing.T, v Value) []byte {
	t.Helper()
	p, err := LowerPlan(v)
	require.NoError(t, err)
	b, err := EncodePlan(p)
	require.NoError(t, err)
	return b
}

// 🎯T26: equivalent sets must serialize identically regardless of insertion
// order. Before 815e720, encodeSetEnum walked frozen's Enumerator, so
// NewSet(3,1,2,9,4) and NewSet(1,2,3,4,9) produced different plan.bin.
func TestPlanSetEncodingIsCanonical(t *testing.T) {
	t.Parallel()
	n := func(x float64) Number { return NewNumber(x) }
	mustSet := func(vs ...Value) Set {
		s, err := NewSet(vs...)
		require.NoError(t, err)
		return s
	}
	a := mustSet(n(3), n(1), n(2), n(9), n(4))
	b := mustSet(n(1), n(2), n(3), n(4), n(9))
	require.Equal(t, mustEncodePlan(t, a), mustEncodePlan(t, b))
}

func TestPlanDictEncodingIsCanonical(t *testing.T) {
	t.Parallel()
	n := func(x float64) Number { return NewNumber(x) }
	e := func(k, v float64) DictEntryTuple { return NewDictEntryTuple(n(k), n(v)) }
	a := MustNewDict(false, e(3, 1), e(1, 2), e(9, 4))
	b := MustNewDict(false, e(1, 2), e(9, 4), e(3, 1))
	require.Equal(t, mustEncodePlan(t, a), mustEncodePlan(t, b))
}
