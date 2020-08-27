package rel

//TODO: reimplement NativeExprFunction
// import (
// 	"reflect"
// 	"unsafe"

// 	"github.com/arr-ai/hash"
// 	"github.com/arr-ai/wbnf/parser"
// )

// // NativeExprFunction represents a binary relation uniquely mapping inputs to outputs.
// type NativeExprFunction struct {
// 	name string
// 	fn   func(Expr, Scope) (Value, error)
// }

// // NewNativeFunction returns a new function.
// func NewNativeExprFunction(name string, fn func(Expr, Scope) (Value, error)) Value {
// 	return &NativeExprFunction{"⦑" + name + "⦒", fn}
// }

// // NewNativeFunction returns a new function.
// func NewNativeExprFunctionAttr(name string, fn func(Expr, Scope) (Value, error)) Attr {
// 	return NewAttr(name, NewNativeExprFunction(name, fn))
// }

// // Name returns a native function's name.
// func (f *NativeExprFunction) Name() string {
// 	return f.name
// }

// // Hash computes a hash for a NativeExprFunction.
// func (f *NativeExprFunction) Hash(seed uintptr) uintptr {
// 	return hash.String(f.String(), hash.Uintptr(9714745597188477233>>(64-8*unsafe.Sizeof(uintptr(0))), seed))
// }

// // Equal tests two Values for equality. Any other type returns false.
// func (f *NativeExprFunction) Equal(i interface{}) bool {
// 	if g, ok := i.(*NativeExprFunction); ok {
// 		return f == g
// 	}
// 	return false
// }

// // String returns a string representation of the expression.
// func (f *NativeExprFunction) String() string {
// 	return f.name
// }

// // Eval returns the Value
// func (f *NativeExprFunction) Eval(ctx context.Context, local Scope) (Value, error) {
// 	return f, nil
// }

// // Source returns an empty scanner since NativeExprFunction doesn't have
// // associated source code.
// func (f *NativeExprFunction) Source() parser.Scanner {
// 	return *parser.NewScanner("")
// }

// var nativeExprFunctionKind = registerKind(210, reflect.TypeOf(NativeExprFunction{}))

// // Kind returns a number that is unique for each major kind of Value.
// func (f *NativeExprFunction) Kind() int {
// 	return nativeExprFunctionKind
// }

// // Bool always returns true.
// func (f *NativeExprFunction) IsTrue() bool {
// 	return true
// }

// // Less returns true iff g is not a number or f.number < g.number.
// func (f *NativeExprFunction) Less(g Value) bool {
// 	if f.Kind() != g.Kind() {
// 		return f.Kind() < g.Kind()
// 	}
// 	return f.String() < g.String()
// }

// // Negate returns {(negateTag): f}.
// func (f *NativeExprFunction) Negate() Value {
// 	return NewTuple(NewAttr(negateTag, f))
// }

// // Export exports a NativeExprFunction.
// func (f *NativeExprFunction) Export(_ context.Context, ) interface{} {
// 	return f.fn
// }

// // Call calls the NativeExprFunction with the given parameter.
// func (f *NativeExprFunction) Call(expr Expr, local Scope) (_ Value, err error) {
// 	return f.fn(expr, local)
// }
