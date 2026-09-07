package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenericSetArrayEnumeratorMixedTypesDoesNotPanic guards against a
// regression where ArrayEnumerator assumed every element was an "@"-indexed
// tuple. That assumption can never hold for a genuine GenericSet: a
// homogeneous set of "@"-indexed tuples canonicalises to a specialised type
// (Array/Dict/String/Bytes) instead of staying a GenericSet, so any actual
// GenericSet necessarily holds a mix of element shapes.
func TestGenericSetArrayEnumeratorMixedTypesDoesNotPanic(t *testing.T) {
	t.Parallel()

	s, err := NewSet(NewNumber(1), NewGoString("a"))
	assert.NoError(t, err)
	assert.IsType(t, GenericSet{}, s)

	var got []Value
	for e := s.(GenericSet).ArrayEnumerator(); e.MoveNext(); {
		got = append(got, e.Current())
	}
	assert.Len(t, got, 2)
}
