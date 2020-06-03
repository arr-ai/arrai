package rel

import (
	"testing"

	"github.com/arr-ai/wbnf/parser"
)

func TestArrowExprAccessors(t *testing.T) {
	t.Parallel()

	lhs := NewNumber(1)
	fn := NewFunction(*parser.NewScanner(""), nil, NewNumber(2))

	expr := NewArrowExpr(*parser.NewScanner(""), lhs, fn).(*ArrowExpr)

	AssertExprsEvalToSameValue(t, lhs, expr.LHS())
	AssertExprsEvalToSameValue(t, fn, expr.Fn())
}
