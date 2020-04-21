package rel

import (
	"reflect"

	"github.com/arr-ai/hash"
	"github.com/arr-ai/wbnf/parser"
)

// NativeFunction represents a binary relation uniquely mapping inputs to outputs.
type NativeFunction struct {
	name string
	fn   func(Value) Value
}

// NewNativeFunction returns a new function.
func NewNativeFunction(name string, fn func(Value) Value) Value {
	return &NativeFunction{"⧼" + name + "⧽", fn}
}

// NewNativeLambda returns a nameless function.
func NewNativeLambda(fn func(Value) Value) Value {
	return NewNativeFunction("", fn)
}

// NewNativeFunction returns a new function.
func NewNativeFunctionAttr(name string, fn func(Value) Value) Attr {
	return NewAttr(name, NewNativeFunction(name, fn))
}

// Name returns a native function's name.
func (f *NativeFunction) Name() string {
	return f.name
}

// Fn returns a native function's implementation.
func (f *NativeFunction) Fn() func(Value) Value {
	return f.fn
}

// Hash computes a hash for a NativeFunction.
func (f *NativeFunction) Hash(seed uintptr) uintptr {
	return hash.String(f.String(), hash.Uintptr(9714745597188477233, seed))
}

// Equal tests two Values for equality. Any other type returns false.
func (f *NativeFunction) Equal(i interface{}) bool {
	if g, ok := i.(*NativeFunction); ok {
		return f == g
	}
	return false
}

// String returns a string representation of the expression.
func (f *NativeFunction) String() string {
	return f.name
}

// Eval returns the Value
func (f *NativeFunction) Eval(local Scope) (Value, error) {
	return f, nil
}

// Scanner returns the scanner of NativeFunction.
func (f *NativeFunction) Scanner() parser.Scanner {
	panic("not implemented")
}

var nativeFunctionKind = registerKind(203, reflect.TypeOf(NativeFunction{}))

// Kind returns a number that is unique for each major kind of Value.
func (f *NativeFunction) Kind() int {
	return nativeFunctionKind
}

// Bool always returns true.
func (f *NativeFunction) IsTrue() bool {
	return true
}

// Less returns true iff g is not a number or f.number < g.number.
func (f *NativeFunction) Less(g Value) bool {
	if f.Kind() != g.Kind() {
		return f.Kind() < g.Kind()
	}
	return f.String() < g.String()
}

// Negate returns {(negateTag): f}.
func (f *NativeFunction) Negate() Value {
	return NewTuple(NewAttr(negateTag, f))
}

// Export exports a NativeFunction.
func (f *NativeFunction) Export() interface{} {
	return f.fn
}

// Call calls the NativeFunction with the given parameter.
func (f *NativeFunction) Call(expr Expr, local Scope) (Value, error) {
	if expr == nil {
		return f.fn(nil), nil
	}
	value, err := expr.Eval(local)
	if err != nil {
		return nil, err
	}
	return f.fn(value), nil
}
