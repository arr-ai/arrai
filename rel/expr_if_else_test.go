package rel

import (
	"testing"

	"github.com/arr-ai/wbnf/parser"
)

func TestNewIfElseExpr(t *testing.T) {
	t.Parallel()
	AssertExprsEvalToSameValue(t, NewNumber(42),
		NewIfElseExpr(*parser.NewScanner("42"), NewNumber(42), True, NewNumber(43)),
	)
	AssertExprsEvalToSameValue(t, NewNumber(43),
		NewIfElseExpr(*parser.NewScanner("42"), NewNumber(42), False, NewNumber(43)),
	)
}
