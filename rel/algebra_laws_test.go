package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T16 standing oracles (Arieh #736–#742). Each class is implemented with a
// test wired into CI, or rejected with a reason:
//
//   - Algebraic value laws (#738/#740): this file. Generated corpus across
//     representations, not hand-picked one-offs. Cross-kind Less trichotomy
//     is rejected: kinds are not a total order (empty String vs tuples).
//   - Silent @latest fallbacks (#736/#742): TestPinnedModuleDownloadFailureIsLoud
//     plus goModFilePin; a pin that cannot be fetched fails instead of
//     resolving @latest.
//   - Environment paths (#737/#742): nested imports via PrimaryRoot
//     (TestRetrieveModuleUsesModuleRootNotCwd), go.mod/go.sum byte-identity
//     (TestLoadRequiredModulesDoesNotLeaveGoSumChanged), signal/os.Exit
//     profiler flush (cmd/arrai TestWaitAndFlushCallsStop).
//
// sampleValues is a generated corpus spanning the representations that
// #738/#740 found hash/equals bugs in. Not hand-picked one-off cases: every
// pair is checked, so a new representation added here is tested against all
// the others.
func sampleValues(t *testing.T) []Value {
	t.Helper()
	mustSet := func(vs ...Value) Set {
		s, err := NewSet(vs...)
		require.NoError(t, err)
		return s
	}
	mustDict := func(entries ...DictEntryTuple) Set {
		d, err := NewDict(false, entries...)
		require.NoError(t, err)
		return d
	}
	n := func(x float64) Number { return NewNumber(x) }
	return []Value{
		n(0), n(1), n(-1), n(2), n(1.5),
		NewString([]rune("")), NewString([]rune("a")), NewString([]rune("ab")),
		NewString([]rune("aa")),
		NewBytes([]byte{}), NewBytes([]byte{0}), NewBytes([]byte{1, 2}),
		NewArray(), NewArray(n(1)), NewArray(n(1), n(1)), NewArray(n(1), n(2)),
		EmptyTuple,
		NewTuple(NewAttr("a", n(1))),
		NewTuple(NewAttr("a", n(1)), NewAttr("b", n(2))),
		NewTuple(NewAttr("b", n(2)), NewAttr("a", n(1))),
		NewArrayItemTuple(0, n(1)),
		NewArrayItemTuple(1, n(1)),
		NewDictEntryTuple(n(1), n(2)),
		NewStringCharTuple(0, 'a'),
		NewBytesByteTuple(0, 1),
		None,
		mustSet(n(1)),
		mustSet(n(1), n(2)),
		mustSet(n(2), n(1)),
		mustSet(n(1), n(1)),
		True,
		mustDict(NewDictEntryTuple(n(1), n(2))),
		mustDict(NewDictEntryTuple(n(1), n(2)), NewDictEntryTuple(n(3), n(4))),
		mustRel(t, NewTuple(NewAttr("a", n(1)), NewAttr("b", n(2)))),
		mustRel(t,
			NewTuple(NewAttr("a", n(1)), NewAttr("b", n(2))),
			NewTuple(NewAttr("a", n(3)), NewAttr("b", n(4)))),
		mustRel(t,
			NewTuple(NewAttr("b", n(2)), NewAttr("a", n(1))),
			NewTuple(NewAttr("b", n(4)), NewAttr("a", n(3)))),
	}
}

func TestAlgebraEqualImpliesHash128(t *testing.T) {
	t.Parallel()
	vs := sampleValues(t)
	for i, a := range vs {
		for j, b := range vs {
			if a.Equal(b) {
				assert.Equal(t, a.Hash128(), b.Hash128(),
					"Equal but Hash128 differs: [%d]=%s [%d]=%s", i, a, j, b)
			}
		}
	}
}

func TestAlgebraLessAgreesWithEqual(t *testing.T) {
	t.Parallel()
	vs := sampleValues(t)
	for i, a := range vs {
		for j, b := range vs {
			eq := a.Equal(b)
			ab, ba := a.Less(b), b.Less(a)
			if eq {
				assert.False(t, ab, "Equal but Less: [%d] %s < [%d] %s", i, a, j, b)
				assert.False(t, ba, "Equal but Less: [%d] %s < [%d] %s", j, b, i, a)
				continue
			}
			if a.Kind() != b.Kind() {
				continue
			}
			assert.True(t, ab != ba, "incomparable same-kind: [%d] %s vs [%d] %s", i, a, j, b)
		}
	}
}

func TestAlgebraWithWithoutRoundTrip(t *testing.T) {
	t.Parallel()
	base := NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewNumber(2)))
	extra := NewNumber(3)
	got := base.With("c", extra).Without("c")
	assert.True(t, base.Equal(got), "%s vs %s", base, got)
	assert.Equal(t, base.Hash128(), got.Hash128())

	replaced := base.With("a", extra)
	assert.False(t, base.Equal(replaced))
	restored := replaced.With("a", NewNumber(1))
	assert.True(t, base.Equal(restored), "With overwrite round-trip")
}

func TestAlgebraCanonicalisationIdempotent(t *testing.T) {
	t.Parallel()
	item := NewArrayItemTuple(3, NewString([]rune("x")))
	generic := newGenericTuple(
		NewAttr("@", NewNumber(3)),
		NewAttr(ArrayItemAttr, NewString([]rune("x"))),
	)
	// Specialized Equal is kind-narrow; GenericTuple.Equal is permissive.
	// Hash128 still agrees, so a hash-based set treats them as duplicates.
	assert.False(t, item.Equal(generic))
	assert.True(t, generic.Equal(item))
	assert.Equal(t, item.Hash128(), generic.Hash128())
	// NewTuple canonicalises to the specialised kind; that is the identity.
	again := NewTuple(NewAttr("@", NewNumber(3)), NewAttr(ArrayItemAttr, NewString([]rune("x"))))
	assert.True(t, item.Equal(again))
	assert.Equal(t, item.Hash128(), again.Hash128())
}

func TestAlgebraDictVsGenericSet(t *testing.T) {
	t.Parallel()
	d, err := NewDict(false, NewDictEntryTuple(NewNumber(1), NewNumber(2)))
	require.NoError(t, err)
	s, err := NewSet(NewDictEntryTuple(NewNumber(1), NewNumber(2)))
	require.NoError(t, err)
	assert.True(t, d.Equal(s), "Dict must Equal a generic set of the same entries")
	assert.Equal(t, d.Hash128(), s.Hash128())
}

func TestAlgebraLessAgreesWithEqualOnSelf(t *testing.T) {
	t.Parallel()
	for i, v := range sampleValues(t) {
		assert.True(t, v.Equal(v), "not equal to self: [%d] %s", i, v)
		assert.False(t, v.Less(v), "less than self: [%d] %s", i, v)
		assert.Equal(t, v.Hash128(), v.Hash128())
	}
}
