package rel

import (
	"context"
	"testing"

	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"

	"github.com/arr-ai/arrai/pkg/arraictx"
)

func TestDotExprAccessors(t *testing.T) {
	t.Parallel()

	lhs := NewTuple(NewAttr("a", NewNumber(1)))
	attr := "a"
	expr := NewDotExpr(*parser.NewScanner(""), lhs, attr).(*DotExpr)

	AssertExprsEvalToSameValue(t, expr.Subject(), lhs)
	assert.Equal(t, expr.Attr(), attr)
}

func TestDotExprErrorOnMissingEmptyAttr(t *testing.T) {
	t.Parallel()

	lhs := NewTuple(NewAttr("a", NewNumber(1)))
	expr := NewDotExpr(*parser.NewScanner(""), lhs, "").(*DotExpr)

	AssertExprErrorEquals(t, expr, `Missing attr "" (available: |a|)`)
}

func TestDotExprErrorOnInvalidStarAttr(t *testing.T) {
	t.Parallel()

	expr := NewDotExpr(*parser.NewScanner(""), NewTuple(), "*")

	AssertExprErrorEquals(t, expr, "expr.* not allowed outside tuple attr")
}

func TestDotExprErrorOnEvalError(t *testing.T) {
	t.Parallel()

	// This will fail to eval, as in the previous test.
	lhs := NewDotExpr(*parser.NewScanner("().*"), NewTuple(), "*")
	_, err := lhs.Eval(arraictx.InitRunCtx(context.Background()), EmptyScope)

	// When this fails, it will propagate the err above, wrapped in expr's context.
	expr := NewDotExpr(*parser.NewScanner("().*.a"), lhs, "a")

	AssertExprErrorEquals(t, expr, err.Error())
}

func TestDotExprErrorOnNonEnumerableSet(t *testing.T) {
	t.Parallel()

	expr := NewDotExpr(*parser.NewScanner("native.a"), NewNativeFunction("native", nil), "a")

	AssertExprErrorEquals(t, expr, `Cannot get attr "a" from native-function`)

	expr = NewDotExpr(*parser.NewScanner("closure.a"), NewClosure(Scope{}, nil), "a")

	AssertExprErrorEquals(t, expr, `Cannot get attr "a" from closure`)
}
