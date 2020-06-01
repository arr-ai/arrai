package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

type AndExpr struct {
	ExprScanner
	a, b Expr
}

func NewAndExpr(scanner parser.Scanner, a, b Expr) Expr {
	return AndExpr{ExprScanner: ExprScanner{Src: scanner}, a: a, b: b}
}

func (e AndExpr) String() string {
	return fmt.Sprintf("(%v) && (%v)", e.a, e.b)
}

func (e AndExpr) Eval(local Scope) (Value, error) {
	a, err := e.a.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e, local)
	}
	if !a.IsTrue() {
		return a, nil
	}

	b, err := e.b.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e, local)
	}

	return b, nil
}

type OrExpr struct {
	ExprScanner
	a, b Expr
}

func NewOrExpr(scanner parser.Scanner, a, b Expr) Expr {
	return OrExpr{ExprScanner: ExprScanner{Src: scanner}, a: a, b: b}
}

func (e OrExpr) String() string {
	return fmt.Sprintf("(%v) || (%v)", e.a, e.b)
}

func (e OrExpr) Eval(local Scope) (Value, error) {
	a, err := e.a.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e, local)
	}
	if a.IsTrue() {
		return a, nil
	}

	b, err := e.b.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e, local)
	}

	return b, nil
}
