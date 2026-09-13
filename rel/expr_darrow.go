package rel

import (
	"context"
	"errors"
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	goerrors "github.com/go-errors/errors"
)

// DArrowExpr returns the set applied elementwise to a function.
type DArrowExpr struct {
	ExprScanner
	lhs Expr
	fn  *Function
}

// NewDArrowExpr returns a new DArrowExpr.
func NewDArrowExpr(scanner parser.Scanner, lhs Expr, fn Expr) Expr {
	e := &DArrowExpr{ExprScanner{scanner}, lhs, ExprAsFunction(fn)}
	if fastPaths {
		if p := pruneStackedProjects(e); p != nil {
			return p
		}
	}
	return e
}

// pruneStackedProjects drops inner => attributes that the outer identDots
// body does not read (🎯T21 projection pruning).
func pruneStackedProjects(outer *DArrowExpr) Expr {
	inner, ok := outer.lhs.(*DArrowExpr)
	if !ok {
		return nil
	}
	oident, ok := outer.fn.arg.(IdentPattern)
	if !ok {
		return nil
	}
	iident, ok := inner.fn.arg.(IdentPattern)
	if !ok {
		return nil
	}
	ote, ok := outer.fn.body.(*TupleExpr)
	if !ok {
		return nil
	}
	ite, ok := inner.fn.body.(*TupleExpr)
	if !ok {
		return nil
	}
	_, osrc, ok := ote.identDots(string(oident))
	if !ok {
		return nil
	}
	idst, _, ok := ite.identDots(string(iident))
	if !ok {
		return nil
	}
	needed := map[string]bool{}
	for _, s := range osrc {
		needed[s] = true
	}
	keep := make([]AttrExpr, 0, len(ite.attrs))
	for i, name := range idst {
		if needed[name] {
			keep = append(keep, ite.attrs[i])
		}
	}
	if len(keep) == 0 || len(keep) == len(ite.attrs) {
		return nil
	}
	innerFn := NewFunction(inner.fn.Src, inner.fn.arg, NewTupleExpr(ite.Src, keep...))
	newInner := NewDArrowExpr(inner.Src, inner.lhs, innerFn)
	return NewDArrowExpr(outer.Src, newInner, outer.fn)
}

func projectIdentDots(r Relation, dst, src []string) (Set, bool) {
	if isCanonicalTupleShape(dst) {
		return r.projectCanonical(dst, src)
	}
	return r.projectDots(dst, src)
}

// tuplePatternDots recognises `\(:a, :b, ...) (x: a, y: b)` as project/rename
// (🎯T29.4). Leftover `...` may be present only if unused in the body.
func tuplePatternDots(f *Function) (dst, src []string, exact bool, names []string, ok bool) {
	tp, is := f.arg.(TuplePattern)
	if !is {
		return nil, nil, false, nil, false
	}
	te, is := f.body.(*TupleExpr)
	if !is || len(te.attrs) == 0 {
		return nil, nil, false, nil, false
	}
	bound := map[string]string{}
	exact = true
	for _, attr := range tp.attrs {
		if attr.pattern.fallback != nil {
			return nil, nil, false, nil, false
		}
		switch p := attr.pattern.pattern.(type) {
		case ExtraElementPattern:
			if p.ident != "" {
				return nil, nil, false, nil, false
			}
			exact = false
		case IdentPattern:
			bound[string(p)] = attr.name
			names = append(names, attr.name)
		default:
			return nil, nil, false, nil, false
		}
	}
	if len(bound) == 0 {
		return nil, nil, false, nil, false
	}
	dst = make([]string, 0, len(te.attrs))
	src = make([]string, 0, len(te.attrs))
	for _, attr := range te.attrs {
		if attr.IsWildcard() {
			return nil, nil, false, nil, false
		}
		id, is := attr.expr.(IdentExpr)
		if !is {
			return nil, nil, false, nil, false
		}
		from, has := bound[id.ident]
		if !has {
			return nil, nil, false, nil, false
		}
		dst = append(dst, attr.name)
		src = append(src, from)
	}
	return dst, src, exact, names, true
}

