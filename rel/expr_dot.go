package rel

import (
	"fmt"

	"github.com/go-errors/errors"
)

// DotExpr returns the tuple or set with a single field replaced by an
// expression.
type DotExpr struct {
	lhs  Expr
	attr string
}

// NewDotExpr returns a new DotExpr that fetches the given attr from the
// lhs, which is expected to be a tuple.
func NewDotExpr(lhs Expr, attr string) Expr {
	return &DotExpr{lhs, attr}
}

// Subject returns the DotExpr's LHS.
func (e *DotExpr) Subject() Expr {
	return e.lhs
}

// Attr returns the name of the attribute accessed by the DotExpr.
func (e *DotExpr) Attr() string {
	return e.attr
}

// String returns a string representation of the expression.
func (e *DotExpr) String() string {
	return fmt.Sprintf("(%s.%s)", e.lhs, e.attr)
}

// Eval returns the lhs
func (e *DotExpr) Eval(local Scope) (Value, error) {
	if e.attr == "*" {
		return nil, errors.Errorf("expr.* not allowed outside tuple attr")
	}
	a, err := e.lhs.Eval(local)
	if err != nil {
		return nil, err
	}
	get := func(t Tuple) (Value, error) {
		if value, found := t.Get(e.attr); found {
			return value, nil
		}
		if e.attr[:1] != "&" {
			if value, found := t.Get("&" + e.attr); found {
				tupleScope := local.With("self", t)
				return value.(*Function).Call(nil, tupleScope)
			}
		}
		return nil, errors.Errorf("Missing attr %s", e.attr)
	}

	switch x := a.(type) {
	case Tuple:
		return get(x)
	case Set:
		if !x.Bool() {
			return nil, errors.Errorf("Cannot get attr from empty set")
		}
		e := x.Enumerator()
		e.MoveNext()
		v := e.Current()
		if e.MoveNext() {
			return nil, errors.Errorf("Too many elts to get attr from set")
		}
		if t, ok := v.(Tuple); ok {
			return get(t)
		}
		return nil, errors.Errorf(
			"Cannot get attr from non-tuple set elt %s", v)
	default:
		return nil, errors.Errorf(
			"(%s).%s: lhs must be a Tuple, not %T", e.lhs, e.attr, a)
	}
}
