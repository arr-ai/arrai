package rel

import (
	"bytes"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// DictExpr represents an expression that yields a dict.
type DictExpr struct {
	ExprScanner
	entryExprs   []DictEntryTupleExpr
	allowDupKeys bool
}

// NewDictExpr returns a new DictExpr from pairs.
func NewDictExpr(
	scanner parser.Scanner,
	allowDupKeys bool, dictExprAlways bool, entryExprs ...DictEntryTupleExpr,
) Expr {
	entries := make([]DictEntryTuple, 0, len(entryExprs))
	for _, expr := range entryExprs {
		if !dictExprAlways {
			if at, is := exprIsValue(expr.at); is {
				if value, is := exprIsValue(expr.value); is {
					entries = append(entries, NewDictEntryTuple(at, value))
					continue
				}
			}
		}
		return DictExpr{ExprScanner: ExprScanner{Src: scanner}, entryExprs: entryExprs, allowDupKeys: allowDupKeys}
	}
	return NewLiteralExpr(scanner, NewDict(allowDupKeys, entries...))
}

// String returns a string representation of the expression.
func (e DictExpr) String() string {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, expr := range e.entryExprs {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%v: %v", expr.at.String(), expr.value.String())
	}
	b.WriteByte('}')
	return b.String()
}

// Eval returns the subject
func (e DictExpr) Eval(local Scope) (Value, error) {
	entryExprs := make([]DictEntryTuple, 0, len(e.entryExprs))
	for _, expr := range e.entryExprs {
		at, err := expr.at.Eval(local)
		if err != nil {
			return nil, wrapContext(err, e, local)
		}
		value, err := expr.value.Eval(local)
		if err != nil {
			return nil, wrapContext(err, e, local)
		}
		entryExprs = append(entryExprs, NewDictEntryTuple(at, value))
	}
	return NewDict(e.allowDupKeys, entryExprs...), nil
}
