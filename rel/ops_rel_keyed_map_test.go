package rel

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// appRows returns perApp rows for each of apps app names, keyed by a fresh
// (appName: ['AppN']) array so no two rows share a key pointer.
func appRows(t *testing.T, apps, perApp int) Set {
	t.Helper()
	b := NewSetBuilder()
	for i := 0; i < apps; i++ {
		for j := 0; j < perApp; j++ {
			name := NewArray(NewString([]rune("App" + strconv.Itoa(i))))
			b.Add(NewTuple(NewAttr("appName", name), NewAttr("v", NewNumber(float64(j)))))
		}
	}
	rows, err := b.Finish()
	require.NoError(t, err)
	require.Equal(t, apps*perApp, rows.Count())
	return rows
}

func appNameKey(v Value) Value { return v.(Tuple).Project(NewNames("appName")) }

// Reduce and GenericJoin key a frozen map by Value, so the map must hash
// keys by content. frozen v1 expects Hash128() or Hash(seed), neither of
// which rel implements since the move to seedless Hash(), and its reflective
// fallback is not content-stable: at this size it returned 1,016 groups for
// 1,000 app names (🎯T34). The failure needs thousands of rows to show, so
// these tests are deliberately not small.
func TestReduceGroupsManyEqualKeys(t *testing.T) {
	t.Parallel()

	const apps, perApp = 1000, 3
	got := Reduce(appRows(t, apps, perApp), appNameKey,
		func(key Value, tuples Set) Set {
			return MustNewSet(Merge(key.(Tuple), NewTuple(NewAttr("n", NewNumber(float64(tuples.Count()))))))
		},
	)
	require.Equal(t, apps, got.Count(), "one group per app name")
	for e := got.Enumerator(); e.MoveNext(); {
		require.Equal(t, NewNumber(perApp), e.Current().(Tuple).MustGet("n"), "%v", e.Current())
	}
}

func TestGenericJoinMatchesManyEqualKeys(t *testing.T) {
	t.Parallel()

	const apps, perApp = 1000, 3
	got := GenericJoin(appRows(t, apps, perApp), appRows(t, apps, 1), appNameKey,
		func(key Value, as, bs Set) Set {
			return MustNewSet(Merge(key.(Tuple), NewTuple(
				NewAttr("l", NewNumber(float64(as.Count()))),
				NewAttr("r", NewNumber(float64(bs.Count()))),
			)))
		},
	)
	require.Equal(t, apps, got.Count(), "one joined group per app name")
	for e := got.Enumerator(); e.MoveNext(); {
		row := e.Current().(Tuple)
		require.Equal(t, NewNumber(perApp), row.MustGet("l"), "%v", row)
		require.Equal(t, NewNumber(1), row.MustGet("r"), "%v", row)
	}
}
