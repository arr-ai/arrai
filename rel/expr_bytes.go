package rel

import (
	"context"
	"strings"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// BytesExpr is an expression that evaluates to a Byte Array
type BytesExpr struct {
	ExprScanner
	elements []Expr
}

func NewBytesExpr(scanner parser.Scanner, elements ...Expr) Expr {
	bytes := make([]byte, 0, len(elements))
	for _, expr := range elements {
		if value, is := exprIsValue(expr); is {
			if byteNum, is := value.(Number); is && isByteNumber(byteNum) {
				bytes = append(bytes, byte(int(byteNum)))
				continue
			}
		}
		return BytesExpr{ExprScanner{scanner}, elements}
	}
	return NewBytes(bytes)
}

func (b BytesExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	bytes := make([]byte, 0)
	for _, expr := range b.elements {
		value, err := expr.Eval(ctx, local)
		if err != nil {
			return nil, WrapContextErr(err, b, local)
		}
		switch v := value.(type) {
		case Number:
			if !isByteNumber(v) {
				return nil, WrapContextErr(errors.Errorf("BytesExpr.Eval: Number does not represent a byte: %v", v), b, local)
			}
			bytes = append(bytes, byte(v))
		case String:
			if err := b.handleOffset(v); err != nil {
				return nil, WrapContextErr(err, b, local)
			}
			bytes = append(bytes, []byte(string(v.s))...)
		case GenericSet:
			if s, isString := AsString(v); isString {
				if err := b.handleOffset(s); err != nil {
					return nil, WrapContextErr(err, b, local)
				}
				bytes = append(bytes, []byte(string(s.s))...)
				continue
			}
			return nil, WrapContextErr(errors.Errorf("BytesExpr.Eval: Set %v is not supported", expr), b, local)
		default:
			return nil, WrapContextErr(errors.Errorf("BytesExpr.Eval: %s is not supported", ValueTypeAsString(v)), b, local)
		}
	}
	return NewBytes(bytes), nil
}

func (b BytesExpr) handleOffset(s String) error {
	if s.offset != 0 {
		return errors.Errorf("BytesExpr.Eval: offset string is not supported: %v", s)
	}
	return nil
}

func (b BytesExpr) String() string {
	s := strings.Builder{}
	s.WriteString("<<")
	for i, expr := range b.elements {
		if i > 0 {
			s.WriteString(", ")
		}
		s.WriteString(expr.String())
	}
	s.WriteString(">>")
	return s.String()
}

func isByteNumber(n Number) bool {
	return int(n) >= 0 && int(n) <= 255
}
