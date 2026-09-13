package rel

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"

	"github.com/arr-ai/arrai/pkg/deprecate"
)

var (
	plusDeprecation = deprecate.MustNewDeprecator(
		"use of + for concatenation",
		"2021-03-03", "2021-04-03", "2021-06-03",
	)
)

type binEval func(ctx context.Context, a, b Value, local Scope) (Value, error)

// BinExpr represents a range of operators.
type BinExpr struct {
	ExprScanner
	a, b   Expr
	op     string
	format string
	eval   binEval
}

func newBinExpr(scanner parser.Scanner, a, b Expr, op, format string, eval binEval) Expr {
	return &BinExpr{ExprScanner{scanner}, a, b, op, format, eval}
}

type valueEval func(a, b Value) Value

// MakeBinValExpr returns a function that creates a binExpr for the given
// logical operator.
func MakeBinValExpr(op string, eval valueEval) func(scanner parser.Scanner, a, b Expr) Expr {
	return func(scanner parser.Scanner, a, b Expr) Expr {
		return newBinExpr(scanner, a, b, op, "(%s "+op+" %s)",
			func(ctx context.Context, a, b Value, _ Scope) (Value, error) {
				return eval(a, b), nil
			})
	}
}

type arithEval func(a, b float64) float64

func newArithExpr(scanner parser.Scanner, a, b Expr, op string, eval arithEval) Expr {
	return newBinExpr(scanner, a, b, op, "(%s "+op+" %s)",
		func(_ context.Context, a, b Value, _ Scope) (Value, error) {
			if a, ok := a.(Number); ok {
				if b, ok := b.(Number); ok {
					return NewNumber(eval(a.Float64(), b.Float64())), nil
				}
			}
			return nil, errors.Errorf(
				"Both args to %q must be numbers, not %s and %s",
				op, ValueTypeAsString(a), ValueTypeAsString(b))
		})
}

func addValues(ctx context.Context, scanner parser.Scanner, a, b Value) (Value, error) {
	if a, ok := a.(Number); ok {
		if b, ok := b.(Number); ok {
			return NewNumber(a.Float64() + b.Float64()), nil
		}
	}
	if a, ok := a.(Tuple); ok {
		if b, ok := b.(Tuple); ok {
			if err := plusDeprecation.Deprecate(ctx, scanner); err != nil {
				return nil, err
			}
			return MergeLeftToRight(a, b), nil
		}
	}
	if a, ok := a.(Set); ok {
		if b, ok := b.(Set); ok {
			if err := plusDeprecation.Deprecate(ctx, scanner); err != nil {
				return nil, err
			}
			return Concatenate(a, b)
		}
	}
	return nil, errors.Errorf(
		"Both args to + must be numbers or tuples, not %s and %s",
		ValueTypeAsString(a), ValueTypeAsString(b))
}

// NewAddExpr evaluates a + b, given two Numbers.
func NewAddExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newBinExpr(scanner, a, b, "+", "(%s + %s)",
		func(ctx context.Context, a, b Value, _ Scope) (Value, error) {
			return addValues(ctx, scanner, a, b)
		})
}

// NewAddArrowExpr returns a new BinExpr which supports operator `+>`.
func NewAddArrowExpr(scanner parser.Scanner, lhs, rhs Expr) Expr {
	return newBinExpr(scanner, lhs, rhs, "+>", "(%s +> %s)",
		func(_ context.Context, lhs, rhs Value, _ Scope) (Value, error) {
			return evalValForAddArrow(lhs, rhs)
		})
}

// NewSubExpr evaluates a - b, given two Numbers.
func NewSubExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newArithExpr(scanner, a, b, "-", func(a, b float64) float64 { return a - b })
}

// NewMulExpr evaluates a * b, given two Numbers.
func NewMulExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newArithExpr(scanner, a, b, "*", func(a, b float64) float64 { return a * b })
}

// NewDivExpr evaluates a / b, given two Numbers.
func NewDivExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newArithExpr(scanner, a, b, "/", func(a, b float64) float64 { return a / b })
}

// NewIdivExpr evaluates ⎣a / b⎦, given two Numbers.
func NewIdivExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newArithExpr(scanner, a, b, "//", func(a, b float64) float64 {
		return math.Floor(a / b)
	})
}

// NewModExpr evaluates a % b, given two Numbers.
func NewModExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newArithExpr(scanner, a, b, "%", func(a, b float64) float64 {
		return math.Mod(a, b)
	})
}

// NewSubModExpr evaluates a % b, given two Numbers.
func NewSubModExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newArithExpr(scanner, a, b, "-%", func(a, b float64) float64 {
		return a - math.Mod(a, b)
	})
}

