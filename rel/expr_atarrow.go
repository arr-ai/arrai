package rel

import (
	"fmt"

	"github.com/go-errors/errors"
)

// AtArrowExpr returns the tuple applied to a function.
type AtArrowExpr struct {
	lhs Expr
	fn  *Function
}

// NewAngleArrowExpr returns a new AtArrowExpr.
func NewAngleArrowExpr(lhs Expr, fn Expr) Expr {
	return &AtArrowExpr{lhs, ExprAsFunction(fn)}
}

// LHS returns the LHS of the AtArrowExpr.
func (e *AtArrowExpr) LHS() Expr {
	return e.lhs
}

// Fn returns the function to be applied to the LHS.
func (e *AtArrowExpr) Fn() *Function {
	return e.fn
}

// String returns a string representation of the expression.
func (e *AtArrowExpr) String() string {
	return fmt.Sprintf("(%s @> %s)", e.lhs, e.fn)
}

// Eval returns the lhs
func (e *AtArrowExpr) Eval(local, global Scope) (Value, error) {
	value, err := e.lhs.Eval(local, global)
	if err != nil {
		return nil, err
	}
	if set, ok := value.(Set); ok {
		result := NewSet()
		for i := set.Enumerator(); i.MoveNext(); {
			t := i.Current().(Tuple)
			pos, _ := t.Get("@")
			attr := t.Names().Without("@").Any()
			item, _ := t.Get(attr)
			v, err := e.fn.body.Eval(local.With(e.fn.arg, item), global)
			if err != nil {
				return nil, err
			}
			result = result.With(NewTuple(Attr{"@", pos}, Attr{attr, v}))
		}
		return result, nil
	}
	return nil, errors.Errorf("=> not applicable to %T", value)
}
