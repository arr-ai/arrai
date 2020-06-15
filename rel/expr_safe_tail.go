package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

type SafeTailCallback func(Value, Scope) (Value, error)

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

func (s *SafeTailExpr) Eval(local Scope) (value Value, err error) {
	value, err = s.base.Eval(local)
	if err != nil {
		return nil, WrapContext(err, s, local)
	}
	for _, t := range s.tailExprs {
		value, err = t(value, local)
		if err != nil {
			return nil, WrapContext(err, s, local)
		}
		if value == nil {
			return s.fallbackValue.Eval(local)
		}
	}
	return
}

func (s *SafeTailExpr) String() string {
	//FIXME: printing not very descriptive
	return fmt.Sprintf("%s...TODO...:%s", s.base.String(), s.fallbackValue.String())
}
