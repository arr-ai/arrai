package rel

import (
	"testing"

	"github.com/arr-ai/wbnf/parser"
)

func TestDArrowExprAccessors(t *testing.T) {
	t.Parallel()

	lhs := NewNumber(1)
	fn := NewFunction(*parser.NewScanner(""), nil, NewNumber(2))

	expr := NewDArrowExpr(*parser.NewScanner(""), lhs, fn).(*DArrowExpr)

	AssertExprsEvalToSameValue(t, lhs, expr.LHS())
	AssertExprsEvalToSameValue(t, fn, expr.Fn())
}
