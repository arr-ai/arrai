package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

const safeCallOp = "safe_call"

type SafeTailExpr struct {
	ExprScanner
	fallbackValue, base Expr
	tailExprs           []func(Expr) Expr
}

func NewSafeTailExpr(
	scanner parser.Scanner,
	fallback, base Expr,
	tailExprs []func(Expr) Expr,
) Expr {
	if len(tailExprs) == 0 {
		panic("exprs cannot be empty")
	}
	return &SafeTailExpr{ExprScanner{scanner}, fallback, base, tailExprs}
}

func (s *SafeTailExpr) Eval(local Scope) (value Value, err error) {
	value, err = s.base.Eval(local)
	if err != nil {
		return nil, wrapContext(err, s)
	}
	for _, t := range s.tailExprs {
		expr := t(value)
		if call, isCall := expr.(*BinExpr); isCall && call.op == "safe_call" {
			value, err = call.Eval(local)
			if err != nil {
				return nil, wrapContext(err, s)
			}

			for e, i := value.(Set).Enumerator(), 1; e.MoveNext(); i++ {
				if i > 1 {
					return s.fallbackValue.Eval(local)
				}
			}
			if !value.IsTrue() {
				return s.fallbackValue.Eval(local)
			}
			value = SetAny(value.(Set))
		} else if safeDot, isSafeDot := expr.(*SafeDotExpr); isSafeDot {
			value, err = safeDot.Eval(local)
			if err != nil {
				if _, isMissingAttr := err.(missingAttrError); isMissingAttr {
					return s.fallbackValue.Eval(local)
				}
				return nil, wrapContext(err, s)
			}
		} else {
			value, err = expr.Eval(local)
			if err != nil {
				return nil, wrapContext(err, s)
			}
		}
	}
	return value, err
}

func (s *SafeTailExpr) String() string {
	finalExpr := s.tailExprs[0](s.base)
	if len(s.tailExprs) > 1 {
		for _, e := range s.tailExprs[1:] {
			finalExpr = e(finalExpr)
		}
	}
	return finalExpr.String() + fmt.Sprintf(":%s", s.fallbackValue.String())
}
