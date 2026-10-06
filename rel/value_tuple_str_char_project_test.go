package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringCharTupleProject(t *testing.T) {
	tuple := NewStringCharTuple(3, 'x')
	for _, names := range []Names{
		NewNames(),
		NewNames("@"),
		NewNames(StringCharAttr),
		NewNames("@", StringCharAttr),
	} {
		assert.True(t, tuple.asGenericTuple().Project(names).Equal(tuple.Project(names)), "%v", names)
	}
	assert.Nil(t, tuple.Project(NewNames("missing")))
	assert.Nil(t, tuple.Project(NewNames("@", "missing")))
}

func TestStringCharTupleProjectAllocatesLessThanGenericRoute(t *testing.T) {
	tuple := NewStringCharTuple(3, 'x')
	names := NewNames("@")
	require.NotNil(t, tuple.Project(names))

	direct := testing.AllocsPerRun(100, func() { _ = tuple.Project(names) })
	viaGeneric := testing.AllocsPerRun(100, func() { _ = tuple.asGenericTuple().Project(names) })
	assert.Less(t, direct, viaGeneric)
}
