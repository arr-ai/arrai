package rel

import "fmt"

// dictEntryTupleExpr represents a single name:expr in a dictEntryTupleExpr.
type DictEntryTupleExpr struct {
	at    Expr
	value Expr
}

// NewDictEntryTupleExpr returns a new dictEntryTupleExpr.
func NewDictEntryTupleExpr(at, value Expr) DictEntryTupleExpr {
	// TODO: Optimise for literals.
	// if at, ok := at.(Value); ok {
	// 	if value, ok := value.(Value); ok {
	// 		return NewDictTuple(at, value)
	// 	}
	// }
	return DictEntryTupleExpr{at: at, value: value}
}

// String returns a string representation of the expression.
func (e DictEntryTupleExpr) String() string {
	return fmt.Sprintf("(@: %v, %s: %v)", e.at, DictEntryAttr, e.value)
}

// Eval returns the subject
func (e DictEntryTupleExpr) Eval(local Scope) (Value, error) {
	at, err := e.at.Eval(local)
	if err != nil {
		return nil, err
	}
	value, err := e.value.Eval(local)
	if err != nil {
		return nil, err
	}
	return NewDictTuple(at, value), nil
}
