package rel

// CondControlVarExpr returns the tuple applied to a function, the expression looks like:
// let a = 1 + 1; cond (1 : 2 + 1, 2 : 5, *: 10).
type CondControlVarExpr struct {
	controlVarExpr Expr
	standardExpr   CondExpr
}

// NewCondControlVarExpr returns a new normal CondExpr.
func NewCondControlVarExpr(controlVar Expr, dictExpr Expr, defaultExpr Expr) Expr {
	return &CondControlVarExpr{controlVar, CondExpr{dictExpr, defaultExpr, func(condition Value, local Scope) bool {
		controlVarVal, has := local.Get("controlVarVal")
		if !has {
			return false
		}
		return condition.Equal(controlVarVal)
	}}}
}

// String returns a string representation of the expression.
func (e *CondControlVarExpr) String() string {
	return ""
}

// Eval returns the true condition. It must have only one true condition.
func (e *CondControlVarExpr) Eval(local Scope) (Value, error) {
	controlVarVal, err := e.controlVarExpr.Eval(local)
	if err != nil {
		return nil, err
	}

	local = local.With("controlVarVal", controlVarVal)
	return e.standardExpr.Eval(local)
}
