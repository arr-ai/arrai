package rel

import (
	"context"
	"errors"
)

// errNeedRow means the body used ident as a Value; the caller should inflate.
var errNeedRow = errors.New("row cursor needs the row as a value")

// usesIdentAsValue reports whether ident appears as a Value, not only as the
// lhs of ident.attr (or as the row being extended by +>). 🎯T29.9
func usesIdentAsValue(e Expr, ident string) bool {
	if e == nil {
		return false
	}
	switch e := e.(type) {
	case IdentExpr:
		return e.ident == ident
	case *DotExpr:
		if id, ok := e.lhs.(IdentExpr); ok && id.ident == ident {
			return false
		}
		return usesIdentAsValue(e.lhs, ident)
	case *BinExpr:
		if e.op == "+>" {
			if id, ok := e.a.(IdentExpr); ok && id.ident == ident {
				return usesIdentAsValue(e.b, ident)
			}
		}
		return usesIdentAsValue(e.a, ident) || usesIdentAsValue(e.b, ident)
	case *UnaryExpr:
		return usesIdentAsValue(e.a, ident)
	case AndExpr:
		return usesIdentAsValue(e.a, ident) || usesIdentAsValue(e.b, ident)
	case OrExpr:
		return usesIdentAsValue(e.a, ident) || usesIdentAsValue(e.b, ident)
	case *IfElseExpr:
		return usesIdentAsValue(e.ifTrue, ident) || usesIdentAsValue(e.cond, ident) || usesIdentAsValue(e.ifFalse, ident)
	case CompareExpr:
		for _, a := range e.args {
			if usesIdentAsValue(a, ident) {
				return true
			}
		}
		return false
	case *TupleExpr:
		for _, a := range e.attrs {
			if usesIdentAsValue(a.expr, ident) {
				return true
			}
		}
		return false
	case *Function:
		if shadowsIdent(e.arg, ident) {
			return false
		}
		return usesIdentAsValue(e.body, ident)
	case *DArrowExpr:
		return usesIdentAsValue(e.lhs, ident) || usesIdentAsValue(e.fn, ident)
	case *ArrowExpr:
		return usesIdentAsValue(e.lhs, ident) || usesIdentAsValue(e.fn, ident)
	default:
		return countIdentUses(e, ident) > 0
	}
}

func indexOfName(names NamesSlice, name string) (int, bool) {
	for i, n := range names {
		if n == name {
			return i, true
		}
	}
	return 0, false
}

func evalAtRow(
	ctx context.Context, e Expr, ident string, row Values, cols map[string]int, local Scope,
) (Value, error) {
	if e == nil {
		return nil, errNeedRow
	}
	switch e := e.(type) {
	case Value:
		return e, nil
	case LiteralExpr:
		v, err := e.Eval(ctx, local)
		return v, err
	case IdentExpr:
		if e.ident == ident {
			return nil, errNeedRow
		}
		return e.Eval(ctx, local)
	case *DotExpr:
		if id, ok := e.lhs.(IdentExpr); ok && id.ident == ident {
			i, ok := cols[e.attr]
			if !ok {
				return nil, errNeedRow
			}
			return row[i], nil
		}
		lhs, err := evalAtRow(ctx, e.lhs, ident, row, cols, local)
		if err != nil {
			return nil, err
		}
		t, ok := lhs.(Tuple)
		if !ok {
			return e.Eval(ctx, local)
		}
		v, found := t.Get(e.attr)
		if !found {
			return nil, errNeedRow
		}
		return v, nil
	case *BinExpr:
		a, err := evalAtRow(ctx, e.a, ident, row, cols, local)
		if err != nil {
			return nil, err
		}
		b, err := evalAtRow(ctx, e.b, ident, row, cols, local)
		if err != nil {
			return nil, err
		}
		return e.eval(ctx, a, b, local)
	case *UnaryExpr:
		a, err := evalAtRow(ctx, e.a, ident, row, cols, local)
		if err != nil {
			return nil, err
		}
		return e.eval(ctx, a, local)
	case AndExpr:
		a, err := evalAtRow(ctx, e.a, ident, row, cols, local)
		if err != nil {
			return nil, err
		}
		if !a.IsTrue() {
			return a, nil
		}
		return evalAtRow(ctx, e.b, ident, row, cols, local)
	case OrExpr:
		a, err := evalAtRow(ctx, e.a, ident, row, cols, local)
		if err != nil {
			return nil, err
		}
		if a.IsTrue() {
			return a, nil
		}
		return evalAtRow(ctx, e.b, ident, row, cols, local)
	case CompareExpr:
		lhs, err := evalAtRow(ctx, e.args[0], ident, row, cols, local)
		if err != nil {
			return nil, err
		}
		for i, arg := range e.args[1:] {
			rhs, err := evalAtRow(ctx, arg, ident, row, cols, local)
			if err != nil {
				return nil, err
			}
			sat, err := e.comps[i](lhs, rhs)
			if err != nil {
				return nil, err
			}
			if !sat {
				return False, nil
			}
			lhs = rhs
		}
		return True, nil
	case *IfElseExpr:
		cond, err := evalAtRow(ctx, e.cond, ident, row, cols, local)
		if err != nil {
			return nil, err
		}
		if cond.IsTrue() {
			return evalAtRow(ctx, e.ifTrue, ident, row, cols, local)
		}
		return evalAtRow(ctx, e.ifFalse, ident, row, cols, local)
	case *TupleExpr:
		if e.attrSet.namesRep != nil {
			vals := make([]Value, len(e.slots))
			for i, attr := range e.attrs {
				v, err := evalAtRow(ctx, attr.expr, ident, row, cols, local)
				if err != nil {
					return nil, err
				}
				vals[e.slots[i]] = v
			}
			return newShapedTuple(e.attrSet, vals), nil
		}
		tuple := EmptyTuple
		for _, attr := range e.attrs {
			v, err := evalAtRow(ctx, attr.expr, ident, row, cols, local)
			if err != nil {
				return nil, err
			}
			if attr.IsWildcard() {
				return nil, errNeedRow
			}
			tuple = tuple.With(attr.name, v)
		}
		return tuple, nil
	default:
		return nil, errNeedRow
	}
}

