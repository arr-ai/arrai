package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// SeqArrowExpr returns the tuple applied to a function.
type SeqArrowExpr struct {
	ExprScanner
	lhs    Expr
	fn     *Function
	withAt bool
	op     string
}

// NewSequenceMapExpr returns a new SequenceMapExpr.
func NewSeqArrowExpr(withAt bool) func(scanner parser.Scanner, lhs Expr, fn Expr) Expr {
	op := ">>"
	if withAt {
		op = ">>>"
	}
	return func(scanner parser.Scanner, lhs Expr, fn Expr) Expr {
		return &SeqArrowExpr{
			ExprScanner: ExprScanner{scanner},
			lhs:         lhs,
			fn:          ExprAsFunction(fn),
			withAt:      withAt,
			op:          op,
		}
	}
}

// String returns a string representation of the expression.
func (e *SeqArrowExpr) String() string {
	return fmt.Sprintf("(%s >> %s)", e.lhs, e.fn)
}

// Eval returns the lhs
func (e *SeqArrowExpr) Eval(ctx context.Context, local Scope) (_ Value, err error) {
	value, err := e.lhs.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	var closure Set = NewClosure(local, e.fn)
	var call func(at, v Value) (Value, error)
	if e.withAt {
		call = func(at, v Value) (Value, error) {
			s, err := SetCall(ctx, closure, at)
			if err != nil {
				return nil, WrapContextErr(err, e, local)
			}
			return SetCall(ctx, s.(Set), v)
		}
	} else {
		call = func(_, v Value) (Value, error) {
			return SetCall(ctx, closure, v)
		}
	}

	switch value := value.(type) {
	case String: //nolint:dupl
		runes := make([]rune, len(value.s))
		for at, char := range value.s {
			newChar, err := call(NewNumber(float64(value.offset+at)), NewNumber(float64(char)))
			if err != nil {
				return nil, WrapContextErr(err, e, local)
			}
			if n, is := newChar.(Number); is {
				if r := rune(n.Float64()); float64(r) == n.Float64() {
					runes[at] = r
					continue
				}
			}
			return nil, WrapContextErr(fmt.Errorf("string %s ... must produce valid chars", e.op), e, local)
		}
		return NewOffsetString(runes, value.offset), nil
	case Bytes: //nolint:dupl
		bytes := make([]byte, len(value.b))
		for at, byt := range value.b {
			newByte, err := call(NewNumber(float64(value.offset+at)), NewNumber(float64(byt)))
			if err != nil {
				return nil, WrapContextErr(err, e, local)
			}
			if n, is := newByte.(Number); is {
				if b := byte(n.Float64()); float64(b) == n.Float64() {
					bytes[at] = b
					continue
				}
			}
			return nil, WrapContextErr(fmt.Errorf("bytes %s ... must produce valid bytes", e.op), e, local)
		}
		return NewOffsetBytes(bytes, value.offset), nil
	case Array:
		items := make([]Value, len(value.values))
		for at, item := range value.values {
			if item != nil {
				items[at], err = call(NewNumber(float64(value.offset+at)), item)
				if err != nil {
					return nil, WrapContextErr(err, e, local)
				}
			}
		}
		return NewOffsetArray(value.offset, items...), nil
	case Dict:
		entries := make([]DictEntryTuple, 0, value.m.Count())
		for i := value.Enumerator(); i.MoveNext(); {
			entry := i.Current().(DictEntryTuple)
			newValue, err := call(entry.at, entry.value)
			if err != nil {
				return nil, WrapContextErr(err, e, local)
			}
			entries = append(entries, NewDictEntryTuple(entry.at, newValue))
		}
		d, err := NewDict(true, entries...)
		if err != nil {
			return nil, WrapContextErr(err, e, local)
		}
		return d, nil
	case Set:
		values := []Value{}
		for i := value.Enumerator(); i.MoveNext(); {
			t, ok := i.Current().(Tuple)
			if !ok {
				return nil, WrapContextErr(errors.Errorf(
					"%s lhs must be an indexed type, not %s", e.op, ValueTypeAsString(value)), e, local)
			}
			at, has := t.Get("@")
			if !has {
				return nil, WrapContextErr(errors.Errorf(
					"%s lhs must be an indexed type, not %s", e.op, ValueTypeAsString(value)), e, local)
			}
			attr := t.Names().Without("@").Any()
			item, _ := t.Get(attr)
			newItem, err := call(at, item)
			if err != nil {
				return nil, WrapContextErr(err, e, local)
			}
			values = append(values, NewTuple(Attr{"@", at}, Attr{attr, newItem}))
		}
		s, err := NewSet(values...)
		if err != nil {
			return nil, WrapContextErr(err, e, local)
		}
		return s, nil
	}
	return nil, WrapContextErr(errors.Errorf(
		"%s lhs must be an indexed type, not %s", e.op, ValueTypeAsString(value)), e, local)
}
