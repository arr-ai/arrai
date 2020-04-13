package rel

// SwitchExpr returns the tuple applied to a function.
type SwitchExpr struct {
}

// NewSwitchExpr returns a new SwitchExpr.
func NewSwitchExpr() Expr {
	return nil
}

// Eval returns the true condition
func (e *SwitchExpr) Eval(local Scope) (Value, error) {
	return nil, nil
}
