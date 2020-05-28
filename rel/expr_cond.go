package rel

import (
	"bytes"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// CondExpr returns the tuple applied to a function, the expression looks like:
// cond (1 > 0: 2 + 1, 3 > 4: 5, *: 10).
// The original keyword was switch (SwitchExpr), and finally it was chanaged from switch to cond.
// The keyword cond is more unusual, therefore less likely to be chosen as a regular name,
// and avoids accidental comparisons with the switch statement in other languages.
type CondExpr struct {
	ExprScanner
	dicExpr         Expr
	validValidation func(condition Value, local Scope) (bool, error) // Valid condition validation.
}

// NewCondExpr returns a new CondExpr.
func NewCondExpr(scanner parser.Scanner, dict Expr) Expr {
	return CondExpr{ExprScanner{scanner}, dict, func(condition Value, local Scope) (bool, error) {
		return condition.IsTrue(), nil
	}}
}

// String returns a string representation of the expression.
func (e CondExpr) String() string {
	var b bytes.Buffer
	b.WriteByte('{')
	var i int = -1
	var expr DictEntryTupleExpr
	switch c := e.dicExpr.(type) {
	case DictExpr:
		for i, expr = range c.entryExprs {
			if i > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "%v: %v", expr.at.String(), expr.value.String())
		}
	}
	b.WriteByte('}')
	return b.String()
}

// Eval returns the value of true condition, or default condition value.
func (e CondExpr) Eval(local Scope) (Value, error) {
	var trueCond *DictEntryTupleExpr

	switch c := e.dicExpr.(type) {
	case DictExpr:
		for _, expr := range c.entryExprs {
			tempExpr := expr
			if expr.at.String() == "_" {
				trueCond = &tempExpr
				break
			}

			cond, err := tempExpr.at.Eval(local)
			if err != nil {
				return nil, wrapContext(err, e)
			}

			valid, err := e.validValidation(cond, local)
			if err != nil {
				return nil, wrapContext(err, e)
			}
			if cond != nil && valid {
				trueCond = &tempExpr
				break
			}
		}
	}

	var finalCond Expr
	if trueCond != nil {
		finalCond = trueCond.value
	} else {
		return None, nil
	}

	value, err := finalCond.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e)
	}
	return value, nil
}
