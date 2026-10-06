package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecialisedTupleProjectRequiresEveryName(t *testing.T) {
	tuples := map[string]struct {
		tuple Tuple
		attr  string
	}{
		"char":  {NewStringCharTuple(1, 'x'), StringCharAttr},
		"array": {NewArrayItemTuple(1, NewNumber(2)), ArrayItemAttr},
		"bytes": {NewBytesByteTuple(1, 2), BytesByteAttr},
		"dict":  {NewDictEntryTuple(NewNumber(1), NewNumber(2)), DictValueAttr},
	}
	for name, tc := range tuples {
		tc := tc
		t.Run(name, func(t *testing.T) {
			names := NewNames("@", tc.attr)
			assert.True(t, tc.tuple.Project(names).Equal(tc.tuple))
			assert.Nil(t, tc.tuple.Project(NewNames("@", tc.attr, "extra")))
			assert.Nil(t, tc.tuple.Project(NewNames("extra")))
			assert.True(t, tc.tuple.Project(NewNames("@")).Equal(NewTuple(NewAttr("@", tc.tuple.MustGet("@")))))
		})
	}
}
