package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

type MissingAttrError struct {
	ctxErr error
}

func (m MissingAttrError) Error() string {
	return m.ctxErr.Error()
}

// DotExpr returns the tuple or set with a single field replaced by an
// expression.
type DotExpr struct {
	ExprScanner
	lhs  Expr
	attr string
}

// NewDotExpr returns a new DotExpr that fetches the given attr from the
// lhs, which is expected to be a tuple.
func NewDotExpr(scanner parser.Scanner, lhs Expr, attr string) Expr {
	return &DotExpr{ExprScanner{scanner}, lhs, attr}
}

// Subject returns the DotExpr's LHS.
func (x *DotExpr) Subject() Expr {
	return x.lhs
}

// Attr returns the name of the attribute accessed by the DotExpr.
func (x *DotExpr) Attr() string {
	return x.attr
}

// String returns a string representation of the expression.
func (x *DotExpr) String() string {
	return fmt.Sprintf("(%s.%s)", x.lhs, x.attr)
}

// Eval returns the lhs
func (x *DotExpr) Eval(ctx context.Context, local Scope) (_ Value, err error) {
	if x.attr == "*" {
		return nil, WrapContextErr(errors.Errorf("expr.* not allowed outside tuple attr"), x, local)
	}
	a, err := x.lhs.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, x, local)
	}
	get := func(ctx context.Context, t Tuple) (Value, error) {
		if value, found := t.Get(x.attr); found {
			return value, nil
		}
		if len(x.attr) > 0 && x.attr[:1] != "&" {
			if value, found := t.Get("&" + x.attr); found {
				//TODO: add tupleScope self to allow accessing itself
				switch f := value.(type) {
				case Closure:
					return SetCall(ctx, f, nil)
				case *NativeFunction:
					return SetCall(ctx, f, nil)
				default:
					panic(fmt.Errorf("not a function: %v", f))
				}
			}
		}
		return nil, WrapContextErr(MissingAttrError{errors.Errorf("Missing attr %q (available: %v)", x.attr, t.Names())},
			x, local)
	}

	switch t := a.(type) {
	case Tuple:
		return get(ctx, t)
	case Set:
		if !t.IsTrue() {
			return nil, WrapContextErr(errors.Errorf("Cannot get attr %q from empty set", x.attr), x, local)
		}
		e := t.Enumerator()
		e.MoveNext()
		v := e.Current()
		if e.MoveNext() {
			return nil, WrapContextErr(errors.Errorf("Too many elts to get attr %q from set", x.attr), x, local)
		}
		if t, ok := v.(Tuple); ok {
			return get(ctx, t)
		}
		return nil, WrapContextErr(errors.Errorf("Cannot get attr %q from non-tuple set elt", x.attr), x, local)
	default:
		return nil, WrapContextErr(errors.Errorf("(%s).%s: lhs must be a tuple, not %s",
			x.lhs, x.attr, ValueTypeAsString(a)), x, local)
	}
}
