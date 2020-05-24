package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDictEntryTupleLess(t *testing.T) {
	// Test for a panic when comparing Strings using !=.
	assert.NotPanics(t, func() {
		NewDict(false,
			NewDictEntryTuple(NewString([]rune("b")), NewNumber(2)),
			NewDictEntryTuple(NewString([]rune("a")), NewNumber(1)),
		).(Dict).OrderedEntries()
	})
}

func TestDictEntryTupleOrdered(t *testing.T) {
	entries := NewDict(true,
		NewDictEntryTuple(NewString([]rune("b")), NewNumber(2)),
		NewDictEntryTuple(NewString([]rune("a")), NewNumber(2)),
		NewDictEntryTuple(NewString([]rune("b")), NewNumber(1)),
		NewDictEntryTuple(NewString([]rune("a")), NewNumber(1)),
	).(Dict).OrderedEntries()

	AssertEqualValues(t, NewDictEntryTuple(NewString([]rune("a")), NewNumber(1)), entries[0])
	AssertEqualValues(t, NewDictEntryTuple(NewString([]rune("a")), NewNumber(2)), entries[1])
	AssertEqualValues(t, NewDictEntryTuple(NewString([]rune("b")), NewNumber(1)), entries[2])
	AssertEqualValues(t, NewDictEntryTuple(NewString([]rune("b")), NewNumber(2)), entries[3])
}

func TestDictLess(t *testing.T) {
	kv := func(k, v float64) DictEntryTuple {
		return NewDictEntryTuple(NewNumber(k), NewNumber(v))
	}
	assertLess := func(a, b Set) {
		assert.True(t, a.Less(b))
		assert.False(t, b.Less(a))
	}
	assertLess(NewDict(true, kv(1, 42)), NewDict(true, kv(1, 43)))
	assertLess(NewDict(true, kv(1, 42)), NewDict(true, kv(1, 43), kv(2, 44)))
	assertLess(NewDict(true, kv(1, 42)), NewDict(true, kv(1, 42), kv(1, 44)))
	assertLess(NewDict(true, kv(1, 41), kv(1, 42)), NewDict(true, kv(1, 42)))
	assertLess(NewDict(true, kv(1, 42)), NewDict(true, kv(1, 43), kv(2, 42)))
	assertLess(NewDict(true, kv(1, 42), kv(2, 43)), NewDict(true, kv(1, 42), kv(3, 43)))

	assertSame := func(a, b Set) {
		assert.False(t, a.Less(b))
		assert.False(t, b.Less(a))
	}
	assertSame(NewDict(true, kv(1, 43), kv(2, 42)), NewDict(true, kv(1, 43), kv(2, 42)))
}

func TestDictCallAll(t *testing.T) {
	t.Parallel()

	kv := func(k, v float64) DictEntryTuple {
		return NewDictEntryTuple(NewNumber(k), NewNumber(v))
	}
	dict := NewDict(false, kv(1, 10), kv(2, 20), kv(3, 30))

	AssertEqualValues(t, NewSet(NewNumber(10)), dict.CallAll(NewNumber(1)))
	AssertEqualValues(t, NewSet(NewNumber(20)), dict.CallAll(NewNumber(2)))
	AssertEqualValues(t, NewSet(NewNumber(30)), dict.CallAll(NewNumber(3)))
	AssertEqualValues(t, None, dict.CallAll(NewNumber(4)))

	dict = NewDict(true, kv(1, 10), kv(1, 11), kv(2, 20), kv(3, 30))

	AssertEqualValues(t, NewSet(NewNumber(10), NewNumber(11)), dict.CallAll(NewNumber(1)))
	AssertEqualValues(t, NewSet(NewNumber(20)), dict.CallAll(NewNumber(2)))
	AssertEqualValues(t, NewSet(NewNumber(30)), dict.CallAll(NewNumber(3)))
	AssertEqualValues(t, None, dict.CallAll(NewNumber(4)))
}
