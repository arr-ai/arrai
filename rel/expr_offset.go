package rel

import (
	"fmt"

	"github.com/go-errors/errors"
)

// OffsetExpr is an expression which offsets the provided array by the
// provided offset
type OffsetExpr struct {
	offset, array Expr
}

// NewOffsetExpr returns a new OffsetExpr
func NewOffsetExpr(n, s Expr) Expr {
	return &OffsetExpr{n, s}
}

func (o *OffsetExpr) Eval(local Scope) (Value, error) {
	offset, err := o.offset.Eval(local)
	if err != nil {
		return nil, err
	}
	_, isNumber := offset.(Number)
	if !isNumber {
		return nil, errors.Errorf("\\ not applicable to %T", offset)
	}

	array, err := o.array.Eval(local)
	if err != nil {
		return nil, err
	}
	switch a := array.(type) {
	case Array:
		return NewOffsetArray(a.offset+int(offset.(Number)), a.values...), nil
	case Bytes:
		return NewOffsetBytes(a.Bytes(), a.offset+int(offset.(Number))), nil
	case String:
		return NewOffsetString(a.s, a.offset+int(offset.(Number))), nil
	}
	return nil, errors.Errorf("\\ not applicable to %T", array)
}

func (o *OffsetExpr) String() string {
	return fmt.Sprintf("(%s <: %s)", o.offset, o.array)
}
