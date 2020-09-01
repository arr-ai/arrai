package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// DynLetExpr implements "let bindings; expr", where bindings must evaluates to
// a tuple whose attributes are all dynamic names. These names will be bound
// dynamically for evaluation of expr.
type DynLetExpr struct {
	ExprScanner
	bindings Expr
	expr     Expr
}

// NewDynLetExpr returns a new DynLetExpr.
func NewDynLetExpr(scanner parser.Scanner, bindings, expr Expr) Expr {
	return &DynLetExpr{ExprScanner: ExprScanner{scanner}, bindings: bindings, expr: expr}
}

// String returns a string representation of the expression.
func (e *DynLetExpr) String() string {
	return fmt.Sprintf("(let %s; %s)", e.bindings, e.expr)
}

// Eval evaluates expr with the contents of bindings bound as dynamic variables.
func (e *DynLetExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	value, err := e.bindings.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	t, is := value.(Tuple)
	if !is {
		return nil, fmt.Errorf(`bindings not a tuple in "let bindings; expr": %v`, value)
	}
	for e := t.Enumerator(); e.MoveNext(); {
		name, value := e.Current()
		if !isDynIdent(name) {
			return nil, fmt.Errorf(`%q not a dynamic name in "let bindings; expr"`, name)
		}
		ctx = context.WithValue(ctx, DynIdent(name), value)
	}
	return e.expr.Eval(ctx, local)
}
