package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

type SafeTailCallback func(context.Context, Value, Scope) (Value, error)

type SafeTailExpr struct {
	ExprScanner
	fallbackValue, base Expr
	tailExprs           []SafeTailCallback
}

func NewSafeTailExpr(scanner parser.Scanner, fallback, base Expr, tailExprs []SafeTailCallback) Expr {
	if len(tailExprs) == 0 {
		panic("exprs cannot be empty")
	}
	return &SafeTailExpr{ExprScanner{scanner}, fallback, base, tailExprs}
}

func (s *SafeTailExpr) Eval(ctx context.Context, local Scope) (value Value, err error) {
	value, err = s.base.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, s, local)
	}
	for _, t := range s.tailExprs {
		value, err = t(ctx, value, local)
		if err != nil {
			return nil, WrapContextErr(err, s, local)
		}
		if value == nil {
			return s.fallbackValue.Eval(ctx, local)
		}
	}
	return
}

func (s *SafeTailExpr) String() string {
	//FIXME: printing not very descriptive
	return fmt.Sprintf("%s...TODO...:%s", s.base.String(), s.fallbackValue.String())
}
