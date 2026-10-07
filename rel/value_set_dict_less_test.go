package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestDict(t *testing.T, kv ...int) Dict {
	t.Helper()
	entries := make([]DictEntryTuple, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		entries = append(entries, NewDictEntryTuple(NewNumber(float64(kv[i])), NewNumber(float64(kv[i+1]))))
	}
	set, err := NewDict(false, entries...)
	require.NoError(t, err)
	d, ok := set.(Dict)
	require.True(t, ok)
	return d
}

func TestDictLessOrdersByKeysThenValues(t *testing.T) {
	a := newTestDict(t, 1, 10, 2, 20)
	b := newTestDict(t, 1, 10, 3, 20)
	c := newTestDict(t, 1, 11, 2, 20)
	prefix := newTestDict(t, 1, 10)

	for i := 0; i < 2; i++ { // the second round runs against memoised keys
		assert.True(t, a.Less(b))
		assert.False(t, b.Less(a))
		assert.True(t, a.Less(c))
		assert.False(t, c.Less(a))
		assert.True(t, prefix.Less(a))
		assert.False(t, a.Less(a))
	}
}

func TestDictLessDoesNotReorderKeysPerComparison(t *testing.T) {
	a := newTestDict(t, 1, 10, 2, 20, 3, 30, 4, 40, 5, 50, 6, 60)
	b := newTestDict(t, 1, 10, 2, 20, 3, 30, 4, 40, 5, 50, 7, 70)
	require.True(t, a.Less(b))

	less := testing.AllocsPerRun(100, func() { _ = a.Less(b) })
	reorder := testing.AllocsPerRun(100, func() { _ = a.m.Keys().OrderedElements(ValueLess) })
	assert.Less(t, less, reorder)
}
