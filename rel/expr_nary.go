package rel

import (
	"context"
	"fmt"
	"strings"

	"github.com/arr-ai/wbnf/parser"
)

type CompareFunc func(a, b Value) (bool, error)

// CompareExpr represents a range of operators.
type CompareExpr struct {
	ExprScanner
	args  []Expr
	comps []CompareFunc
	ops   []string
}

func NewCompareExpr(scanner parser.Scanner, args []Expr, comps []CompareFunc, ops []string) CompareExpr {
	return CompareExpr{ExprScanner: ExprScanner{Src: scanner}, args: args, comps: comps, ops: ops}
}

// String returns a string representation of the expression.
func (e CompareExpr) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "(%v", e.args[0])
	for i, arg := range e.args[1:] {
		fmt.Fprintf(&sb, " %s %v", e.ops[i], arg)
	}
	sb.WriteString(")")
	return sb.String()
}

// Eval returns the subject
func (e CompareExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	lhs, err := e.args[0].Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	for i, arg := range e.args[1:] {
		rhs, err := arg.Eval(ctx, local)
		if err != nil {
			return nil, WrapContextErr(err, e, local)
		}
		sat, err := e.comps[i](lhs, rhs)
		if err != nil {
			return nil, err
		}
		if !sat {
			return False, nil
		}
		lhs = rhs
	}
	return True, nil
}
