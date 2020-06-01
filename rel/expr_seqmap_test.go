package rel

import (
	"testing"

	"github.com/arr-ai/wbnf/parser"
)

func TestSequenceMapExprIdent(t *testing.T) {
	t.Parallel()

	AssertExprsEvalToSameValue(t,
		NewSequenceMapExpr(
			*parser.NewScanner(""),
			NewArray(
				NewNumber(float64(1)),
				NewNumber(float64(2)),
				NewNumber(float64(3)),
			),
			NewIdentExpr(*parser.NewScanner("."), "."),
		),
		NewArray(
			NewNumber(float64(1)),
			NewNumber(float64(2)),
			NewNumber(float64(3)),
		),
	)
}
