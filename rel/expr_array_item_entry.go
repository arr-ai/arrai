package rel

import (
	"context"
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
	return ArrayItemTupleExpr{ExprScanner: ExprScanner{Src: scanner}, at: at, item: value}
}

// String returns a string representation of the expression.
func (e ArrayItemTupleExpr) String() string {
	return fmt.Sprintf("(@: %v, %s: %v)", e.at, ArrayItemAttr, e.item)
}

// Eval returns the subject.
func (e ArrayItemTupleExpr) Eval(ctx context.Context, local Scope) (_ Value, err error) {
	at, err := e.at.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	value, err := e.item.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	return NewArrayItemTuple(int(at.(Number).Float64()), value), nil
}
