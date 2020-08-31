package rel

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/arr-ai/hash"
	"github.com/arr-ai/wbnf/parser"
)

// Function represents a binary relation uniquely mapping inputs to outputs.
type Function struct {
	ExprScanner
	arg  Pattern
	body Expr
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

// Hash computes a hash for a Function.
func (f *Function) Hash(seed uintptr) uintptr {
	//TODO: function should be an expr but hash is called by Closure
	return hash.String(f.String(), hash.Uintptr(17297263775284131973>>(64-8*unsafe.Sizeof(uintptr(0))), seed))
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
	// Function equality is undecidable in the general case. Should we panic?
	return f.body == g.body
}

// String returns a string representation of the expression.
func (f *Function) String() string {
	if f.arg.String() == "-" {
		return fmt.Sprintf("(&%s)", f.body)
	}
	return fmt.Sprintf("(\\%s %s)", f.arg, f.body)
}

// Eval returns the Value
func (f *Function) Eval(ctx context.Context, local Scope) (Value, error) {
	return NewClosure(local, f), nil
}
