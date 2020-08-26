package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
)

// StringCharTupleExpr represents an expr that evaluates to a StringCharTuple.
type StringCharTupleExpr struct {
	ExprScanner
	at, char Expr
}

// NewStringCharTupleExpr returns a new dictEntryTupleExpr.
func NewStringCharTupleExpr(scanner parser.Scanner, at, char Expr) StringCharTupleExpr {
	// TODO: Optimise for literals.
	// if at, ok := at.(Value); ok {
	// 	if char, ok := char.(Value); ok {
	// 		return NewDictTuple(at, char)
	// 	}
	// }
	return StringCharTupleExpr{ExprScanner: ExprScanner{Src: scanner}, at: at, char: char}
}

// String returns a string representation of the expression.
func (e StringCharTupleExpr) String() string {
	return fmt.Sprintf("(@: %v, %s: %v)", e.at, StringCharAttr, e.char)
}

// Eval returns the subject
func (e StringCharTupleExpr) Eval(ctx context.Context, local Scope) (_ Value, err error) {
	at, err := e.at.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	char, err := e.char.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	return NewStringCharTuple(int(at.(Number).Float64()), rune(char.(Number).Float64())), nil
}
