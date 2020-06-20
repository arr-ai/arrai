package rel

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/anz-bank/pkg/log"
	"github.com/arr-ai/wbnf/parser"
	"github.com/pkg/errors"
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
func (expr CondPatternControlVarExpr) Eval(scope Scope) (Value, error) {
	varVal, err := expr.controlVarExpr.Eval(scope)
	if err != nil {
		return nil, WrapContext(err, expr.controlVarExpr, scope)
	}

	// cond is order dependent. It tries patterns from left to right and stops at the first match.
	// If the first pattern is ident, it can't reach the subsequent patterns.
	// So print out an error message to remind user if the order of cond patterns is correct.
	if len(expr.conditionPairs) > 1 {
		if ident, isIdent := expr.conditionPairs[0].expr.(IdentExpr); isIdent {
			loggingOnce.Do(func() {
				log.Error(context.Background(),
					errors.Errorf("the first cond pattern is %s, it is causing subsequent patterns are unreachable, "+
						"please make sure if the order of pattern is correct.", ident))
			})
		}
	}

	for _, conditionPair := range expr.conditionPairs {
		bindings, err := conditionPair.Bind(scope, varVal)
		if err == nil {
			l, err := scope.MatchedUpdate(bindings)
			if err != nil {
				return nil, WrapContext(err, expr.controlVarExpr, scope)
			}
			val, err := conditionPair.Eval(l)
			if err != nil {
				return nil, WrapContext(err, expr.controlVarExpr, l)
			}
			return val, nil
		}
	}

	return None, nil
}

var loggingOnce sync.Once
