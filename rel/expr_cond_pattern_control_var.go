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

// Control returns the cond's control expr.
func (e CondPatternControlVarExpr) Control() Expr {
	return e.controlVarExpr
}

// Conditions returns the cond's conditions.
func (e CondPatternControlVarExpr) Conditions() []PatternExprPair {
	return e.conditionPairs
}

func (e CondPatternControlVarExpr) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "cond %v ", e.controlVarExpr.String())

	b.WriteByte('{')
	for i, conditionPair := range e.conditionPairs {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%v", conditionPair.String())
	}

	b.WriteByte('}')
	return b.String()
}

// Eval evaluates to find the first valid condition and return its value.
func (e CondPatternControlVarExpr) Eval(ctx context.Context, scope Scope) (Value, error) {
	varVal, err := e.controlVarExpr.Eval(ctx, scope)
	if err != nil {
		return nil, WrapContextErr(err, e.controlVarExpr, scope)
	}

	for _, conditionPair := range e.conditionPairs {
		ctx, bindings, err := conditionPair.Bind(ctx, scope, varVal)
		if err == nil {
			l, err := scope.MatchedUpdate(bindings)
			if err != nil {
				return nil, WrapContextErr(err, e.controlVarExpr, scope)
			}
			val, err := conditionPair.eval(ctx, l)
			if err != nil {
				return nil, WrapContextErr(err, e.controlVarExpr, l)
			}
			return val, nil
		}
	}

	return None, nil
}
