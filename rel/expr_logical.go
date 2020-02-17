package rel

import "fmt"

type AndExpr struct {
	a, b Expr
}

func NewAndExpr(a, b Expr) Expr {
	return AndExpr{a: a, b: b}
}

func (e AndExpr) String() string {
	return fmt.Sprintf("(%v) && (%v)", e.a, e.b)
}

func (e AndExpr) Eval(local Scope) (Value, error) {
	a, err := e.a.Eval(local)
	if err != nil {
		return nil, err
	}
	if !a.Bool() {
		return a, nil
	}

	b, err := e.b.Eval(local)
	if err != nil {
		return nil, err
	}

	return b, nil
}

type OrExpr struct {
	a, b Expr
}

func NewOrExpr(a, b Expr) Expr {
	return OrExpr{a: a, b: b}
}

func (e OrExpr) String() string {
	return fmt.Sprintf("(%v) || (%v)", e.a, e.b)
}

func (e OrExpr) Eval(local Scope) (Value, error) {
	a, err := e.a.Eval(local)
	if err != nil {
		return nil, err
	}
	if a.Bool() {
		return a, nil
	}

	b, err := e.b.Eval(local)
	if err != nil {
		return nil, err
	}

	return b, nil
}
