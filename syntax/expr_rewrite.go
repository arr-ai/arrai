package syntax

import "github.com/arr-ai/arrai/rel"

// Tree-rewrite support for this package's own Expr types (🎯T28); see
// rel.Rewritable and rel.Simplify.

// effectfulPackages are the stdlib packages whose functions have side
// effects or read the environment, so a rewrite must not move a reference
// to them relative to the rest of the program.
var effectfulPackages = map[string]bool{
	"log": true, "os": true, "net": true, "eval": true, "test": true, "tst": true, "runtime": true,
}

func (e PackageExpr) Children() []rel.Child {
	return []rel.Child{{Expr: e.a, Kind: rel.ChildOnce}}
}

func (e PackageExpr) WithChildren(kids []rel.Expr) rel.Expr {
	return NewPackageExpr(e.Src, kids[0])
}

// HasSideEffects reports whether the package is one whose functions act on
// the outside world; `//str`, `//seq` and friends are pure.
func (e PackageExpr) HasSideEffects() bool {
	if d, is := e.a.(*rel.DotExpr); is {
		return effectfulPackages[d.Attr()]
	}
	return true
}

// An import is a leaf: its module evaluates in its own empty scope.
func (i ImportExpr) Children() []rel.Child { return nil }

func (i ImportExpr) WithChildren([]rel.Expr) rel.Expr { return i }

func (e *xstrExpr) Children() []rel.Child {
	var kids []rel.Child
	for _, p := range e.parts {
		if p.expr != nil {
			kids = append(kids, rel.Child{Expr: p.expr, Kind: rel.ChildOnce})
		}
	}
	return kids
}

func (e *xstrExpr) WithChildren(kids []rel.Expr) rel.Expr {
	n := &xstrExpr{ExprScanner: e.ExprScanner, parts: make([]xstrPart, len(e.parts))}
	j := 0
	for i, p := range e.parts {
		n.parts[i] = p
		if p.expr != nil {
			n.parts[i].expr = kids[j]
			j++
		}
	}
	return n
}
