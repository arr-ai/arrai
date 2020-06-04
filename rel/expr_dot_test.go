package rel

import (
	"testing"

	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
)

func TestDotExprAccessors(t *testing.T) {
	t.Parallel()

	lhs := NewTuple(NewAttr("a", NewNumber(1)))
	attr := "a"
	expr := NewDotExpr(*parser.NewScanner(""), lhs, attr).(*DotExpr)

	AssertExprsEvalToSameValue(t, expr.Subject(), lhs)
	assert.Equal(t, expr.Attr(), attr)
}

func TestDotExprErrorOnInvalidStarAttr(t *testing.T) {
	t.Parallel()

	expr := NewDotExpr(*parser.NewScanner(""), NewTuple(), "*")

	AssertExprErrorEquals(t, expr, "expr.* not allowed outside tuple attr")
}

func TestDotExprErrorOnEvalError(t *testing.T) {
	t.Parallel()

	// This will fail to eval, as in the previous test.
	lhs := NewDotExpr(*parser.NewScanner(""), NewTuple(), "*")
	_, err := lhs.Eval(EmptyScope)

	// When this fails, it will propagate the err above, wrapped in expr's context.
	expr := NewDotExpr(*parser.NewScanner(""), lhs, "a")

	AssertExprErrorEquals(t, expr, err.Error())
}