func (r Relation) whereAtRow(ctx context.Context, scope Scope, fn *Function, ident string) (Value, bool, error) {
	ids := make([]uint32, 0, r.rows.n)
	for i := 0; i < r.rows.n; i++ {
		v, err := evalAtRow(ctx, fn.body, ident, r.rows.rowAt(i), r.attrMap, scope)
		if err != nil {
			if errors.Is(err, errNeedRow) {
				return nil, false, nil
			}
			return nil, false, err
		}
		if v.IsTrue() {
			ids = append(ids, r.rows.arenaID(i))
		}
	}
	if len(ids) == 0 {
		return None, true, nil
	}
	return r.newBody(r.rows.selView(ids)), true, nil
}

func (r Relation) mapAtRow(ctx context.Context, scope Scope, fn *Function, ident string) (Set, bool, error) {
	b := NewSetBuilder()
	for i := 0; i < r.rows.n; i++ {
		v, err := evalAtRow(ctx, fn.body, ident, r.rows.rowAt(i), r.attrMap, scope)
		if err != nil {
			if errors.Is(err, errNeedRow) {
				return nil, false, nil
			}
			return nil, false, err
		}
		b.Add(v)
	}
	s, err := b.Finish()
	if err != nil {
		return nil, false, err
	}
	return s, true, nil
}

// mapAddArrow evaluates `. +> (name: expr, ...)` as a wider Relation (🎯T29.9).
func (r Relation) mapAddArrow(ctx context.Context, scope Scope, te *TupleExpr, ident string) (Set, bool, error) {
	outNames := append(NamesSlice{}, r.attrs...)
	slot := make([]int, len(te.attrs))
	for j, a := range te.attrs {
		if a.IsWildcard() || a.name == "" {
			return nil, false, nil
		}
		if k, has := indexOfName(outNames, a.name); has {
			slot[j] = k
			continue
		}
		slot[j] = len(outNames)
		outNames = append(outNames, a.name)
	}
	sb := newStoreBuilder(len(outNames), r.rows.n, true)
	for i := 0; i < r.rows.n; i++ {
		row := r.rows.rowAt(i)
		out := make(Values, len(outNames))
		for j := range r.attrs {
			out[j] = row[r.p[j]]
		}
		for j, a := range te.attrs {
			v, err := evalAtRow(ctx, a.expr, ident, row, r.attrMap, scope)
			if err != nil {
				if errors.Is(err, errNeedRow) {
					return nil, false, nil
				}
				return nil, false, err
			}
			out[slot[j]] = v
		}
		sb.add(out)
	}
	return newRelation(outNames, makeIdentity(len(outNames)), sb.finish()), true, nil
}
