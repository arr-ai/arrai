package rel

import (
	"bytes"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// CondPatternControlVarExpr which is used for `cond` pattern matching.
type CondPatternControlVarExpr struct {
	ExprScanner
	controlVarExpr Expr
	conditions     []Pattern
	values         []Expr
}

// NewCondPatternControlVarExpr returns a new CondPatternControlVarExpr.
func NewCondPatternControlVarExpr(scanner parser.Scanner, controlVar Expr, conditions []Pattern,
	values []Expr) Expr {
	return CondPatternControlVarExpr{ExprScanner{scanner}, controlVar, conditions, values}
}

func (expr CondPatternControlVarExpr) String() string {
	var b bytes.Buffer
	b.WriteByte('(')
	fmt.Fprintf(&b, "(control_var: %v)", expr.controlVarExpr.String())

	if len(expr.conditions) > 0 {
		b.WriteByte(',')
	}

	b.WriteByte('{')
	for i, cond := range expr.conditions {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%v: %v", cond.String(), expr.values[i].String())
	}

	b.WriteByte('}')
	b.WriteByte(')')
	return b.String()
}

// Eval evaluates to find the first valid condition and return its value.
func (expr CondPatternControlVarExpr) Eval(local Scope) (Value, error) {
	varVal, err := expr.controlVarExpr.Eval(local)
	if err != nil {
		return nil, wrapContext(err, expr.controlVarExpr)
	}

	for cIndex, condition := range expr.conditions {
		local, err = condition.Bind(local, varVal)
		if err == nil {
			val, err := expr.values[cIndex].Eval(local)
			if err != nil {
				return nil, err
			}
			return val, nil
		}
	}

	return None, nil
}
