package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

type RecursionExpr struct {
	ExprScanner
	name      ExprPattern
	fn        Expr
	fix, fixt Value
}

func NewRecursionExpr(scanner parser.Scanner, name Expr, fn Expr, fix, fixt Value) Expr {
	return RecursionExpr{ExprScanner{scanner}, NewExprPattern(name), fn, fix, fixt}
}

func (r RecursionExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	val, err := r.fn.Eval(ctx, local)
	if err != nil {
		return nil, WrapContext(err, r, local)
	}

	//TODO: optimise, get it to load either fix or fixt not both
	switch f := val.(type) {
	case Tuple:
		t := f
		for e := f.Enumerator(); e.MoveNext(); {
			attr, val := e.Current()
			if fn, isFunction := val.(Closure); isFunction {
				f = f.With(attr, NewClosure(local, NewFunction(fn.Source(), r.name, fn.f).(*Function)))
				continue
			}
			return nil, WrapContext(errors.Errorf("Recursion requires a tuple of functions: %v", t.String()), r, local)
		}
		return Call(r.fixt, f, local)
	case Closure:
		return Call(r.fix, NewClosure(local, NewFunction(f.Source(), r.name, f.f).(*Function)), local)
	}
	return nil, WrapContext(errors.Errorf("Recursion does not support %T", val), r, local)
}

func (r RecursionExpr) String() string {
	return fmt.Sprintf("\\%s %s", r.name, r.fn)
}
