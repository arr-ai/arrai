package rel

import (
	"testing"

	"github.com/arr-ai/wbnf/parser"
)

func TestSequenceMapExprAccessors(t *testing.T) {
	t.Parallel()

	lhs := NewArray(NewNumber(float64(1)))
	fn := NewNumber(2)

	expr := NewSequenceMapExpr(*parser.NewScanner(""), lhs, fn).(*SequenceMapExpr)

	AssertExprsEvalToSameValue(t, lhs, expr.LHS())
	AssertExprsEvalToSameValue(t, fn, expr.Fn())
}
