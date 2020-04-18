package rel

import (
	"fmt"
	"strings"
)

type CompareFunc func(a, b Value) bool

// CompareExpr represents a range of operators.
type CompareExpr struct {
	args  []Expr
	comps []CompareFunc
	ops   []string
}

func NewCompareExpr(args []Expr, comps []CompareFunc, ops []string) CompareExpr {
	return CompareExpr{args: args, comps: comps, ops: ops}
}

// String returns a string representation of the expression.
func (e CompareExpr) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "(%v ", e.args[0])
	for i, arg := range e.args[1:] {
		fmt.Fprintf(&sb, " %s %v", e.ops[i], arg)
	}
	sb.WriteString(")")
	return sb.String()
}

// Eval returns the subject
func (e CompareExpr) Eval(local Scope) (Value, error) {
	lhs, err := e.args[0].Eval(local)
	if err != nil {
		return nil, err
	}
	for i, arg := range e.args[1:] {
		rhs, err := arg.Eval(local)
		if err != nil {
			return nil, err
		}
		if !e.comps[i](lhs, rhs) {
			return False, nil
		}
		lhs = rhs
	}
	return True, nil
}
