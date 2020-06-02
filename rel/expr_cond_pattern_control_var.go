package rel

import (
	"bytes"
	"fmt"

	"github.com/go-errors/errors"

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
func (expr CondPatternControlVarExpr) Eval(scope Scope) (Value, error) {
	varVal, err := expr.controlVarExpr.Eval(scope)
	if err != nil {
		return nil, wrapContext(err, expr.controlVarExpr, scope)
	}

	for _, conditionPair := range expr.conditionPairs {
		bindings, err := expr.binding(conditionPair, scope, varVal)
		if err == nil {
			l := scope.MatchedUpdate(bindings)
			val, err := conditionPair.Eval(l)
			if err != nil {
				return nil, wrapContext(err, expr.controlVarExpr, l)
			}
			return val, nil
		}
	}

	return None, nil
}

// binding will recover and return EmptyScope and Error whose error message is raised by panic.
// Added this method as some `Bind` will execute panic.
func (expr CondPatternControlVarExpr) binding(conditionPair PatternExpr,
	scope Scope, controlVarVal Value) (rtScope Scope, err error) {
	defer func() {
		if r := recover(); r != nil {
			rtScope = EmptyScope
			switch r := r.(type) {
			case string:
				err = errors.Errorf(r)
			default:
				err = errors.Errorf("calling to method Bind was failed")
			}
		}
	}()

	return conditionPair.Bind(scope, controlVarVal)
}
