package rel

import (
	"fmt"
	"reflect"

	"github.com/arr-ai/hash"
	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// Function represents a binary relation uniquely mapping inputs to outputs.
type Function struct {
	arg  string
	body Expr
}

// NewFunction returns a new function.
func NewFunction(arg string, body Expr) Expr {
	return &Function{arg, body}
}

// ExprAsFunction returns a function for an expr. If the expr is already a
// function, returns expr. Otherwise, returns expr wrapper in a function with
// arg '.'.
func ExprAsFunction(expr Expr) *Function {
	if fn, ok := expr.(*Function); ok {
		return fn
	}
	return NewFunction(".", expr).(*Function)
}

// Arg returns a function's formal argument.
func (f *Function) Arg() string {
	return f.arg
}

// Body returns a function's body.
func (f *Function) Body() Expr {
	return f.body
}

// Hash computes a hash for a Function.
func (f *Function) Hash(seed uintptr) uintptr {
	return hash.String(f.String(), hash.Uintptr(17297263775284131973, seed))
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
	if f.arg == "-" {
		return fmt.Sprintf("(&%s)", f.body)
	}
	return fmt.Sprintf("(\\%s %s)", f.arg, f.body)
}

// Eval returns the Value
func (f *Function) Eval(local Scope) (Value, error) {
	return NewClosure(local, f), nil
}

// Source returns a scanner locating the Function's source code.
func (f *Function) Source() parser.Scanner {
	return *parser.NewScanner("")
}

var functionKind = registerKind(202, reflect.TypeOf(Function{}))

// Kind returns a number that is unique for each major kind of Value.
func (f *Function) Kind() int {
	return functionKind
}

// Bool returns true iff the tuple has attributes.
func (f *Function) IsTrue() bool {
	return true
}

// Less returns true iff g is not a number or f.number < g.number.
func (f *Function) Less(g Value) bool {
	if f.Kind() != g.Kind() {
		return f.Kind() < g.Kind()
	}
	return f.String() < g.String()
}

// Negate returns {(negateTag): f}.
func (f *Function) Negate() Value {
	return NewTuple(NewAttr(negateTag, f))
}

// Export exports a Function.
func (f *Function) Export() interface{} {
	if f.arg == "-" {
		return func(local Scope) (Value, error) {
			return f.Call(None, local)
		}
	}
	return func(e Value, local Scope) (Value, error) {
		return f.body.Eval(local.With(f.arg, e))
	}
}

// Call calls the Function with the given parameter.
func (f *Function) Call(expr Expr, local Scope) (Value, error) {
	niladic := f.arg == "-"
	noArg := expr == nil
	if niladic != noArg {
		return nil, errors.Errorf(
			"nullary-vs-unary function arg mismatch (%s vs %s)", f.arg, expr)
	}
	if niladic {
		return f.body.Eval(local)
	}
	return f.body.Eval(local.With(f.arg, expr))
}
