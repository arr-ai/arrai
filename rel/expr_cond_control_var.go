package rel

import (
	"bytes"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// CondControlVarExpr returns the tuple applied to a function, the expression looks like:
// let a = 1 + 1; a cond (1 : 2 + 1, 2 : 5, *: 10)
// let a = 1 + 1; let b = a cond (1 : 2 + 1, 2 : 5, *: 10); b
type CondControlVarExpr struct {
	ExprScanner
	controlVarExpr Expr
	standardExpr   CondExpr
}

// NewCondControlVarExpr returns a new CondControlVarExpr.
func NewCondControlVarExpr(scanner parser.Scanner, controlVar Expr, dictExpr Expr, defaultExpr Expr) Expr {
	return &CondControlVarExpr{ExprScanner{scanner}, controlVar,
		CondExpr{ExprScanner{scanner}, dictExpr, defaultExpr, func(condition Value, local Scope) bool {
			controlVarVal, has := local.Get("controlVarVal")
			if !has {
				return false
			}
			return condition.Equal(controlVarVal)
		}}}
}

// String returns a string representation of the expression.
func (e *CondControlVarExpr) String() string {
	var b bytes.Buffer
	b.WriteByte('(')
	fmt.Fprintf(&b, "(control_var: %v)", e.controlVarExpr.String())
	standardExprStr := e.standardExpr.String()
	if len(standardExprStr) != 0 {
		b.WriteString("," + standardExprStr)
	}
	b.WriteByte(')')
	return b.String()
}

// Eval returns the value of valid condition and whose value equals to the control var value,
// or the value of default condition.
func (e *CondControlVarExpr) Eval(local Scope) (Value, error) {
	controlVarVal, err := e.controlVarExpr.Eval(local)
	if err != nil {
		return nil, err
	}

	local = local.With("controlVarVal", controlVarVal)
	return e.standardExpr.Eval(local)
}
