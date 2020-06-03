package rel

import (
	"testing"

	"github.com/arr-ai/wbnf/parser"
)

func TestTupleMapExprAccessors(t *testing.T) {
	t.Parallel()

	lhs := NewTuple(NewAttr("a", NewNumber(1)))
	fn := NewFunction(*parser.NewScanner(""), nil, NewNumber(2))

	expr := NewTupleMapExpr(*parser.NewScanner(""), lhs, fn).(*TupleMapExpr)

	AssertExprsEvalToSameValue(t, lhs, expr.LHS())
	AssertExprsEvalToSameValue(t, fn, expr.Fn())
}
