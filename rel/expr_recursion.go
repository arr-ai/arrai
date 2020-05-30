package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

type RecursionExpr struct {
	ExprScanner
	name          Pattern
	fn, fix, fixt Expr
}

func NewRecursionExpr(scanner parser.Scanner, name Pattern, fn, fix, fixt Expr) Expr {
	return RecursionExpr{ExprScanner{scanner}, name, fn, fix, fixt}
}

func (r RecursionExpr) Eval(local Scope) (Value, error) {
	if _, isIdent := r.name.(IdentExpr); !isIdent {
		return nil, errors.Errorf("Does not evaluate to a variable name: %v", r.name)
	}

	val, err := r.fn.Eval(local)
	if err != nil {
		return nil, err
	}

	argName := ExprAsPattern(NewIdentExpr(r.Source(), r.name.String()))

	//TODO: optimise, get it to load either fix or fixt not both
	switch f := val.(type) {
	case Tuple:
		valString := f.String()
		for e := f.Enumerator(); e.MoveNext(); {
			attr, val := e.Current()
			if fn, isFunction := val.(Closure); isFunction {
				f = f.With(attr, NewClosure(local, NewFunction(fn.Source(), argName, fn.f).(*Function)))
				continue
			}
			return nil, errors.Errorf("Recursion requires a tuple of functions: %v", valString)
		}
		return NewCallExpr(r.Source(), r.fixt, f).Eval(local)
	case Closure:
		return NewCallExpr(r.Source(), r.fix, NewFunction(f.Source(), argName, f.f)).Eval(local)
	}
	return nil, errors.Errorf("Recursion does not support %T", val)
}

func (r RecursionExpr) String() string {
	return fmt.Sprintf("\\%s %s", r.name, r.fn)
}
