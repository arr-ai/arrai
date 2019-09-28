package rel

import (
	"testing"
)

func TestNewIfElseExpr(t *testing.T) {
	t.Parallel()
	AssertExprsEvalToSameValue(t, NewNumber(42),
		NewIfElseExpr(NewNumber(42), True, NewNumber(43)),
	)
	AssertExprsEvalToSameValue(t, NewNumber(43),
		NewIfElseExpr(NewNumber(42), False, NewNumber(43)),
	)
}
