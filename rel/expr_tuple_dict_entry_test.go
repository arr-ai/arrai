package rel

import (
	"context"
	"strings"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
)

func TestDictEntryTupleExprString(t *testing.T) {
	t.Parallel()

	e := NewDictEntryTupleExpr(*parser.NewScanner(""), NewNumber(1), NewNumber(2))

	assert.Equal(t, "(@: 1, @value: 2)", e.String())
}

func TestDictEntryTupleExprEval(t *testing.T) {
	t.Parallel()

	e := NewDictEntryTupleExpr(*parser.NewScanner(""), NewNumber(1), NewNumber(2))
	val, err := e.Eval(arraictx.InitRunCtx(context.Background()), EmptyScope)

	assert.Equal(t, NewDictEntryTuple(NewNumber(1), NewNumber(2)), val)
	assert.NoError(t, err)
}

func TestDictEntryTupleExprEvalErrorOnAtEvalError(t *testing.T) {
	t.Parallel()

	// On Eval, this will return an error that will be propagated.
	badExpr := NewDotExpr(*parser.NewScanner(""), NewNumber(0), "*")
	_, err := badExpr.Eval(arraictx.InitRunCtx(context.Background()), EmptyScope)
	msg := strings.Split(err.Error(), "\n")[0]

	e := NewDictEntryTupleExpr(*parser.NewScanner(""), badExpr, NewNumber(2))

	AssertExprErrorEquals(t, e, msg)
}

func TestDictEntryTupleExprEvalErrorOnValueEvalError(t *testing.T) {
	t.Parallel()

	// On Eval, this will return an error that will be propagated.
	badExpr := NewDotExpr(*parser.NewScanner(""), NewNumber(0), "*")
	_, err := badExpr.Eval(arraictx.InitRunCtx(context.Background()), EmptyScope)
	msg := strings.Split(err.Error(), "\n")[0]

	e := NewDictEntryTupleExpr(*parser.NewScanner(""), NewNumber(1), badExpr)

	AssertExprErrorEquals(t, e, msg)
}
