package syntax

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanEmptyPiecesDropsAllEmptyComputedPieces(t *testing.T) {
	t.Parallel()

	actual := cleanEmptyPieces([]xstrPiece{{}, {}, {}})
	require.NotNil(t, actual)
	require.Equal(t, 0, len(actual))
}
