package rel

import (
	"github.com/arr-ai/wbnf/parser"
)

// CondPatternControlVarExpr
type CondPatternControlVarExpr struct {
	ExprScanner
	controlVarExpr Expr
	conditions     []interface{}
	values         []Expr
}

// NewCondPatternControlVarExpr returns a new CondPatternControlVarExpr.
func NewCondPatternControlVarExpr(scanner parser.Scanner, controlVar Expr, conditions []interface{}, values []Expr) Expr {
	return &CondPatternControlVarExpr{ExprScanner{scanner}, controlVar, conditions, values}
}

func (expr *CondPatternControlVarExpr) String() string {
	return ""
}

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
					val, err := expr.values[index].Eval(local)
					if err != nil {
						return nil, wrapContext(err, condition)
					}
					return val, nil
				}
			}
		case Pattern:
			condition.Bind(local, nil)
		}
	}
	return None, nil
}
