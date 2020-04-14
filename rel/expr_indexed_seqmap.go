package rel

import (
	"fmt"

	"github.com/go-errors/errors"
)

type IndexedSequenceMapExpr struct {
	lhs Expr
	fn  *Function
}

func NewIndexedSequenceMapExpr(lhs, rhs Expr) Expr {
	return &IndexedSequenceMapExpr{lhs, ExprAsFunction(rhs)}
}

func (is *IndexedSequenceMapExpr) Eval(local Scope) (Value, error) {
	value, err := is.lhs.Eval(local)
	if err != nil {
		return nil, err
	}
	// TODO: implement directly for String, Array and Dict.
	if set, ok := value.(Set); ok {
		values := []Value{}
		for i := set.Enumerator(); i.MoveNext(); {
			indexed, err := is.fn.Body().Eval(local)
			if err != nil {
				return nil, err
			}

			t := i.Current().(Tuple)
			pos, isIndexed := t.Get("@")
			if !isIndexed {
				return nil, errors.Errorf("=> not applicable to unindexed type %v", value)
			}
			attr := t.Names().Without("@").Any()
			item, _ := t.Get(attr)
			itemFn := indexed.(Closure).f
			v, err := itemFn.Body().Eval(local.With(itemFn.Arg(), item).With(is.fn.arg, pos))
			if err != nil {
				return nil, err
			}
			values = append(values, NewTuple(Attr{"@", pos}, Attr{attr, v}))
		}
		return NewSet(values...), nil
	}
	return nil, errors.Errorf("=> not applicable to %T", value)
}

// String returns a string representation of the expression.
func (is *IndexedSequenceMapExpr) String() string {
	return fmt.Sprintf("(%s >>> %s)", is.lhs, is.fn)
}
