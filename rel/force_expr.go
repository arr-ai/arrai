package rel

import (
	"context"
)

// forcedExpr is a top-level Expr whose Eval result is Forced, so Go code
// that evaluates it gets concrete values (arr-ai/arrai#789). It only ever
// wraps the root of a program handed to Go: in-language evaluation never
// sees one, so >> fusion inside the program is unaffected.
type forcedExpr struct {
	Expr
}

// ForceResult returns an Expr that evaluates e and then Forces the result:
// see Force. It is idempotent, and e's String and Source are preserved.
// Constants (a Value or LiteralExpr) are forced once, here, and keep their
// type; any other Expr is wrapped, so its concrete type changes: use
// Unforced to inspect the original.
//
// syntax.Compile already applies this to the Expr it returns.
func ForceResult(e Expr) (Expr, error) {
	switch e := e.(type) {
	case nil, forcedExpr:
		return e, nil
	case Value:
		return Force(e)
	case LiteralExpr:
		v, changed, err := force(e.literal)
		if err != nil || !changed {
			return e, err
		}
		return NewLiteralExpr(e.Source(), v), nil
	}
	return forcedExpr{e}, nil
}

// Eval evaluates the wrapped Expr and Forces the result.
func (e forcedExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	v, err := e.Expr.Eval(ctx, local)
	if err != nil {
		return nil, err
	}
	return Force(v)
}

// Unforced returns the Expr that ForceResult wrapped, or e itself.
func Unforced(e Expr) Expr {
	if f, ok := e.(forcedExpr); ok {
		return f.Expr
	}
	return e
}

// CallForced calls f with arg, as the arr.ai expression f(arg) would, and
// Forces the result. Use it, rather than Call or SetCall, to call an arr.ai
// function from Go and inspect the result's concrete type.
func CallForced(ctx context.Context, f, arg Value) (Value, error) {
	v, err := Call(ctx, f, arg, EmptyScope)
	if err != nil {
		return nil, err
	}
	return Force(v)
}