// NewPowExpr evaluates a to the power of b, given two Numbers.
func NewPowExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newArithExpr(scanner, a, b, "^", func(a, b float64) float64 {
		return math.Pow(a, b)
	})
}

// NewWithExpr evaluates a with b, given a set lhs.
func NewWithExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newBinExpr(scanner, a, b, "with", "(%s with %s)",
		func(_ context.Context, a, b Value, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				return x.With(b), nil
			}
			return nil, errors.Errorf("'with' lhs must be a set, not %s", ValueTypeAsString(a))
		})
}

// NewWithoutExpr evaluates a without b, given a set lhs.
func NewWithoutExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newBinExpr(scanner, a, b, "without", "(%s without %s)",
		func(_ context.Context, a, b Value, _ Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if s := x.Without(b); s.IsTrue() {
					return s, nil
				}
				return None, nil
			}
			return nil, errors.Errorf("'without' lhs must be a set, not %s", ValueTypeAsString(a))
		})
}

// eqAttrPredicate describes a where-predicate of the form `.attr = key`
// (or `key = .attr`) where key does not depend on `.`. Such predicates can be
// answered from a relation's cached attribute index instead of a full scan.
type eqAttrPredicate struct {
	attr string
	key  Expr
}

// matchEqAttrPredicates recognises `\. .attr = key` and And-trees of those
// (🎯T29.1). Any non-eq conjunct refuses the whole match so we never index
// only part of a mixed predicate.
func matchEqAttrPredicates(f *Function) []eqAttrPredicate {
	if ident, is := f.arg.(IdentPattern); !is || ident != "." {
		return nil
	}
	var out []eqAttrPredicate
	if !collectEqAttrPredicates(f.body, &out) || len(out) == 0 {
		return nil
	}
	return out
}

func collectEqAttrPredicates(e Expr, out *[]eqAttrPredicate) bool {
	switch e := e.(type) {
	case AndExpr:
		return collectEqAttrPredicates(e.a, out) && collectEqAttrPredicates(e.b, out)
	case CompareExpr:
		p := matchOneEqAttr(e)
		if p == nil {
			return false
		}
		*out = append(*out, *p)
		return true
	default:
		return false
	}
}

func matchOneEqAttr(cmp CompareExpr) *eqAttrPredicate {
	if len(cmp.args) != 2 || cmp.ops[0] != "=" {
		return nil
	}
	if attr, ok := isDotAttr(cmp.args[0]); ok && isDotFreeKey(cmp.args[1]) {
		return &eqAttrPredicate{attr: attr, key: cmp.args[1]}
	}
	if attr, ok := isDotAttr(cmp.args[1]); ok && isDotFreeKey(cmp.args[0]) {
		return &eqAttrPredicate{attr: attr, key: cmp.args[0]}
	}
	return nil
}

func isDotAttr(e Expr) (string, bool) {
	if d, is := e.(*DotExpr); is {
		if id, is := d.lhs.(IdentExpr); is && id.ident == "." {
			return d.attr, true
		}
	}
	return "", false
}

func isDotFreeKey(e Expr) bool {
	// Conservative: identifiers other than `.`, literals, and values.
	switch e := e.(type) {
	case IdentExpr:
		return e.ident != "."
	case LiteralExpr:
		return true
	case Value:
		return true
	}
	return false
}

type eqAttrBound struct {
	index int
	key   Value
}

func eqAttrFact(bounds []eqAttrBound) string {
	if len(bounds) == 1 {
		// Keep the single-pred fact format so existing plan-cache tests hold.
		return strconv.Itoa(bounds[0].index) + "=" + bounds[0].key.String()
	}
	ord := append([]eqAttrBound(nil), bounds...)
	for i := 1; i < len(ord); i++ {
		for j := i; j > 0 && ord[j].index < ord[j-1].index; j-- {
			ord[j], ord[j-1] = ord[j-1], ord[j]
		}
	}
	var b strings.Builder
	for i, bound := range ord {
		if i > 0 {
			b.WriteByte('&')
		}
		b.WriteString(strconv.Itoa(bound.index))
		b.WriteByte('=')
		b.WriteString(bound.key.String())
	}
	return b.String()
}

// intersectIDs returns the ascending intersection of two ascending id slices.
func intersectIDs(a, b []uint32) []uint32 {
	out := make([]uint32, 0, min(len(a), len(b)))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] == b[j]:
			out = append(out, a[i])
			i++
			j++
		case a[i] < b[j]:
			i++
		default:
			j++
		}
	}
	return out
}

