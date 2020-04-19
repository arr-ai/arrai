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
	dicExpr DictExpr
	expr    string
}

// NewCondExpr returns a new CondExpr.
func NewCondExpr(dict DictExpr, expr string) Expr {
	return &CondExpr{dict, expr}
}

// String returns a string representation of the expression.
func (e *CondExpr) String() string {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, expr := range e.dicExpr.entryExprs {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%v: %v", expr.at.String(), expr.value.String())
	}
	b.WriteByte('}')
	return b.String()
}

// Eval returns the true condition. It must have only one true condition.
func (e *CondExpr) Eval(local Scope) (Value, error) {
	var trueCond, defaultCond, finalCond *DictEntryTupleExpr
	// Evaluates the valid condition, only one condition whose isTrue() == true can be valid.
	// If there is not any valida condition, the condtion whose String() == '*' will work.
	for _, expr := range e.dicExpr.entryExprs {
		tempExpr := expr
		cond, err := tempExpr.at.Eval(local)
		if err != nil {
			return nil, err
		}
		switch cond.(type) {
		case String:
			// Condition is "*", means matching anything
			defaultCond = &tempExpr
		default:
			if cond.IsTrue() {
				if trueCond == nil {
					trueCond = &tempExpr
				} else {
					panic("it expects only one condition is true, but there are more thant 1 conditions are true in 'cond' expression:" +
						e.expr)
				}
			}
		}
	}

	if trueCond != nil {
		finalCond = trueCond
	} else if trueCond == nil && defaultCond != nil {
		finalCond = defaultCond
	} else {
		// trueCond == nil && defaultCond == nil
		panic("it expects only one condition is true, but there is not any condition is true in 'cond' expression:" +
			e.expr)
	}

	value, err := finalCond.value.Eval(local)
	if err != nil {
		return nil, err
	}
	return value, nil
}
