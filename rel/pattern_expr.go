package rel

import (
	"bytes"
	"context"
	"fmt"
)

type ExprPattern struct {
	Expr Expr
}

func NewExprPattern(expr Expr) Pattern {
	switch x := expr.(type) {
	case IdentExpr:
		return IdentPattern(x.ident)
	case DynIdentExpr:
		return DynIdentPattern(x.ident)
	}
	if value, is := exprIsValue(expr); is {
		return ExprPattern{Expr: value}
	}
	return ExprPattern{Expr: expr}
}

func (p ExprPattern) Bind(ctx context.Context, scope Scope, value Value, b *scopeBuilder) (context.Context, error) {
	if identExpr, is := p.Expr.(IdentExpr); is {
		// Bind value for identexpr in Pattern, like `let (a: x, b: y) = (a: 4, b: 7); x`
		return ctx, b.add(identExpr.ident, value)
	}

	v, err := p.Expr.Eval(ctx, scope)
	if err != nil {
		return ctx, err
	}
	if v.Equal(value) {
		return ctx, nil
	}
	if !b.explain {
		return ctx, errPatternMismatch
	}
	return ctx, lazyErrorf("no match: %v != %v", v, value)
}

func (p ExprPattern) String() string {
	return p.Expr.String()
}

// Bindings returns nil: ExprPattern matches by value equality and binds no
// names (NewExprPattern routes an IdentExpr to IdentPattern instead).
func (p ExprPattern) Bindings() []string {
	return nil
}

type ExprsPattern struct {
	exprs []Expr
}

func NewExprsPattern(exprs ...Expr) ExprsPattern {
	return ExprsPattern{exprs: exprs}
}

func (p ExprsPattern) Bind(ctx context.Context, scope Scope, value Value, b *scopeBuilder) (context.Context, error) {
	if len(p.exprs) == 0 {
		if !b.explain {
			return ctx, errPatternMismatch
		}
		return ctx, lazyErrorf("there is not any rel.Expr in rel.ExprsPattern")
	}

	incomingVal, err := value.Eval(ctx, scope)
	if err != nil {
		return ctx, err
	}

	for _, e := range p.exprs {
		val, err := e.Eval(ctx, scope)
		if err != nil {
			return ctx, err
		}
		if incomingVal.Equal(val) {
			return ctx, nil
		}
	}

	if !b.explain {
		return ctx, errPatternMismatch
	}
	return ctx, lazyErrorf("didn't find matched value")
}

func (p ExprsPattern) String() string {
	if len(p.exprs) == 0 {
		panic("there is not any rel.Expr in rel.ExprsPattern")
	}

	if len(p.exprs) == 1 {
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

// Bindings returns nil: ExprsPattern matches by value equality against one
// of several alternatives and binds no names.
func (p ExprsPattern) Bindings() []string {
	return nil
}
