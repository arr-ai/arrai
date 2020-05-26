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
	conditions     []interface{}
	values         []Expr
}

// NewCondPatternControlVarExpr returns a new CondPatternControlVarExpr.
func NewCondPatternControlVarExpr(scanner parser.Scanner, controlVar Expr, conditions []interface{},
	values []Expr) Expr {
	return &CondPatternControlVarExpr{ExprScanner{scanner}, controlVar, conditions, values}
}

func (expr *CondPatternControlVarExpr) String() string {
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
		switch cond := cond.(type) {
		case Expr:
			fmt.Fprintf(&b, "%v: %v", cond.String(), expr.values[i].String())
		case Pattern:
			fmt.Fprintf(&b, "%v: %v", cond.String(), expr.values[i].String())
		}
	}

	b.WriteByte('}')
	b.WriteByte(')')
	return b.String()
}

// Eval evaluates with pattern matching.
func (expr *CondPatternControlVarExpr) Eval(local Scope) (Value, error) {
	varVal, err := expr.controlVarExpr.Eval(local)
	if err != nil {
		return nil, wrapContext(err, expr.controlVarExpr)
	}

	for index, condition := range expr.conditions {
		switch condition := condition.(type) {
		case Expr:
			switch condition := condition.(type) {
			case IdentExpr:
				if condition.String() == "_" {
					val, err := expr.values[index].Eval(local)
					if err != nil {
						return nil, wrapContext(err, condition)
					}
					return val, nil
				}
			case Array:
				for _, exprVal := range condition.Values() {
					if exprVal.Equal(varVal) {
						val, err := expr.values[index].Eval(local)
						if err != nil {
							return nil, wrapContext(err, condition)
						}
						return val, nil
					}
				}
			case ArrayExpr:
				for _, exprVal := range condition.Elements() {
					val, err := exprVal.Eval(local)
					if err != nil {
						return nil, wrapContext(err, condition)
					}
					if val.Equal(varVal) {
						val, err := expr.values[index].Eval(local)
						if err != nil {
							return nil, wrapContext(err, condition)
						}
						return val, nil
					}
				}
			default:
				cond, err := condition.Eval(local)
				if err != nil {
					return nil, wrapContext(err, condition)
				}
				if varVal.Equal(cond) {
					val, err := expr.values[index].(Expr).Eval(local)
					if err != nil {
						return nil, wrapContext(err, condition)
					}
					return val, nil
				}
			}
		case Pattern:
			// TODO: now binding can't check types, see this case `let a = {"a":3}; a cond {(a:x): x + 5,_:2}`
			// It will panic and stop the process, it is not good.
			local = condition.Bind(local, varVal)
			val, err := expr.values[index].Eval(local)
			if err != nil {
				return nil, err
			}
			return val, nil
		}
	}

	return None, nil
}
