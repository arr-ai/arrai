package rel

import (
	"fmt"
	"strings"
)

// SliceExpr is an expression that evaluates to a slice of the setToSlice.
type SliceExpr struct {
	setToSlice Expr
	dataRange  *RangeData
}

// NewSliceExpr returns a SliceExpr
func NewSliceExpr(setToSlice Expr, dataRange *RangeData) SliceExpr {
	return SliceExpr{setToSlice, dataRange}
}

// Eval evaluates SliceExpr to the slice of the set.
func (s SliceExpr) Eval(local Scope) (Value, error) {
	data, err := s.dataRange.eval(local)
	if err != nil {
		return nil, err
	}
	start, end, step := data.start, data.end, int(data.step.(Number))

	set, err := s.setToSlice.Eval(local)
	if err != nil {
		return nil, err
	}

	return set.(Set).CallSlice(start, end, step, data.isInclusive()), nil
}

func (s SliceExpr) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("%s(%s)", s.setToSlice, s.dataRange.string()))
	return str.String()
}
