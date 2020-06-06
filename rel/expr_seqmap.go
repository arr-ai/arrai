package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// SeqArrowExpr returns the tuple applied to a function.
type SeqArrowExpr struct {
	ExprScanner
	lhs Expr
	fn  *Function
}

// NewSequenceMapExpr returns a new SequenceMapExpr.
func NewSeqArrowExpr(scanner parser.Scanner, lhs Expr, fn Expr) Expr {
	return &SeqArrowExpr{ExprScanner{scanner}, lhs, ExprAsFunction(fn)}
}

// String returns a string representation of the expression.
func (e *SeqArrowExpr) String() string {
	return fmt.Sprintf("(%s >> %s)", e.lhs, e.fn)
}

// Eval returns the lhs
func (e *SeqArrowExpr) Eval(local Scope) (_ Value, err error) {
	defer wrapPanic(e, &err, local)
	value, err := e.lhs.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e, local)
	}
	closure := NewClosure(local, e.fn)

	// TODO: implement directly for String, Array and Dict.
	switch value := value.(type) {
	// case Array:
	// 	values := value.clone().values
	// 	for i, item := range values {
	// 		values[i] =
	// 	}
	case Set:
		values := []Value{}
		for i := value.Enumerator(); i.MoveNext(); {
			t := i.Current().(Tuple)
			pos, _ := t.Get("@")
			attr := t.Names().Without("@").Any()
			item, _ := t.Get(attr)
			newItem := SetCall(closure, item)
			values = append(values, NewTuple(Attr{"@", pos}, Attr{attr, newItem}))
		}
		return NewSet(values...), nil
	}
	return nil, wrapContext(errors.Errorf(">> not applicable to %T", value), e, local)
}
