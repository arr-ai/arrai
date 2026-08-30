package rel

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/arr-ai/wbnf/parser"
)

// IdentExpr returns the variable referenced by ident.
type IdentExpr struct {
	ExprScanner
	ident string
	// where caches the binding's address as (hops << 32 | slot) + 1, or 0
	// when unknown. Scoping is lexical, so the address is stable across
	// evaluations; it is verified on every use and refreshed on a miss, so
	// a wrong entry costs a walk, never a wrong answer. Shared by value
	// copies of the expression; atomic because predicates can be evaluated
	// concurrently.
	where *atomic.Uint64
}

// NewIdentExpr returns a new identifier.
func NewIdentExpr(scanner parser.Scanner, ident string) Expr {
	if isDynIdent(ident) {
		return DynIdentExpr{IdentExpr: IdentExpr{ExprScanner: ExprScanner{scanner}, ident: ident}}
	}
	return IdentExpr{ExprScanner: ExprScanner{scanner}, ident: ident, where: &atomic.Uint64{}}
}

func NewDotIdent(source parser.Scanner) Expr {
	return NewIdentExpr(source, ".")
}

// Ident returns the ident for the IdentExpr.
func (e IdentExpr) Ident() string {
	return e.ident
}

// String returns a string representation of the expression.
func (e IdentExpr) String() string {
	if e.ident == "." {
		return "(" + e.ident + ")"
	}
	return e.ident
}

// Eval returns the value from scope corresponding to the ident.
func (e IdentExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	if a := e.resolve(local); a != nil {
		return a.Eval(ctx, local)
	}
	return nil, WrapContextErr(
		fmt.Errorf("name %q not found in {%s}", e.ident, strings.Join(local.OrderedNames(), ", ")), e, local)
}

// resolve returns the ident's binding in local, or nil. It tries the cached
// address first and verifies the name at that slot before trusting it.
func (e IdentExpr) resolve(local Scope) Expr {
	if e.where != nil && fastPaths {
		if w := e.where.Load(); w != 0 {
			w--
			hops, slot := int(w>>32), int(w&0xffffffff)
			if f := local.at(hops); f != nil && slot < len(f.names) && f.names[slot] == e.ident {
				return f.vals[slot]
			}
		}
	}
	f, slot, hops := local.lookup(e.ident)
	if f == nil {
		return nil
	}
	if e.where != nil {
		e.where.Store(uint64(hops)<<32 | uint64(slot) + 1)
	}
	return f.vals[slot]
}

type DynIdentExpr struct {
	IdentExpr
}

// Eval returns the value from scope corresponding to the ident.
func (e DynIdentExpr) Eval(ctx context.Context, local Scope) (Value, error) {
	if a := ctx.Value(DynIdent(e.ident)); a != nil {
		return a.(Value), nil
	}
	return nil, WrapContextErr(fmt.Errorf("dynamic variable %s not found", e.ident), e, local)
}
