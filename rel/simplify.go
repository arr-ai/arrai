package rel

// Simplify is the tree-level rewrite layer (🎯T28): it runs on the compiled
// Expr tree after compilation and before evaluation or plan encoding, and
// is independent of value representation. See docs/simplifier.md.
//
// Rewrites, applied bottom-up to a fixpoint:
//
//   - One-shot let folding. `let x = rhs; body` compiles to an ArrowExpr
//     applying \x body once. When x occurs exactly once in body, in a
//     position evaluated exactly once (not under a lambda, a mapped body,
//     a conditional branch or a short-circuit rhs), no binder between the
//     let and that use captures a name rhs mentions, and rhs has no side
//     effects, the let becomes body[x := rhs]. The use site then sees the
//     rhs's shape, which the construction-time rewrites (predicate
//     pushdown, projection pruning, column extracts) can act on.
//   - Unused let dropping. When x does not occur in body and rhs is total
//     (cannot fail: a literal, a lambda, a bound identifier, or a tuple,
//     array or set of those), the let becomes body.
//
// Timing. A folded rhs is evaluated at its use instead of before body, so
// if both rhs and an earlier part of body would fail, the reported error
// changes; nothing else moves, because an rhs with side effects (a stdlib
// package with effects, or any call) is never folded, and a dropped rhs is
// one that could not have failed or acted.
//
// Occurrence counting is syntactic and conservative: an unknown node type
// (one children() cannot walk) makes the enclosing let ineligible.
func Simplify(e Expr) Expr {
	if !fastPaths {
		// The slowpath build keeps the tree as compiled, so the
		// differential-oracle CI job also covers every rewrite here.
		return e
	}
	for i := 0; i < simplifyMaxPasses; i++ {
		var changed bool
		e, changed = simplifyNode(e, nil)
		if !changed {
			break
		}
	}
	return e
}

const simplifyMaxPasses = 4

// SimplifyEnabled reports whether Simplify rewrites anything in this build;
// it is off under -tags slowpath. Tests that assert a tree's printed form
// choose their expectation with it.
func SimplifyEnabled() bool {
	return fastPaths
}

// simplifyNode rewrites e's children, then e itself. bound is the set of
// identifiers bound by enclosing binders.
func simplifyNode(e Expr, bound []string) (Expr, bool) {
	kids, rebuild, ok := children(e)
	if !ok {
		return e, false
	}
	changed := false
	var newKids []Expr
	for i, k := range kids {
		kb := bound
		if k.Binds != nil {
			kb = append(append([]string(nil), bound...), patternIdents(k.Binds)...)
		}
		nk, kc := simplifyNode(k.Expr, kb)
		if kc && newKids == nil {
			newKids = make([]Expr, len(kids))
			for j := range kids[:i] {
				newKids[j] = kids[j].Expr
			}
		}
		if newKids != nil {
			newKids[i] = nk
		}
		changed = changed || kc
	}
	if changed {
		e = rebuild(newKids)
	}
	if folded, did := simplifyLet(e, bound); did {
		return folded, true
	}
	return e, changed
}

// simplifyLet applies the let rewrites to e when e is a non-recursive
// `let ident = rhs; body`.
func simplifyLet(e Expr, bound []string) (Expr, bool) {
	arrow, is := e.(*ArrowExpr)
	if !is {
		return e, false
	}
	ident, is := arrow.fn.arg.(IdentPattern)
	if !is {
		return e, false
	}
	switch arrow.lhs.(type) {
	case RecursionExpr, *RecursionExpr:
		return e, false
	}
	name := string(ident)
	body := arrow.fn.body

	var a occurrences
	a.free = identsIn(arrow.lhs)
	a.walk(body, name, ChildOnce, nil)
	if a.opaque {
		return e, false
	}
	if a.count == 0 {
		if isTotal(arrow.lhs, bound) {
			return body, true
		}
		return e, false
	}
	if a.count == 1 && a.foldable && !hasEffects(arrow.lhs) {
		return substitute(body, name, arrow.lhs), true
	}
	return e, false
}

// occurrences is the result of occurrence analysis for one identifier.
type occurrences struct {
	free     map[string]bool // identifiers the rhs mentions
	count    int             // syntactic occurrences
	foldable bool            // the sole occurrence is once-evaluated and uncaptured
	opaque   bool            // an unwalkable node may hide occurrences
}