// matchColumnExtract reports `ident.attr` as the whole => body (🎯T29.2).
func matchColumnExtract(f *Function, ident string) (string, bool) {
	d, ok := f.body.(*DotExpr)
	if !ok {
		return "", false
	}
	id, ok := d.lhs.(IdentExpr)
	if !ok || id.ident != ident {
		return "", false
	}
	return d.attr, true
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
		if fastPaths {
			if r, is := set.(Relation); is {
				if v, ok, err := e.evalRelationFast(ctx, r, ident, isIdent, local); ok || err != nil {
					return v, err
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
	return nil, WrapContextErr(goerrors.Errorf(
		"=> lhs must be set, not %s: %v", ValueTypeAsString(value), value), e, local)
}

func (e *DArrowExpr) evalRelationFast(
	ctx context.Context, r Relation, ident IdentPattern, isIdent bool, local Scope,
) (Value, bool, error) {
	if isIdent {
		if te, ok := e.fn.body.(*TupleExpr); ok {
			if dst, src, ok := te.identDots(string(ident)); ok {
				if v, ok := projectIdentDots(r, dst, src); ok {
					return v, true, nil
				}
			}
		}
		if attr, ok := matchColumnExtract(e.fn, string(ident)); ok {
			if v, ok := r.projectColumn(attr); ok {
				return v, true, nil
			}
		}
		identStr := string(ident)
		if bin, ok := e.fn.body.(*BinExpr); ok && bin.op == "+>" {
			if id, ok := bin.a.(IdentExpr); ok && id.ident == identStr {
				if te, ok := bin.b.(*TupleExpr); ok && !usesIdentAsValue(te, identStr) {
					if v, ok, err := r.mapAddArrow(ctx, local, te, identStr); ok || err != nil {
						return v, ok, err
					}
				}
			}
		}
		if e.fn.isColumnOnly() {
			if v, ok, err := r.mapAtRow(ctx, local, e.fn, identStr); ok || err != nil {
				return v, ok, err
			}
		}
	}
	if dst, src, exact, names, ok := tuplePatternDots(e.fn); ok {
		if !exact || r.hasOnlyAttrs(names) {
			if v, ok := projectIdentDots(r, dst, src); ok {
				return v, true, nil
			}
		}
	}
	return nil, false, nil
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
	if r, ok := set.(Relation); ok {
		return e.evalParallelRelation(ctx, r, ident, local, ranges)
	}
	elems := make([]Value, 0, set.Count())
	for elem := range All(set) {
		elems = append(elems, elem)
	}
	return e.finishParallel(ctx, local, ident, ranges, len(elems), func(i int) Value {
		return elems[i]
	})
}

func (e *DArrowExpr) evalParallelRelation(
	ctx context.Context, r Relation, ident string, local Scope, ranges [][2]int,
) (Value, bool, error) {
	if e.fn.isColumnOnly() {
		n := r.rows.n
		out := make([]Value, n)
		errs := make([]error, len(ranges))
		runRanges(ranges, func(w, lo, hi int) {
			for i := lo; i < hi; i++ {
				v, err := evalAtRow(ctx, e.fn.body, ident, r.rows.rowAt(i), r.attrMap, local)
				if err != nil {
					errs[w] = err
					return
				}
				out[i] = v
			}
		})
		if err := firstErr(errs); err != nil {
			if errors.Is(err, errNeedRow) {
				return e.finishParallel(ctx, local, ident, ranges, n, func(i int) Value {
					return r.tuple(r.rows.rowAt(i))
				})
			}
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
	return e.finishParallel(ctx, local, ident, ranges, r.rows.n, func(i int) Value {
		return r.tuple(r.rows.rowAt(i))
	})
}

func (e *DArrowExpr) finishParallel(
	ctx context.Context, local Scope, ident string, ranges [][2]int, n int, elem func(int) Value,
) (Value, bool, error) {
	out := make([]Value, n)
	errs := make([]error, len(ranges))
	runRanges(ranges, func(w, lo, hi int) {
		for i := lo; i < hi; i++ {
			v, err := e.fn.body.Eval(ctx, local.With(ident, elem(i)))
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