func (r Relation) idsForEqAttr(index int, key Value) ([]uint32, bool) {
	if n, is := key.(Number); is {
		z := r.rows.zone(index)
		if z.ok && (n.Less(z.min) || z.max.Less(n)) {
			return nil, false
		}
	}
	return r.rows.groupBy(valueProjector{index}).getKey(key)
}

// whereByIndex answers `r where .attr = key` (and And of those) from r's
// cached groupBy indexes, intersecting id lists when there are several.
func (r Relation) whereByIndex(ctx context.Context, scope Scope, preds []eqAttrPredicate) (Value, bool, error) {
	if len(preds) == 0 {
		return nil, false, nil
	}
	bounds := make([]eqAttrBound, 0, len(preds))
	for _, p := range preds {
		index, has := r.attrMap[p.attr]
		if !has {
			return nil, false, nil
		}
		key, err := p.key.Eval(ctx, scope)
		if err != nil {
			return nil, false, err
		}
		bounds = append(bounds, eqAttrBound{index: index, key: key})
	}
	n0 := r.rows.n
	// Key by store column, not dest attr: injective projectDots siblings share
	// the positionalRelation, so dest names like k/x collide across projectors.
	fact := eqAttrFact(bounds)
	guard := func() bool { return r.rows.n == n0 }
	if view, hit := r.rows.planGet(fact, guard); hit {
		if view == nil {
			return None, true, nil
		}
		return r.newBody(view), true, nil
	}
	var ids []uint32
	for i, bound := range bounds {
		got, has := r.idsForEqAttr(bound.index, bound.key)
		if !has {
			r.rows.planPut(fact, guard, nil)
			return None, true, nil
		}
		if i == 0 {
			ids = got
			continue
		}
		ids = intersectIDs(ids, got)
		if len(ids) == 0 {
			r.rows.planPut(fact, guard, nil)
			return None, true, nil
		}
	}
	view := r.rows.selView(ids)
	r.rows.planPut(fact, guard, view)
	return r.newBody(view), true, nil
}

