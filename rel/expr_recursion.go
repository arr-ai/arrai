package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

type RecursionExpr struct {
	ExprScanner
	name      Pattern
	fn        Expr
	fix, fixt Value
}

func NewRecursionExpr(scanner parser.Scanner, name string, fn Expr, fix, fixt Value) Expr {
	return RecursionExpr{ExprScanner{scanner}, NewIdentPattern(name), fn, fix, fixt}
}

func (r RecursionExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	val, err := r.fn.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, r, local)
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
			return nil, WrapContextErr(errors.Errorf("Recursion requires a tuple of functions: %v", t.String()), r, local)
		}
		return Call(ctx, r.fixt, f, local)
	case Closure:
		return Call(ctx, r.fix, NewClosure(local, NewFunction(f.Source(), r.name, f.f).(*Function)), local)
	}
	return nil, WrapContextErr(errors.Errorf("Recursion does not support %s", ValueTypeAsString(val)), r, local)
}

func (r RecursionExpr) String() string {
	return fmt.Sprintf("\\%s %s", r.name, r.fn)
}
