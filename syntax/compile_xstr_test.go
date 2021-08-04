package syntax

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/rel"
)

func TestCleanEmptyValTransformsArrayOfThreeEmptySetsIntoEmptyArray(t *testing.T) {
	t.Parallel()

	a, b, c := rel.None, rel.None, rel.None
	values := rel.NewArray(a, b, c)
	array, ok := rel.AsArray(values)
	require.True(t, ok)

	actual := cleanEmptyVal(array)
	require.NotNil(t, actual)
	require.Equal(t, 0, len(actual))
}