// pushWhereThroughProject rewrites `(rel => (a: .a, ...)) where .a = k` to
// `(rel where .a = k) => (a: .a, ...)` when every predicate attr is projected
// unchanged (🎯T21, 🎯T29.1).
func pushWhereThroughProject(scanner parser.Scanner, a, pred Expr, eqs []eqAttrPredicate) Expr {
	if len(eqs) == 0 {
		return nil
	}
	d, ok := a.(*DArrowExpr)
	if !ok {
		return nil
	}
	ident, is := d.fn.arg.(IdentPattern)
	if !is {
		return nil
	}
	te, ok := d.fn.body.(*TupleExpr)
	if !ok {
		return nil
	}
	dst, src, ok := te.identDots(string(ident))
	if !ok {
		return nil
	}
	for _, eq := range eqs {
		found := false
		for i, name := range dst {
			if name == eq.attr && src[i] == eq.attr {
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	inner := NewWhereExpr(scanner, d.lhs, pred)
	return NewDArrowExpr(scanner, inner, d.fn)
}

type colPredKind int

const (
	colTruthy colPredKind = iota
	colFalsy
	colNE
	colIn
	colNotIn
)

// colPred is a column predicate that can be evaluated on the arena (🎯T29.5).
type colPred struct {
	attr string
	kind colPredKind
	key  Expr
}

func matchColPred(f *Function) *colPred {
	if ident, is := f.arg.(IdentPattern); !is || ident != "." {
		return nil
	}
	if attr, ok := isDotAttr(f.body); ok {
		return &colPred{attr: attr, kind: colTruthy}
	}
	if u, ok := f.body.(*UnaryExpr); ok && u.op == "!" {
		if attr, ok := isDotAttr(u.a); ok {
			return &colPred{attr: attr, kind: colFalsy}
		}
	}
	cmp, ok := f.body.(CompareExpr)
	if !ok || len(cmp.args) != 2 {
		return nil
	}
	attr, leftDot := isDotAttr(cmp.args[0])
	other := cmp.args[1]
	if !leftDot {
		attr, leftDot = isDotAttr(cmp.args[1])
		other = cmp.args[0]
		if !leftDot || !isDotFreeKey(other) {
			return nil
		}
		if cmp.ops[0] == "!=" {
			return &colPred{attr: attr, kind: colNE, key: other}
		}
		return nil
	}
	if !isDotFreeKey(other) {
		return nil
	}
	switch cmp.ops[0] {
	case "!=":
		return &colPred{attr: attr, kind: colNE, key: other}
	case "<:":
		return &colPred{attr: attr, kind: colIn, key: other}
	case "!<:":
		return &colPred{attr: attr, kind: colNotIn, key: other}
	}
	return nil
}

func (r Relation) whereByColumn(ctx context.Context, scope Scope, p *colPred) (Value, bool, error) {
	index, has := r.attrMap[p.attr]
	if !has {
		return nil, false, nil
	}
	var key Value
	var set Set
	if p.key != nil {
		var err error
		key, err = p.key.Eval(ctx, scope)
		if err != nil {
			return nil, false, err
		}
		if p.kind == colIn || p.kind == colNotIn {
			var ok bool
			set, ok = key.(Set)
			if !ok {
				return nil, false, fmt.Errorf("%s rhs not a set: %v", colPredOp(p.kind), key)
			}
		}
	}
	ids := make([]uint32, 0, r.rows.n)
	for i := 0; i < r.rows.n; i++ {
		v := r.rows.rowAt(i)[index]
		ok := false
		switch p.kind {
		case colTruthy:
			ok = v.IsTrue()
		case colFalsy:
			ok = !v.IsTrue()
		case colNE:
			ok = !v.Equal(key)
		case colIn:
			ok = set.Has(v)
		case colNotIn:
			ok = !set.Has(v)
		}
		if ok {
			ids = append(ids, r.rows.arenaID(i))
		}
	}
	if len(ids) == 0 {
		return None, true, nil
	}
	return r.newBody(r.rows.selView(ids)), true, nil
}

func colPredOp(k colPredKind) string {
	switch k {
	case colIn:
		return "<:"
	case colNotIn:
		return "!<:"
	default:
		return "!="
	}
}

// NewWhereExpr evaluates a where pred, given a set lhs.
func NewWhereExpr(scanner parser.Scanner, a, pred Expr) Expr {
	predFn := ExprAsFunction(pred)
	pred = predFn
	eqPreds := matchEqAttrPredicates(predFn)
	colPred := matchColPred(predFn)
	if fastPaths {
		if pushed := pushWhereThroughProject(scanner, a, pred, eqPreds); pushed != nil {
			return pushed
		}
	}
	return newBinExpr(scanner, a, pred, "where", "(%s where %s)",
		func(ctx context.Context, a, pred Value, local Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if p, ok := pred.(Closure); ok {
					if r, is := x.(Relation); is && fastPaths && len(eqPreds) > 0 {
						if v, done, err := r.whereByIndex(ctx, p.scope, eqPreds); done || err != nil {
							return v, err
						}
					}
					if r, is := x.(Relation); is && fastPaths && colPred != nil {
						if v, done, err := r.whereByColumn(ctx, p.scope, colPred); done || err != nil {
							return v, err
						}
					}
					if r, is := x.(Relation); is && fastPaths {
						if ident, ok := predFn.arg.(IdentPattern); ok && predFn.isColumnOnly() {
							if v, done, err := r.whereAtRow(ctx, p.scope, predFn, string(ident)); done || err != nil {
								return v, err
							}
						}
					}
					s, err := x.Where(func(v Value) (bool, error) {
						r, err := SetCall(ctx, p, v)
						if err != nil {
							return false, err
						}
						return r.IsTrue(), nil
					})
					if err != nil {
						return nil, err
					}
					if !s.IsTrue() {
						return None, nil
					}
					return s, nil
				}
				return nil, errors.Errorf("'where' rhs must be a function, not %s", ValueTypeAsString(a))
			}
			return nil, errors.Errorf("'where' lhs must be a set, not %s", ValueTypeAsString(a))
		})
}

// NewOrderByExpr evaluates a orderby key, given a set lhs, returning an array.
func NewOrderByExpr(scanner parser.Scanner, a, key Expr) Expr {
	keyFn := ExprAsFunction(key)
	return newBinExpr(scanner, a, keyFn, "orderby", "(%s orderby %s)",
		func(ctx context.Context, a, key Value, local Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if r, is := x.(Relation); is && fastPaths {
					if ident, ok := keyFn.arg.(IdentPattern); ok {
						if attr, ok := matchColumnExtract(keyFn, string(ident)); ok {
							if values, ok := r.orderByColumn(attr); ok {
								return NewArray(values...), nil
							}
						}
					}
				}
				if k, ok := key.(Closure); ok {
					values, err := OrderBy(x,
						func(value Value) (Value, error) {
							return SetCall(ctx, k, value)
						},
						func(a, b Value) bool {
							return a.Less(b)
						})
					if err != nil {
						return nil, err
					}
					return NewArray(values...), nil
				}
				return nil, errors.Errorf("'orderby' rhs must be a function, not %s", ValueTypeAsString(a))
			}
			return nil, errors.Errorf("'orderby' lhs must be a set, not %s", ValueTypeAsString(a))
		})
}

// NewOrderExpr evaluates a order less, given a set lhs, returning an array.
func NewOrderExpr(scanner parser.Scanner, a, key Expr) Expr {
	key = ExprAsFunction(key)
	return newBinExpr(scanner, a, key, "order", "(%s orderby %s)",
		func(ctx context.Context, a, less Value, local Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if l, ok := less.(Closure); ok {
					values, err := OrderBy(x,
						func(value Value) (Value, error) {
							return value, nil
						},
						func(a, b Value) bool {
							c, err := SetCall(ctx, l, a)
							if err != nil {
								panic(err)
							}
							less, err := SetCall(ctx, c.(Closure), b)
							if err != nil {
								panic(err)
							}
							return less.IsTrue()
						})
					if err != nil {
						return nil, err
					}
					return NewArray(values...), nil
				}
				return nil, errors.Errorf("'order' rhs must be a function, not %s", ValueTypeAsString(a))
			}
			return nil, errors.Errorf("'order' lhs must be a set, not %s", ValueTypeAsString(a))
		})
}

