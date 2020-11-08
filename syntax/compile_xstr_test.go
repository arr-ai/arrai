package syntax

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
	"github.com/stretchr/testify/require"
)

func TestCleanEmptyValTransformsArrayOfThreeEmptySetsIntoEmptyArray(t *testing.T) {
	t.Parallel()

	a, err := rel.NewSet()
	require.NoError(t, err)
	b, err := rel.NewSet()
	require.NoError(t, err)
	c, err := rel.NewSet()
	require.NoError(t, err)
	values := rel.NewArray(a, b, c)
	array, ok := rel.AsArray(values)
	require.True(t, ok)

	actual := cleanEmptyVal(array)
	require.NotNil(t, actual)
	require.Equal(t, 0, len(actual))
}
