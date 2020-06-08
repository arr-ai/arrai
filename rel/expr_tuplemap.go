package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// TupleMapExpr returns the tuple applied to a function.
type TupleMapExpr struct {
	ExprScanner
	lhs Expr
	fn  *Function
}

// NewAngleArrowExpr returns a new AtArrowExpr.
func NewTupleMapExpr(scanner parser.Scanner, lhs Expr, fn Expr) Expr {
	return &TupleMapExpr{ExprScanner{scanner}, lhs, ExprAsFunction(fn)}
}

// String returns a string representation of the expression.
func (e *TupleMapExpr) String() string {
	return fmt.Sprintf("(%s :> %s)", e.lhs, e.fn)
}

// Eval returns the lhs
func (e *TupleMapExpr) Eval(local Scope) (_ Value, err error) {
	value, err := e.lhs.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e, local)
	}
	defer func() {
		if r := recover(); r != nil {
			if err == nil {
				panic(r)
			}
		}
	}()
	return value.(Tuple).Map(func(v Value) Value {
		scope, _ := e.fn.arg.Bind(local, v) //nolint: errcheck
		v, err = e.fn.body.Eval(local.Update(scope))
		if err != nil {
			panic(err)
		}
		return v
	}), nil
}
