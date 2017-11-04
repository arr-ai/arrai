package rel

import (
	"encoding/binary"
	"fmt"

	"github.com/OneOfOne/xxhash"
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
func (f *Function) Hash(seed uint32) uint32 {
	xx := xxhash.NewS32(seed ^ 0x734be5ff)
	binary.Write(xx, binary.LittleEndian, f.String())
	return xx.Sum32()
}

// Equal tests two Values for equality. Any other type returns false.
func (f *Function) Equal(i interface{}) bool {
	if g, ok := i.(*Function); ok {
		return f.body == g.body
	}
	return false
}

// String returns a string representation of the expression.
func (f *Function) String() string {
	if f.arg == "-" {
		return fmt.Sprintf("(& %s)", f.body)
	}
	return fmt.Sprintf("(\\%s %s)", f.arg, f.body)
}

// Eval returns the Value
func (f *Function) Eval(local, global *Scope) (Value, error) {
	return f, nil
}

// Kind returns a number that is unique for each major kind of Value.
func (f *Function) Kind() int {
	return 202
}

// Bool returns true iff the tuple has attributes.
func (f *Function) Bool() bool {
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
		return func(local *Scope, global *Scope) (Value, error) {
			return f.Call(None, local, global)
		}
	}
	return func(e Value, local *Scope, global *Scope) (Value, error) {
		return f.body.Eval(local.With(f.arg, e), global)
	}
}

// Call calls the Function with the given parameter.
func (f *Function) Call(expr Expr, local, global *Scope) (Value, error) {
	niladic := f.arg == "-"
	noArg := expr == nil
	if niladic != noArg {
		return nil, errors.Errorf(
			"nullary-vs-unary function arg mismatch (%s vs %s)", f.arg, expr)
	}
	if niladic {
		return f.body.Eval(local, global)
	}
	return f.body.Eval(local.With(f.arg, expr), global)
}
