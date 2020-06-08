package rel

import (
	"reflect"
	"unsafe"

	"github.com/arr-ai/hash"
	"github.com/arr-ai/wbnf/parser"
)

// NativeFunction represents a binary relation uniquely mapping inputs to outputs.
type NativeFunction struct {
	name string
	fn   func(Value) (Value, error)
}

// NewNativeFunction returns a new function.
func NewNativeFunction(name string, fn func(Value) (Value, error)) Value {
	return &NativeFunction{"⦑" + name + "⦒", fn}
}

// NewNativeLambda returns a nameless function.
func NewNativeLambda(fn func(Value) (Value, error)) Value {
	return NewNativeFunction("", fn)
}

// NewNativeFunction returns a new function.
func NewNativeFunctionAttr(name string, fn func(Value) (Value, error)) Attr {
	return NewAttr(name, NewNativeFunction(name, fn))
}

// Name returns a native function's name.
func (f *NativeFunction) Name() string {
	return f.name
}

// Hash computes a hash for a NativeFunction.
func (f *NativeFunction) Hash(seed uintptr) uintptr {
	return hash.String(f.String(), hash.Uintptr(9714745597188477233>>(64-8*unsafe.Sizeof(uintptr(0))), seed))
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

// Source returns an empty scanner since NativeFunction doesn't have associated
// source code.
func (f *NativeFunction) Source() parser.Scanner {
	return *parser.NewScanner("")
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

func (*NativeFunction) Count() int {
	return 1
}

func (*NativeFunction) Has(Value) bool {
	panic("unimplemented")
}

func (*NativeFunction) Enumerator() ValueEnumerator {
	panic("unimplemented")
}

func (*NativeFunction) With(Value) Set {
	panic("unimplemented")
}

func (*NativeFunction) Without(Value) Set {
	panic("unimplemented")
}

func (*NativeFunction) Map(func(Value) Value) Set {
	panic("unimplemented")
}

func (*NativeFunction) Where(p func(v Value) (bool, error)) (Set, error) {
	panic("unimplemented")
}

// Call calls the NativeFunction with the given parameter.
func (f *NativeFunction) CallAll(arg Value) (Set, error) {
	v, err := f.fn(arg)
	if err != nil {
		return nil, err
	}
	return NewSet(v), nil
}

func (*NativeFunction) ArrayEnumerator() (OffsetValueEnumerator, bool) {
	panic("unimplemented")
}
