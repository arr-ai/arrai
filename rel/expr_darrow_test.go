package rel

import (
	"context"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/wbnf/parser"
)

func TestDArrowExprErrorOnInvalidType(t *testing.T) {
	t.Parallel()

	ident := NewIdentExpr(*parser.NewScanner("."), ".")
	expr := NewDArrowExpr(*parser.NewScanner(""), NewTuple(), ident)

	AssertExprErrorEquals(t, expr, "=> not applicable to tuple: ()")
}

func TestDArrowExprErrorOnFnEvalError(t *testing.T) {
	t.Parallel()

	ident := NewIdentExpr(*parser.NewScanner("."), ".")
	// This will fail to eval, as in the previous test.
	badFn := NewDArrowExpr(*parser.NewScanner("() => ."), NewTuple(), ident)
	_, err := badFn.Eval(arraictx.InitRunCtx(context.Background()), EmptyScope)

	// When this fails, it will propagate the err above, wrapped in expr's context.
	wrapper := NewDArrowExpr(*parser.NewScanner("{1} => () => ."), MustNewSet(NewNumber(1)), badFn)

	AssertExprErrorEquals(t, wrapper, err.Error())
}
