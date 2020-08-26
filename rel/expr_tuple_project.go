package rel

import (
	"context"
	"fmt"
	"strings"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

type TupleProjectExpr struct {
	ExprScanner
	base    Expr
	inverse bool
	attrs   Names
}

func NewTupleProjectExpr(scanner parser.Scanner, base Expr, inverse bool, attrs []string) Expr {
	return &TupleProjectExpr{ExprScanner{scanner}, base, inverse, NewNames(attrs...)}
}

func (tp *TupleProjectExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	val, err := tp.base.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, tp, local)
	}
	tuple, isTuple := val.(Tuple)
	if !isTuple {
		return nil, WrapContextErr(errors.Errorf("lhs does not evaluate to tuple: %s", val), tp, local)
	}
	if !tp.attrs.IsSubsetOf(tuple.Names()) {
		return nil, WrapContextErr(errors.Errorf("names are not subset of lhs: %s", tuple.Names()), tp, local)
	}

	if tp.inverse {
		return TupleProjectAllBut(tuple, tp.attrs), nil
	}

	return tuple.Project(tp.attrs), nil
}

func (tp *TupleProjectExpr) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("(%s).", tp.base))
	if tp.inverse {
		str.WriteRune('~')
	}
	str.WriteString(tp.attrs.String())
	return str.String()
}
