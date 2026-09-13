package rel

import (
	"context"
	"fmt"
	"sync/atomic"
	"unsafe"

	"github.com/arr-ai/arrai/pkg/fu"

	"github.com/arr-ai/wbnf/parser"
)

// Function represents a binary relation uniquely mapping inputs to outputs.
type Function struct {
	ExprScanner
	arg  Pattern
	body Expr
	// columnOnly is 0 unknown, 1 yes, 2 no (ident used as a Value). 🎯T29.9
	columnOnly uint32
}

// NewFunction returns a new function.
func NewFunction(scanner parser.Scanner, arg Pattern, body Expr) Expr {
	return &Function{ExprScanner: ExprScanner{Src: scanner}, arg: arg, body: body}
}

// ExprAsFunction returns a function for an expr. If the expr is already a
// function, returns expr. Otherwise, returns expr wrapper in a function with
// arg '.'.
func ExprAsFunction(expr Expr) *Function {
	if fn, ok := expr.(*Function); ok {
		return fn
	}
	return NewFunction(expr.Source(), IdentPattern("."), expr).(*Function)
}

// Arg returns a function's formal argument.
func (f *Function) Arg() string {
	return f.arg.String()
}

// Body returns a function's body.
func (f *Function) Body() Expr {
	return f.body
}

// isColumnOnly reports whether the ident formal is only read as ident.attr
// (or as the row extended by +>), never as a Value. Cached on the Function
// so plan eval does not re-walk the body (🎯T29.9).
func (f *Function) isColumnOnly() bool {
	if ident, ok := f.arg.(IdentPattern); ok {
		switch atomic.LoadUint32(&f.columnOnly) {
		case 1:
			return true
		case 2:
			return false
		}
		if !usesIdentAsValue(f.body, string(ident)) {
			atomic.StoreUint32(&f.columnOnly, 1)
			return true
		}
		atomic.StoreUint32(&f.columnOnly, 2)
	}
	return false
}

// Hash computes a hash for a Function. Functions are compared by identity
// (see EqualFunction), so the hash is derived from the node's address rather
// than from formatting its body, which is expensive and, for recursive
// functions, self-referential.
func (f *Function) Hash() uintptr {
	return mix(funcSalt, uintptr(unsafe.Pointer(f)))
}

// Equal tests two Values for equality. Any other type returns false.
func (f *Function) Equal(i interface{}) bool {
	// Function equality is undecidable in the general case. Should we panic?
	if g, ok := i.(*Function); ok {
		return f.EqualFunction(g)
	}
	return false
}

// Equal tests two Values for equality. Any other type returns false.
func (f *Function) EqualFunction(g *Function) bool {
	// Function equality is undecidable in the general case, so functions are
	// equal iff they are the same compiled node. (Comparing bodies with ==
	// panics when the body is an uncomparable expression type.)
	return f == g
}

// String returns a string representation of the expression.
func (f *Function) String() string {
	return fu.String(f)
}

// Format formats the expression.
func (f *Function) Format(s fmt.State, verb rune) {
	if f.arg.String() == "-" {
		fu.Fprintf(s, "(&%s)", f.body)
	} else {
		fu.Fprintf(s, "(\\%s %s)", f.arg, f.body)
	}
}

// Eval returns the Value
func (f *Function) Eval(ctx context.Context, local Scope) (Value, error) {
	return NewClosure(local, f), nil
}
