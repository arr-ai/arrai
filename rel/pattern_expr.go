package rel

import (
	"bytes"
	"fmt"

	"github.com/go-errors/errors"
)

type ExprPattern struct {
	Expr Expr
}

func NewExprPattern(expr Expr) ExprPattern {
	if value, is := exprIsValue(expr); is {
		return ExprPattern{Expr: value}
	}
	return ExprPattern{Expr: expr}
}

func (p ExprPattern) Bind(scope Scope, value Value) (Scope, error) {
	switch t := p.Expr.(type) {
	case IdentExpr, Number, GenericSet:
		return t.(Pattern).Bind(EmptyScope, value)
	default:
		return EmptyScope, fmt.Errorf("%s is not a Pattern", t)
	}
}

func (p ExprPattern) String() string {
	return p.Expr.String()
}

type ExprsPattern struct {
	exprs []Expr
}

func NewExprsPattern(exprs ...Expr) ExprsPattern {
	return ExprsPattern{exprs: exprs}
}

func (p ExprsPattern) Bind(scope Scope, value Value) (Scope, error) {
	if len(p.exprs) == 0 {
		return EmptyScope, errors.Errorf("there is not any rel.Expr in rel.ExprsPattern")
	}

	if pe, isPattern := p.exprs[0].(Pattern); len(p.exprs) == 1 && isPattern {
		// Support patterns IDENT and NUM
		return pe.Bind(scope, value)
	}

	incomingVal, err := value.Eval(scope)
	if err != nil {
		return EmptyScope, err
	}

	for _, e := range p.exprs {
		val, err := e.Eval(scope)
		if err != nil {
			return EmptyScope, err
		}
		if incomingVal.Equal(val) {
			return scope, nil
		}
	}

	return EmptyScope, errors.Errorf("didn't find matched value")
}

func (p ExprsPattern) String() string {
	if len(p.exprs) == 0 {
		panic("there is not any rel.Expr in rel.ExprsPattern")
	}

	if len(p.exprs) == 1 {
		// it processes cases IDENT and NUM as syntax, otherwise `let (:x) = (x: 1); x` will fail.
		return p.exprs[0].String()
	}

	var b bytes.Buffer
	b.WriteByte('[')
	for i, e := range p.exprs {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%v", e.String())
	}
	b.WriteByte(']')
	return b.String()
}
