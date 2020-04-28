package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// ArrayItemTupleExpr represents an expr that evaluates to an ArrayItemTuple.
type ArrayItemTupleExpr struct {
	ExprScanner
	at, item Expr
}

// NewArrayItemTupleExpr returns a new ArrayItemTupleExpr.
func NewArrayItemTupleExpr(scanner parser.Scanner, at, value Expr) ArrayItemTupleExpr {
	// TODO: Optimise for literals.
	// if at, ok := at.(Value); ok {
	// 	if value, ok := value.(Value); ok {
	// 		return NewDictTuple(at, value)
	// 	}
	// }
	return ArrayItemTupleExpr{at: at, item: value}
}

// String returns a string representation of the expression.
func (e ArrayItemTupleExpr) String() string {
	return fmt.Sprintf("(@: %v, %s: %v)", e.at, ArrayItemAttr, e.item)
}

// Eval returns the subject.
func (e ArrayItemTupleExpr) Eval(local Scope) (Value, error) {
	at, err := e.at.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e)
	}
	value, err := e.item.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e)
	}
	return NewArrayItemTuple(int(at.(Number).Float64()), value), nil
}
