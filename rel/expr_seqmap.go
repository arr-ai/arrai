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

// NewSequenceMapExpr returns a new SequenceMapExpr.
func NewSequenceMapExpr(scanner parser.Scanner, lhs Expr, fn Expr) Expr {
	return &SequenceMapExpr{ExprScanner{scanner}, lhs, ExprAsFunction(fn)}
}

// String returns a string representation of the expression.
func (e *SequenceMapExpr) String() string {
	return fmt.Sprintf("(%s >> %s)", e.lhs, e.fn)
}

// Eval returns the lhs
func (e *SequenceMapExpr) Eval(local Scope) (_ Value, err error) {
	defer wrapPanic(e, &err, local)
	value, err := e.lhs.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e, local)
	}
	// TODO: implement directly for String, Array and Dict.
	if set, ok := value.(Set); ok {
		values := []Value{}
		for i := set.Enumerator(); i.MoveNext(); {
			t := i.Current().(Tuple)
			pos, _ := t.Get("@")
			attr := t.Names().Without("@").Any()
			item, _ := t.Get(attr)
			scope, err := e.fn.arg.Bind(local, item)
			if err != nil {
				return nil, wrapContext(err, e, local)
			}
			v, err := e.fn.body.Eval(local.Update(scope))
			if err != nil {
				return nil, wrapContext(err, e, local)
			}
			values = append(values, NewTuple(Attr{"@", pos}, Attr{attr, v}))
		}
		return NewSet(values...), nil
	}
	return nil, wrapContext(errors.Errorf(">> not applicable to %T", value), e, local)
}
