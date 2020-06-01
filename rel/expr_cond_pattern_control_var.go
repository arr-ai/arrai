package rel

import (
	"bytes"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// CondPatternControlVarExpr which is used for `cond` pattern matching.
type CondPatternControlVarExpr struct {
	ExprScanner
	controlVarExpr Expr
	conditionPairs []PatternExpr
}

// NewCondPatternControlVarExpr returns a new CondPatternControlVarExpr.
func NewCondPatternControlVarExpr(scanner parser.Scanner, controlVar Expr, patternExprs ...PatternExpr) Expr {
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
func (expr CondPatternControlVarExpr) Eval(local Scope) (Value, error) {
	varVal, err := expr.controlVarExpr.Eval(local)
	if err != nil {
		return nil, wrapContext(err, expr.controlVarExpr, local)
	}

	for _, conditionPair := range expr.conditionPairs {
		bindings, err := conditionPair.Bind(local, varVal)
		if err == nil {
			l := local.MatchedUpdate(bindings)
			val, err := conditionPair.Eval(l)
			if err != nil {
				return nil, wrapContext(err, expr.controlVarExpr, l)
			}
			return val, nil
		}
	}

	return None, nil
}
