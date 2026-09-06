package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// SafeTailStep is one ?.attr or ?(arg) in a safe-tail chain (🎯T25).
// Get is true for attribute access (attr may be empty: .”).
type SafeTailStep struct {
	safe bool
	get  bool
	attr string
	args []Expr
}

func NewSafeTailGet(safe bool, attr string) SafeTailStep {
	return SafeTailStep{safe: safe, get: true, attr: attr}
}

func NewSafeTailCall(safe bool, args ...Expr) SafeTailStep {
	return SafeTailStep{safe: safe, args: args}
}

type SafeTailExpr struct {
	ExprScanner
	fallbackValue, base Expr
	steps               []SafeTailStep
}

func NewSafeTailExpr(scanner parser.Scanner, fallback, base Expr, steps []SafeTailStep) Expr {
	if len(steps) == 0 {
		panic("exprs cannot be empty")
	}
	return &SafeTailExpr{ExprScanner{scanner}, fallback, base, steps}
}

func (s *SafeTailExpr) Eval(ctx context.Context, local Scope) (value Value, err error) {
	value, err = s.base.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, s, local)
	}
	for i := range s.steps {
		value, err = s.steps[i].apply(ctx, value, local)
		if err != nil {
			return nil, WrapContextErr(err, s, local)
		}
		if value == nil {
			return s.fallbackValue.Eval(ctx, local)
		}
	}
	return
}

func (s SafeTailStep) apply(ctx context.Context, v Value, local Scope) (Value, error) {
	var err error
	if s.get {
		v, err = NewDotExpr(v.Source(), v, s.attr).Eval(ctx, local)
	} else {
		for _, arg := range s.args {
			a, aerr := arg.Eval(ctx, local)
			if aerr != nil {
				err = aerr
				break
			}
			set, is := v.(Set)
			if !is {
				return nil, fmt.Errorf("not a set: %v", v)
			}
			v, err = SetCall(ctx, set, a)
			if err != nil {
				break
			}
		}
	}
	if err != nil && s.safe && isSafeTailMiss(err) {
		return nil, nil
	}
	return v, err
}

func isSafeTailMiss(err error) bool {
	if _, ok := err.(NoReturnError); ok {
		return true
	}
	if e, ok := err.(ContextErr); ok {
		_, is := e.NextErr().(MissingAttrError)
		return is
	}
	return false
}

func (s *SafeTailExpr) String() string {
	return fmt.Sprintf("%s...TODO...:%s", s.base.String(), s.fallbackValue.String())
}