// NewRankExpr evaluates a rank tuplef, given a relation lhs, returning a new
// relation with each lhs tuple augmented by the tuplef attrs containing the
// corresponding rank.
func NewRankExpr(scanner parser.Scanner, a, key Expr) Expr {
	key = ExprAsFunction(key)
	return newBinExpr(scanner, a, key, "rank", "(%s rank %s)",
		func(ctx context.Context, a, tuplef Value, local Scope) (Value, error) {
			if x, ok := a.(Set); ok {
				if l, ok := tuplef.(Closure); ok {
					return Rank(x, func(v Tuple) (Tuple, error) {
						result, err := SetCall(ctx, l, v)
						if err != nil {
							return nil, err
						}
						return result.(Tuple), nil
					})
				}
				return nil, errors.Errorf("'rank' rhs must be a function, not %s", ValueTypeAsString(a))
			}
			return nil, errors.Errorf("'rank' lhs must be a set, not %s", ValueTypeAsString(a))
		})
}

func Call(ctx context.Context, a, b Value, _ Scope) (Value, error) {
	if x, ok := a.(Set); ok {
		return SetCall(ctx, x, b)
	}
	return nil, errors.Errorf(
		"call lhs must be a function, not %s", ValueTypeAsString(a))
}

// NewCallExpr evaluates a without b, given a set lhs.
func NewCallExpr(scanner parser.Scanner, a, b Expr) Expr {
	return newBinExpr(scanner, a, b, "call", "«%s»(%s)", Call)
}

func NewCallExprCurry(scanner parser.Scanner, f Expr, args ...Expr) Expr {
	for _, arg := range args {
		f = NewCallExpr(scanner, f, arg)
	}
	return f
}

// String returns a string representation of the expression.
func (e *BinExpr) String() string {
	return fmt.Sprintf(e.format, e.a, e.b)
}

// Eval returns the subject
func (e *BinExpr) Eval(ctx context.Context, local Scope) (_ Value, err error) {
	a, err := e.a.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}

	b, err := e.b.Eval(ctx, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	val, err := e.eval(ctx, a, b, local)
	if err != nil {
		return nil, WrapContextErr(err, e, local)
	}
	return val, nil
}

// evalValForAddArrow evaluates operator `+>`.
func evalValForAddArrow(lhs, rhs Value) (Value, error) {
	switch lhs := lhs.(type) {
	case Tuple:
		if rhs, ok := rhs.(Tuple); ok {
			return MergeLeftToRight(lhs, rhs), nil
		}
	case Dict:
		switch rhs := rhs.(type) {
		case Dict:
			return mergeDicts(lhs, rhs), nil
		case Set:
			if !rhs.IsTrue() {
				return lhs, nil
			}
		}
	case Set:
		if !lhs.IsTrue() {
			switch rhs := rhs.(type) {
			case Dict:
				return rhs, nil
			case Set:
				if !rhs.IsTrue() {
					return lhs, nil
				}
			}
		}
	}

	return nil, errors.Errorf(
		"Args to +> must be both tuples or both dicts, not %s and %s",
		ValueTypeAsString(lhs), ValueTypeAsString(rhs))
}

func mergeDicts(lhs Dict, rhs Dict) Dict {
	tempMap := lhs.m
	for e := rhs.DictEnumerator(); e.MoveNext(); {
		key, value := e.Current()
		tempMap = tempMap.With(key, value)
	}
	return newDict(tempMap)
}
