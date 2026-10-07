package rel

import "slices"

// freeIdents returns the identifiers e refers to without binding them itself,
// in first-use order, and whether that set is known. It is unknown (ok ==
// false) when e holds a node children cannot see into, since such a node may
// refer to anything in scope. Dynamic (@-prefixed) identifiers resolve through
// the context, not the scope, and are not reported.
func freeIdents(e Expr) (names []string, ok bool) {
	ok = true
	var visit func(e Expr, bound []string)
	visit = func(e Expr, bound []string) {
		if e == nil || !ok {
			return
		}
		if id, is := e.(IdentExpr); is {
			if !slices.Contains(bound, id.ident) && !slices.Contains(names, id.ident) {
				names = append(names, id.ident)
			}
			return
		}
		kids, _, known := children(e)
		if !known {
			ok = false
			return
		}
		for _, k := range kids {
			kb := bound
			if k.Binds != nil {
				// Copy on extend: sibling children must not see each other's binders.
				kb = append(bound[:len(bound):len(bound)], patternIdents(k.Binds)...)
			}
			visit(k.Expr, kb)
		}
	}
	visit(e, nil)
	return names, ok
}
