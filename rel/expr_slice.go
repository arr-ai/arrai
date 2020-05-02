package rel

import (
	"fmt"
	"strings"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// SliceExpr is an expression that evaluates to a slice of the setToSlice.
type SliceExpr struct {
	ExprScanner
	setToSlice Expr
	dataRange  *RangeData
}

// NewSliceExpr returns a SliceExpr
func NewSliceExpr(scanner parser.Scanner, setToSlice Expr, dataRange *RangeData) SliceExpr {
	return SliceExpr{ExprScanner{scanner}, setToSlice, dataRange}
}

// Eval evaluates SliceExpr to the slice of the set.
func (s SliceExpr) Eval(local Scope) (Value, error) {
	data, err := s.dataRange.eval(local)
	if err != nil {
		return nil, err
	}
	start, end, step := data.start, data.end, int(data.step.(Number))
	if step == 0 {
		return nil, errors.New("step cannot be 0")
	}

	set, err := s.setToSlice.Eval(local)
	if err != nil {
		return nil, err
	}

	slice, err := set.(Set).CallSlice(start, end, step, data.isInclusive())
	if err != nil {
		return nil, wrapContext(err, s)
	}
	return slice, nil
}

func (s SliceExpr) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("%s(%s)", s.setToSlice, s.dataRange.string()))
	return str.String()
}
