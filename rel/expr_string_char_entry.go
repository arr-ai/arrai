package rel

import "fmt"

// StringCharTupleExpr represents a single name:expr in a dictEntryTupleExpr.
type StringCharTupleExpr struct {
	at   Expr
	char Expr
}

// NewStringCharTupleExpr returns a new dictEntryTupleExpr.
func NewStringCharTupleExpr(at, char Expr) StringCharTupleExpr {
	// TODO: Optimise for literals.
	// if at, ok := at.(Value); ok {
	// 	if char, ok := char.(Value); ok {
	// 		return NewDictTuple(at, char)
	// 	}
	// }
	return StringCharTupleExpr{at: at, char: char}
}

// String returns a string representation of the expression.
func (e StringCharTupleExpr) String() string {
	return fmt.Sprintf("(@: %v, %s: %v)", e.at, StringCharAttr, e.char)
}

// Eval returns the subject
func (e StringCharTupleExpr) Eval(local Scope) (Value, error) {
	at, err := e.at.Eval(local)
	if err != nil {
		return nil, err
	}
	char, err := e.char.Eval(local)
	if err != nil {
		return nil, err
	}
	return NewStringCharTuple(int(at.(Number).Float64()), rune(char.(Number).Float64())), nil
}
