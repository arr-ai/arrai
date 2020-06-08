package rel

import (
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
func (e *SeqArrowExpr) Eval(local Scope) (_ Value, err error) {
	defer wrapPanic(e, &err, local)
	value, err := e.lhs.Eval(local)
	if err != nil {
		return nil, wrapContext(err, e, local)
	}
	var closure Set = NewClosure(local, e.fn)
	var call func(at, v Value) Value
	if e.withAt {
		call = func(at, v Value) Value {
			return SetCall(SetCall(closure, at).(Set), v)
		}
	} else {
		call = func(_, v Value) Value {
			return SetCall(closure, v)
		}
	}

	switch value := value.(type) {
	case String:
		runes := make([]rune, len(value.s))
		for at, char := range value.s {
			newChar := call(NewNumber(float64(value.offset+at)), NewNumber(float64(char)))
			if n, is := newChar.(Number); is {
				if r := rune(n.Float64()); float64(r) == n.Float64() {
					runes[at] = r
					continue
				}
			}
			return nil, wrapContext(fmt.Errorf("string %s ... must produce valid chars", e.op), e, local)
		}
		return NewOffsetString(runes, value.offset), nil
	case Bytes:
		bytes := make([]byte, len(value.b))
		for at, byt := range value.b {
			newByte := call(NewNumber(float64(value.offset+at)), NewNumber(float64(byt)))
			if n, is := newByte.(Number); is {
				if b := byte(n.Float64()); float64(b) == n.Float64() {
					bytes[at] = b
					continue
				}
			}
			return nil, wrapContext(fmt.Errorf("bytes %s ... must produce valid bytes", e.op), e, local)
		}
		return NewOffsetBytes(bytes, value.offset), nil
	case Array:
		items := make([]Value, len(value.values))
		for at, item := range value.values {
			if item != nil {
				items[at] = call(NewNumber(float64(value.offset+at)), item)
			}
		}
		return NewOffsetArray(value.offset, items...), nil
	case Dict:
		entries := make([]DictEntryTuple, 0, value.m.Count())
		for i := value.Enumerator(); i.MoveNext(); {
			entry := i.Current().(DictEntryTuple)
			newValue := call(entry.at, entry.value)
			entries = append(entries, NewDictEntryTuple(entry.at, newValue))
		}
		return NewDict(true, entries...), nil
	case Set:
		values := []Value{}
		for i := value.Enumerator(); i.MoveNext(); {
			t := i.Current().(Tuple)
			at, has := t.Get("@")
			if !has {
				return nil, wrapContext(errors.Errorf("%s not applicable to unindexed type %v", e.op, value), e, local)
			}
			attr := t.Names().Without("@").Any()
			item, _ := t.Get(attr)
			newItem := call(at, item)
			values = append(values, NewTuple(Attr{"@", at}, Attr{attr, newItem}))
		}
		return NewSet(values...), nil
	}
	return nil, wrapContext(errors.Errorf("%s not applicable to %T", e.op, value), e, local)
}
