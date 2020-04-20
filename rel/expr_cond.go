package rel

import (
	"bytes"
	"fmt"
)

// CondExpr returns the tuple applied to a function, the expression looks like:
// cond (1 > 0: 2 + 1, 3 > 4: 5, *: 10).
// The original keyword was switch (SwitchExpr), and finally it was chanaged from switch to cond.
// The keyword cond is more unusual, therefore less likely to be chosen as a regular name,
// and avoids accidental comparisons with the switch statement in other languages.
type CondExpr struct {
	dicExpr, defaultExpr Expr
}

// NewCondExpr returns a new CondExpr.
func NewCondExpr(dict Expr, defaultExpr Expr) Expr {
	return &CondExpr{dict, defaultExpr}
}

// String returns a string representation of the expression.
func (e *CondExpr) String() string {
	var b bytes.Buffer
	b.WriteByte('{')
	var i int = -1
	var expr DictEntryTupleExpr
	for i, expr = range e.dicExpr.(DictExpr).entryExprs {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%v: %v", expr.at.String(), expr.value.String())
	}
	if e.defaultExpr != nil {
		if i == 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%v: %v", "*", e.defaultExpr.String())
	}
	b.WriteByte('}')
	return b.String()
}

// Eval returns the true condition. It must have only one true condition.
func (e *CondExpr) Eval(local Scope) (Value, error) {
	var trueCond *DictEntryTupleExpr
	// Evaluates the valid condition, only one condition whose isTrue() == true can be valid.
	// If there is not any valid condition, the condtion defaultExpr will work.
	for _, expr := range e.dicExpr.(DictExpr).entryExprs {
		tempExpr := expr
		cond, err := tempExpr.at.Eval(local)
		if err != nil {
			return nil, err
		}

		if cond.IsTrue() {
			if trueCond == nil {
				trueCond = &tempExpr
			} else {
				panic("it expects only one true condition in addition to statement '*':valueExpression, " +
					"but actually there are more than 1 true conditions in cond expression")
			}
		}
	}

	var finalCond Expr
	if trueCond != nil {
		finalCond = trueCond.value
	} else if trueCond == nil && e.defaultExpr != nil {
		finalCond = e.defaultExpr
	} else {
		// trueCond == nil && e.defaultCond == nil
		panic("it expects only one true condition or default condition '*':valueExpression, " +
			"but actually there is not any true condition or default condition '*':valueExpression in cond expression")
	}

	value, err := finalCond.Eval(local)
	if err != nil {
		return nil, err
	}
	return value, nil
}
