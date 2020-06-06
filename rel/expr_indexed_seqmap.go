package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

type IndexedSequenceMapExpr struct {
	ExprScanner
	lhs Expr
	fn  *Function
}

func NewISeqArrowExpr(scanner parser.Scanner, lhs, rhs Expr) Expr {
	return &IndexedSequenceMapExpr{ExprScanner{scanner}, lhs, ExprAsFunction(rhs)}
}

func (is *IndexedSequenceMapExpr) Eval(local Scope) (Value, error) {
	value, err := is.lhs.Eval(local)
	if err != nil {
		return nil, wrapContext(err, is, local)
	}
	// TODO: implement directly for String, Array and Dict.
	if set, ok := value.(Set); ok {
		values := []Value{}
		for i := set.Enumerator(); i.MoveNext(); {
			indexed, err := is.fn.Body().Eval(local)
			if err != nil {
				return nil, wrapContext(err, is, local)
			}

			t := i.Current().(Tuple)
			pos, isIndexed := t.Get("@")
			if !isIndexed {
				return nil, wrapContext(errors.Errorf(">>> not applicable to unindexed type %v", value), is, local)
			}
			attr := t.Names().Without("@").Any()
			item, _ := t.Get(attr)
			itemFn := indexed.(Closure).f
			v, err := itemFn.Body().Eval(local.With(itemFn.Arg(), item).With(is.fn.Arg(), pos))
			if err != nil {
				return nil, wrapContext(err, is, local)
			}
			values = append(values, NewTuple(Attr{"@", pos}, Attr{attr, v}))
		}
		return NewSet(values...), nil
	}
	return nil, wrapContext(errors.Errorf(">>> not applicable to %T", value), is, local)
}

// String returns a string representation of the expression.
func (is *IndexedSequenceMapExpr) String() string {
	return fmt.Sprintf("(%s >>> %s)", is.lhs, is.fn)
}
