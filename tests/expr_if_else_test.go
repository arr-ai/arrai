package tests

import (
	"testing"

	"github.com/arr-ai/arrai/rel"
)

// TestNewIfElseExpr tests rel.NewSet.
func TestNewIfElseExpr(t *testing.T) {
	AssertExprsEvalToSameValue(t, rel.NewNumber(42),
		rel.NewIfElseExpr(rel.NewNumber(42), rel.True, rel.NewNumber(43)),
	)
	AssertExprsEvalToSameValue(t, rel.NewNumber(43),
		rel.NewIfElseExpr(rel.NewNumber(42), rel.False, rel.NewNumber(43)),
	)
}
