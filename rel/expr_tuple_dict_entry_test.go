package rel

import (
	"github.com/arr-ai/wbnf/parser"
	"testing"

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
	val, _ := e.Eval(EmptyScope)

	assert.Equal(t, NewDictEntryTuple(NewNumber(1), NewNumber(2)), val)
}

func TestDictEntryTupleExprEvalErrorOnAtEvalError(t *testing.T) {
	t.Parallel()

	// On Eval, this will return an error that will be propagated.
	badExpr := NewDotExpr(*parser.NewScanner(""), NewNumber(0), "*")
	_, err := badExpr.Eval(EmptyScope)

	e := NewDictEntryTupleExpr(*parser.NewScanner(""), badExpr, NewNumber(2))

	AssertExprErrorEquals(t, e, err.Error())
}

func TestDictEntryTupleExprEvalErrorOnValueEvalError(t *testing.T) {
	t.Parallel()

	// On Eval, this will return an error that will be propagated.
	badExpr := NewDotExpr(*parser.NewScanner(""), NewNumber(0), "*")
	_, err := badExpr.Eval(EmptyScope)

	e := NewDictEntryTupleExpr(*parser.NewScanner(""), NewNumber(1), badExpr)

	AssertExprErrorEquals(t, e, err.Error())
}
