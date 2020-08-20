package rel

import (
	"bytes"
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// CondPatternControlVarExpr which is used for `cond` pattern matching.
type CondPatternControlVarExpr struct {
	ExprScanner
	controlVarExpr Expr
	conditionPairs []PatternExprPair
}

// NewCondPatternControlVarExpr returns a new CondPatternControlVarExpr.
func NewCondPatternControlVarExpr(scanner parser.Scanner, controlVar Expr, patternExprs ...PatternExprPair) Expr {
	return CondPatternControlVarExpr{ExprScanner{scanner}, controlVar, patternExprs}
}

func (expr CondPatternControlVarExpr) String() string {
	var b bytes.Buffer
	b.WriteByte('(')
	fmt.Fprintf(&b, "(control_var: %v)", expr.controlVarExpr.String())

	if len(expr.conditionPairs) > 0 {
		b.WriteByte(',')
	}

	b.WriteByte('{')
	for i, conditionPair := range expr.conditionPairs {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%v", conditionPair.String())
	}

	b.WriteByte('}')
	b.WriteByte(')')
	return b.String()
}

// Eval evaluates to find the first valid condition and return its value.
func (expr CondPatternControlVarExpr) Eval(ctx context.Context, scope Scope) (Value, error) {
	varVal, err := expr.controlVarExpr.Eval(ctx, scope)
	if err != nil {
		return nil, WrapContext(err, expr.controlVarExpr, scope)
	}

	for _, conditionPair := range expr.conditionPairs {
		bindings, err := conditionPair.Bind(ctx, scope, varVal)
		if err == nil {
			l, err := scope.MatchedUpdate(bindings)
			if err != nil {
				return nil, WrapContext(err, expr.controlVarExpr, scope)
			}
			val, err := conditionPair.Eval(ctx, l)
			if err != nil {
				return nil, WrapContext(err, expr.controlVarExpr, l)
			}
			return val, nil
		}
	}

	return None, nil
}
