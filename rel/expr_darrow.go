package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// DArrowExpr returns the set applied elementwise to a function.
type DArrowExpr struct {
	ExprScanner
	lhs Expr
	fn  *Function
}

// NewDArrowExpr returns a new DArrowExpr.
func NewDArrowExpr(scanner parser.Scanner, lhs Expr, fn Expr) Expr {
	return &DArrowExpr{ExprScanner{scanner}, lhs, ExprAsFunction(fn)}
}

// String returns a string representation of the expression.
func (e *DArrowExpr) String() string {
	return fmt.Sprintf("(%s => %s)", e.lhs, e.fn)
}

// Eval returns the lhs transformed elementwise by fn.
func (e *DArrowExpr) Eval(ctx context.Context, local Scope) (_ Value, err error) {
	value, err := e.lhs.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	if set, ok := value.(Set); ok {
		ident, isIdent := e.fn.arg.(IdentPattern)
		if fastPaths && isIdent {
			if te, ok := e.fn.body.(*TupleExpr); ok {
				if dst, src, ok := te.identDots(string(ident)); ok {
					if r, is := set.(Relation); is {
						if v, ok := r.projectDots(dst, src); ok {
							return v, nil
						}
					}
				}
			}
		}
		if fastPaths && isIdent {
			// The ident path threads nothing between elements, so a large
			// set can evaluate its bodies in parallel. Other patterns
			// thread ctx through Bind and stay sequential.
			if v, done, err := e.evalParallel(ctx, set, string(ident), local); done || err != nil {
				return v, err
			}
		}
		// NOTE: not converted to range-over-All: this body assigns captured
		// locals (ctx, err), which range-over-func turns into heap cells per
		// call — measured as a regression on small-set-heavy workloads.
		b := NewSetBuilder()
		for i := set.Enumerator(); i.MoveNext(); {
			var v Value
			var err error
			if isIdent {
				// Fast path for `set => \x body`: see Closure.call.
				v, err = e.fn.body.Eval(ctx, local.With(string(ident), i.Current()))
			} else {
				var b scopeBuilder
				ctx, err = e.fn.arg.Bind(ctx, local, i.Current(), &b)
				if err != nil {
					if err == errPatternMismatch {
						err = explainBind(ctx, e.fn.arg, local, i.Current())
					}
					return nil, WrapContextErr(err, e, local)
				}
				v, err = e.fn.body.Eval(ctx, local.updateWith(&b))
			}
			if err != nil {
				return nil, WrapContextErr(err, e, local)
			}
			b.Add(v)
		}
		s, err := b.Finish()
		if err != nil {
			return nil, WrapContextErr(err, e, local)
		}
		return s, nil
	}
	return nil, WrapContextErr(errors.Errorf(
		"=> lhs must be set, not %s: %v", ValueTypeAsString(value), value), e, local)
}

// evalParallel evaluates the transform's body over a large set's elements in
// parallel, returning done == false when the set is below the parallel
// threshold. The error, if any, is the first element's in enumeration
// order, matching the sequential path.
func (e *DArrowExpr) evalParallel(
	ctx context.Context, set Set, ident string, local Scope,
) (_ Value, done bool, err error) {
	ranges := parallelRanges(set.Count())
	if ranges == nil {
		return nil, false, nil
	}
	elems := make([]Value, 0, set.Count())
	for elem := range All(set) {
		elems = append(elems, elem)
	}
	out := make([]Value, len(elems))
	errs := make([]error, len(ranges))
	runRanges(ranges, func(w, lo, hi int) {
		for i := lo; i < hi; i++ {
			v, err := e.fn.body.Eval(ctx, local.With(ident, elems[i]))
			if err != nil {
				errs[w] = err
				return
			}
			out[i] = v
		}
	})
	if err := firstErr(errs); err != nil {
		return nil, true, WrapContextErr(err, e, local)
	}
	b := NewSetBuilder()
	for _, v := range out {
		b.Add(v)
	}
	s, err := b.Finish()
	if err != nil {
		return nil, true, WrapContextErr(err, e, local)
	}
	return s, true, nil
}
