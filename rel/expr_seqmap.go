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
	// pure is whether fn's body is free of calls and other observable
	// effects (see hasEffects). Only a pure fn may be suspended in a
	// seqPipeline/dictPipeline: those types let Count and a later Equal or
	// force answer from metadata or run stages out of order, which is
	// invisible for a pure fn but would drop or reorder an effectful one's
	// actions (e.g. //log.print calls whose result is never observed).
	pure bool
}

// NewSequenceMapExpr returns a new SequenceMapExpr.
func NewSeqArrowExpr(withAt bool) func(scanner parser.Scanner, lhs Expr, fn Expr) Expr {
	op := ">>"
	if withAt {
		op = ">>>"
	}
	return func(scanner parser.Scanner, lhs Expr, fn Expr) Expr {
		f := ExprAsFunction(fn)
		return &SeqArrowExpr{
			ExprScanner: ExprScanner{scanner},
			lhs:         lhs,
			fn:          f,
			withAt:      withAt,
			op:          op,
			pure:        !hasEffects(f.body),
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
	case seqPipeline:
		return e.evalSeqPipeline(local, value, call)
	case String: //nolint:dupl
		runes := make([]rune, value.size())
		for at := range runes {
			char := value.runeAt(at)
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
		return e.evalArrayOrSuspend(local, value, call)
	case dictPipeline:
		return e.evalDictPipeline(local, value, call)
	case Dict:
		return e.evalDictOrSuspend(local, value, call)
	case Set:
		b := NewSetBuilder()
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
			b.Add(NewTuple(Attr{"@", at}, Attr{attr, newItem}))
		}
		s, err := b.Finish()
		if err != nil {
			return nil, WrapContextErr(err, e, local)
		}
		return s, nil
	}
	return nil, WrapContextErr(errors.Errorf(
		"%s lhs must be an indexed type, not %s", e.op, ValueTypeAsString(value)), e, local)
}

// evalSeqPipeline continues a suspended >> chain, or forces it first when
// this stage is not pure: see the pure field.
func (e *SeqArrowExpr) evalSeqPipeline(
	local Scope, value seqPipeline, call func(_, _ Value) (Value, error),
) (Value, error) {
	if fastPaths && e.pure {
		return value.then(call), nil
	}
	arr, err := value.force()
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	return e.evalArray(local, arr, call)
}

// evalArrayOrSuspend suspends the map in a seqPipeline when this stage is
// pure, or evaluates it immediately otherwise: see the pure field.
func (e *SeqArrowExpr) evalArrayOrSuspend(
	local Scope, value Array, call func(_, _ Value) (Value, error),
) (Value, error) {
	if fastPaths && e.pure {
		return newSeqPipeline(value, call), nil
	}
	return e.evalArray(local, value, call)
}

// evalDictPipeline continues a suspended >> chain, or forces it first when
// this stage is not pure: see the pure field.
func (e *SeqArrowExpr) evalDictPipeline(
	local Scope, value dictPipeline, call func(_, _ Value) (Value, error),
) (Value, error) {
	if fastPaths && e.pure {
		return value.then(call), nil
	}
	d, err := value.force()
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	return e.evalDict(local, d, call)
}

// evalDictOrSuspend suspends the map in a dictPipeline when this stage is
// pure, or evaluates it immediately otherwise: see the pure field.
func (e *SeqArrowExpr) evalDictOrSuspend(
	local Scope, value Dict, call func(_, _ Value) (Value, error),
) (Value, error) {
	if fastPaths && e.pure {
		return newDictPipeline(value, call), nil
	}
	return e.evalDict(local, value, call)
}

// evalArray maps an array's elements, in parallel when the array is large:
// call's captures are read-only and each element writes only its own slot.
// The reported error is the first element's in index order, as the
// sequential path reports (no worker is interrupted, so every range below
// the lowest failing one completes).
func (e *SeqArrowExpr) evalArray(
	local Scope, value Array, call func(_, _ Value) (Value, error),
) (Value, error) {
	items := make([]Value, len(value.values))
	if fastPaths {
		if ranges := parallelRanges(len(value.values)); ranges != nil {
			errs := make([]error, len(ranges))
			runRanges(ranges, func(w, lo, hi int) {
				for at := lo; at < hi; at++ {
					if item := value.values[at]; item != nil {
						v, err := call(NewNumber(float64(value.offset+at)), item)
						if err != nil {
							errs[w] = err
							return
						}
						items[at] = v
					}
				}
			})
			if err := firstErr(errs); err != nil {
				return nil, WrapContextErr(err, e, local)
			}
			return NewOffsetArray(value.offset, items...), nil
		}
	}
	for at, item := range value.values {
		if item != nil {
			v, err := call(NewNumber(float64(value.offset+at)), item)
			if err != nil {
				return nil, WrapContextErr(err, e, local)
			}
			items[at] = v
		}
	}
	return NewOffsetArray(value.offset, items...), nil
}

// evalDict maps a dict's entries, in parallel when the dict is large — the
// apps/endpoints maps a model pipeline iterates are dicts, so this is often
// the outermost loop. Enumeration order is deterministic for a given dict,
// so the first-error contract matches the sequential path.
func (e *SeqArrowExpr) evalDict(
	local Scope, value Dict, call func(_, _ Value) (Value, error),
) (Value, error) {
	entries := make([]DictEntryTuple, 0, value.m.Count())
	for i := value.Enumerator(); i.MoveNext(); {
		entries = append(entries, i.Current().(DictEntryTuple))
	}
	if ranges := parallelRanges(len(entries)); fastPaths && ranges != nil {
		errs := make([]error, len(ranges))
		runRanges(ranges, func(w, lo, hi int) {
			for j := lo; j < hi; j++ {
				newValue, err := call(entries[j].at, entries[j].value)
				if err != nil {
					errs[w] = err
					return
				}
				entries[j] = NewDictEntryTuple(entries[j].at, newValue)
			}
		})
		if err := firstErr(errs); err != nil {
			return nil, WrapContextErr(err, e, local)
		}
	} else {
		for j, entry := range entries {
			newValue, err := call(entry.at, entry.value)
			if err != nil {
				return nil, WrapContextErr(err, e, local)
			}
			entries[j] = NewDictEntryTuple(entry.at, newValue)
		}
	}
	d, err := NewDict(true, entries...)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	return d, nil
}
