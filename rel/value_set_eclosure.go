package rel

import (
	"reflect"

	"github.com/arr-ai/wbnf/parser"
)

// ExprClosure represents the closure of an expression over a scope.
type ExprClosure struct {
	scope Scope
	e     Expr
}

// NewFunction returns a new function.
func NewExprClosure(scope Scope, e Expr) Value {
	return ExprClosure{scope: scope, e: e}
}

// Hash computes a hash for a ExprClosure.
func (c ExprClosure) Hash(seed uintptr) uintptr {
	panic("not implemented")
	// TODO: Is this enough?
	// return c.e.Hash(seed)
}

// Equal tests two Values for equality. Any other type returns false.
func (c ExprClosure) Equal(i interface{}) bool {
	if d, ok := i.(ExprClosure); ok {
		return c.EqualExprClosure(d)
	}
	return false
}

// Equal tests two Values for equality. Any other type returns false.
func (c ExprClosure) EqualExprClosure(d ExprClosure) bool {
	panic("not implemented")
	// return c.f.EqualFunction(d.f)
}

// String returns a string representation of the expression.
func (c ExprClosure) String() string {
	return "◖" + c.e.String() + "◗"
}

// Eval returns the Value
func (c ExprClosure) Eval(_ Scope) (Value, error) {
	return c.e.Eval(c.scope)
}

// Scanner returns the scanner of ExprClosure.
func (c ExprClosure) Scanner() parser.Scanner {
	return *parser.NewScanner("")
}

var eclosureKind = registerKind(206, reflect.TypeOf(ExprClosure{}))

// Kind returns a number that is unique for each major kind of Value.
func (c ExprClosure) Kind() int {
	return eclosureKind
}

// Bool returns true iff the tuple has attributes.
func (c ExprClosure) IsTrue() bool {
	return true
}

// Less returns true iff g is not a number or f.number < g.number.
func (c ExprClosure) Less(d Value) bool {
	if c.Kind() != d.Kind() {
		return c.Kind() < d.Kind()
	}
	return c.String() < d.String()
}

// Negate returns {(negateTag): f}.
func (c ExprClosure) Negate() Value {
	return NewTuple(NewAttr(negateTag, c))
}

// Export exports a ExprClosure.
func (c ExprClosure) Export() interface{} {
	return func(_ Value, local Scope) (Value, error) {
		return c.Call(None, local)
	}
}

// Call calls the ExprClosure with the given parameter.
func (c ExprClosure) Call(expr Expr, local Scope) (Value, error) {
	return c.e.Eval(local)
}
