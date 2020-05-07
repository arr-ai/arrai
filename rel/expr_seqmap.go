package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// SequenceMapExpr returns the tuple applied to a function.
type SequenceMapExpr struct {
	ExprScanner
	lhs Expr
	fn  *Function
}

// NewAngleArrowExpr returns a new AtArrowExpr.
func NewSequenceMapExpr(scanner parser.Scanner, lhs Expr, fn Expr) Expr {
	return &SequenceMapExpr{ExprScanner{scanner}, lhs, ExprAsFunction(fn)}
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
		return nil, wrapContext(err, e)
	}
	// TODO: implement directly for String, Array and Dict.
	if set, ok := value.(Set); ok {
		values := []Value{}
		for i := set.Enumerator(); i.MoveNext(); {
			t := i.Current().(Tuple)
			pos, _ := t.Get("@")
			attr := t.Names().Without("@").Any()
			item, _ := t.Get(attr)
			v, err := e.fn.body.Eval(local.With(e.fn.Arg(), item))
			if err != nil {
				return nil, wrapContext(err, e)
			}
			values = append(values, NewTuple(Attr{"@", pos}, Attr{attr, v}))
		}
		return NewSet(values...), nil
	}
	return nil, wrapContext(errors.Errorf(">> not applicable to %T", value), e)
}
