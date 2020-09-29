package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// OffsetExpr is an expression which offsets the provided array by the
// provided offset
type OffsetExpr struct {
	ExprScanner
	offset, array Expr
}

// NewOffsetExpr returns a new OffsetExpr
func NewOffsetExpr(scanner parser.Scanner, n, s Expr) Expr {
	return &OffsetExpr{ExprScanner{scanner}, n, s}
}

func (o *OffsetExpr) Eval(ctx context.Context, local Scope) (_ Value, err error) {
	offset, err := o.offset.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, o, local)
	}
	_, isNumber := offset.(Number)
	if !isNumber {
		return nil, WrapContextErr(errors.Errorf("offset must be a number, not %s", ValueTypeAsString(offset)), o, local)
	}

	array, err := o.array.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, o, local)
	}
	switch a := array.(type) {
	case Array:
		return NewOffsetArray(a.offset+int(offset.(Number)), a.values...), nil
	case Bytes:
		return NewOffsetBytes(a.Bytes(), a.offset+int(offset.(Number))), nil
	case String:
		return NewOffsetString(a.s, a.offset+int(offset.(Number))), nil
	}
	return nil, WrapContextErr(errors.Errorf("offset not applicable to %s", ValueTypeAsString(array)), o, local)
}

func (o *OffsetExpr) String() string {
	return fmt.Sprintf("(%s <: %s)", o.offset, o.array)
}
