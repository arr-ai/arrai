package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// DictEntryTupleExpr represents an expr that evaluates to a DictEntryTuple.
type DictEntryTupleExpr struct {
	ExprScanner
	at, value Expr
}

// NewDictEntryTupleExpr returns a new dictEntryTupleExpr.
func NewDictEntryTupleExpr(scanner parser.Scanner, at, value Expr) DictEntryTupleExpr {
	// TODO: Optimise for literals.
	// if at, ok := at.(Value); ok {
	// 	if value, ok := value.(Value); ok {
	// 		return NewDictTuple(at, value)
	// 	}
	// }
	return DictEntryTupleExpr{at: at, value: value}
}

// String returns a string representation of the expression.
func (e DictEntryTupleExpr) String() string {
	return fmt.Sprintf("(@: %v, %s: %v)", e.at, DictValueAttr, e.value)
}

// Eval returns the subject
func (e DictEntryTupleExpr) Eval(local Scope) (Value, error) {
	at, err := e.at.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e)
	}
	value, err := e.value.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e)
	}
	return NewDictEntryTuple(at, value), nil
}