func (a *occurrences) walk(e Expr, name string, kind ChildKind, binders []Pattern) {
	if id, is := e.(IdentExpr); is {
		if id.ident == name {
			a.count++
			a.foldable = a.count == 1 && kind == ChildOnce && !a.captured(binders)
		}
		return
	}
	kids, _, ok := children(e)
	if !ok {
		a.opaque = true
		return
	}
	for _, k := range kids {
		kb := binders
		if k.Binds != nil {
			if shadows(k.Binds, name) {
				continue
			}
			kb = append(append([]Pattern(nil), binders...), k.Binds)
		}
		a.walk(k.Expr, name, kind.join(k.Kind), kb)
	}
}

// captured reports whether any binder on the path to a use binds a name
// the rhs mentions, which substitution would wrongly rebind.
func (a *occurrences) captured(binders []Pattern) bool {
	for _, p := range binders {
		for _, id := range patternIdents(p) {
			if a.free[id] {
				return true
			}
		}
	}
	return false
}

func shadows(p Pattern, name string) bool {
	for _, id := range patternIdents(p) {
		if id == name {
			return true
		}
	}
	return false
}

// identsIn returns every identifier mentioned anywhere in e, bound or not:
// a superset of its free identifiers, which is all capture checks need.
func identsIn(e Expr) map[string]bool {
	ids := map[string]bool{}
	var visit func(Expr)
	visit = func(e Expr) {
		if id, is := e.(IdentExpr); is {
			ids[id.ident] = true
			return
		}
		kids, _, ok := children(e)
		if !ok {
			return
		}
		for _, k := range kids {
			if k.Binds != nil {
				for _, id := range patternIdents(k.Binds) {
					ids[id] = true
				}
			}
			visit(k.Expr)
		}
	}
	visit(e)
	return ids
}

// substitute replaces the (single, unshadowed) occurrence of name in e.
func substitute(e Expr, name string, repl Expr) Expr {
	out, _ := substituteIn(e, name, repl)
	return out
}

// substituteIn is substitute with a changed flag; Exprs holding slices are
// not comparable, so identity cannot stand in for it.
func substituteIn(e Expr, name string, repl Expr) (Expr, bool) {
	if id, is := e.(IdentExpr); is {
		if id.ident == name {
			return repl, true
		}
		return e, false
	}
	kids, rebuild, ok := children(e)
	if !ok {
		return e, false
	}
	newKids := make([]Expr, len(kids))
	changed := false
	for i, k := range kids {
		newKids[i] = k.Expr
		if k.Binds != nil && shadows(k.Binds, name) {
			continue
		}
		if nk, kc := substituteIn(k.Expr, name, repl); kc {
			newKids[i] = nk
			changed = true
		}
	}
	if !changed {
		return e, false
	}
	return rebuild(newKids), true
}

// hasEffects reports whether evaluating e may have side effects. Stdlib
// packages declare theirs through Effectful; any call is assumed to, since
// the callee may.
func hasEffects(e Expr) bool {
	if ef, is := e.(Effectful); is && ef.HasSideEffects() {
		return true
	}
	if b, is := e.(*BinExpr); is && b.op == opCall {
		return true
	}
	kids, _, ok := children(e)
	if !ok {
		return true
	}
	for _, k := range kids {
		if hasEffects(k.Expr) {
			return true
		}
	}
	return false
}

// isTotal reports whether e cannot fail or act: dropping it unevaluated is
// unobservable.
func isTotal(e Expr, bound []string) bool {
	switch e := e.(type) {
	case LiteralExpr:
		return true
	case *Function:
		return true
	case IdentExpr:
		for _, b := range bound {
			if b == e.ident {
				return true
			}
		}
		return false
	case *TupleExpr:
		for _, a := range e.attrs {
			if !isTotal(a.expr, bound) {
				return false
			}
		}
		return true
	case ArrayExpr:
		for _, el := range e.elements {
			if el != nil && !isTotal(el, bound) {
				return false
			}
		}
		return true
	case *SetExpr:
		for _, el := range e.elements {
			if !isTotal(el, bound) {
				return false
			}
		}
		return true
	case Closure, ExprClosure:
		return true
	case Value:
		return true
	}
	return false
}
