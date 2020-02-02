package rel

import (
	"fmt"

	"github.com/go-errors/errors"
)

// SequenceMapExpr returns the tuple applied to a function.
type SequenceMapExpr struct {
	lhs Expr
	fn  *Function
}

// NewAngleArrowExpr returns a new AtArrowExpr.
func NewSequenceMapExpr(lhs Expr, fn Expr) Expr {
	return &SequenceMapExpr{lhs, ExprAsFunction(fn)}
}

// LHS returns the LHS of the AtArrowExpr.
func (e *SequenceMapExpr) LHS() Expr {
	return e.lhs
}

// Fn returns the function to be applied to the LHS.
func (e *SequenceMapExpr) Fn() *Function {
	return e.fn
}

// String returns a string representation of the expression.
func (e *SequenceMapExpr) String() string {
	return fmt.Sprintf("(%s >> %s)", e.lhs, e.fn)
}

// Eval returns the lhs
func (e *SequenceMapExpr) Eval(local Scope) (Value, error) {
	value, err := e.lhs.Eval(local)
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
			v, err := e.fn.body.Eval(local.With(e.fn.arg, item))
			if err != nil {
				return nil, err
			}
			result = result.With(NewTuple(Attr{"@", pos}, Attr{attr, v}))
		}
		return result, nil
	}
	return nil, errors.Errorf("=> not applicable to %T", value)
}
