package rel

import (
	"reflect"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// Closure represents the closure of a function over a scope.
type Closure struct {
	scope Scope
	f     *Function
}

// NewFunction returns a new function.
func NewClosure(scope Scope, f *Function) Value {
	return Closure{scope: scope, f: f}
}

// Hash computes a hash for a Closure.
func (c Closure) Hash(seed uintptr) uintptr {
	// TODO: Is this enough?
	return c.f.Hash(seed)
}

// Equal tests two Values for equality. Any other type returns false.
func (c Closure) Equal(i interface{}) bool {
	if d, ok := i.(Closure); ok {
		return c.f.Equal(d.f)
	}
	return false
}

// Equal tests two Values for equality. Any other type returns false.
func (c Closure) EqualClosure(d Closure) bool {
	return c.f.EqualFunction(d.f)
}

// String returns a string representation of the expression.
func (c Closure) String() string {
	return "⦇" + c.f.String() + "⦈"
}

// Eval returns the Value
func (c Closure) Eval(local Scope) (Value, error) {
	return c, nil
}

// Source returns a scanner locating the Closure's source code.
func (c Closure) Source() parser.Scanner {
	return *parser.NewScanner("")
}

var closureKind = registerKind(205, reflect.TypeOf(Closure{}))

// Kind returns a number that is unique for each major kind of Value.
func (c Closure) Kind() int {
	return closureKind
}

// Bool returns true iff the tuple has attributes.
func (c Closure) IsTrue() bool {
	return true
}

// Less returns true iff g is not a number or f.number < g.number.
func (c Closure) Less(d Value) bool {
	if c.Kind() != d.Kind() {
		return c.Kind() < d.Kind()
	}
	return c.String() < d.String()
}

// Negate returns {(negateTag): f}.
func (c Closure) Negate() Value {
	return NewTuple(NewAttr(negateTag, c))
}

// Export exports a Closure.
func (c Closure) Export() interface{} {
	if c.f.arg == "-" {
		return func(_ Value, local Scope) (Value, error) {
			return c.Call(None, local)
		}
	}
	return func(e Value, local Scope) (Value, error) {
		return c.f.body.Eval(local.With(c.f.arg, e))
	}
}

// Call calls the Closure with the given parameter.
func (c Closure) Call(expr Expr, local Scope) (Value, error) {
	niladic := c.f.arg == "-"
	noArg := expr == nil
	if niladic != noArg {
		return nil, errors.Errorf(
			"nullary-vs-unary function arg mismatch (%s vs %s)", c.f.arg, expr)
	}
	if niladic {
		return c.f.body.Eval(local)
	}
	return c.f.body.Eval(c.scope.With(c.f.arg, expr))
}
