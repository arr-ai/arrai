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
