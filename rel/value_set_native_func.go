package rel

import (
	"github.com/arr-ai/hash/hash128"

	"context"
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"github.com/arr-ai/arrai/pkg/fu"

	"github.com/arr-ai/wbnf/parser"
)

type NativeFnBody func(context.Context, Value) (Value, error)

var nativeByName sync.Map // string → NativeFnBody, for plan decode (🎯T25)

// NativeFunction represents a binary relation uniquely mapping inputs to outputs.
type NativeFunction struct {
	name string
	fn   NativeFnBody
}

// NewNativeFunction returns a new function.
func NewNativeFunction(name string, fn NativeFnBody) Value {
	full := "⦑" + name + "⦒"
	nativeByName.Store(full, fn)
	if name != "" {
		nativeByName.Store(name, fn)
	}
	return &NativeFunction{full, fn}
}

func lookupNative(name string) *NativeFunction {
	if v, ok := nativeByName.Load(name); ok {
		return &NativeFunction{name, v.(NativeFnBody)}
	}
	return nil
}

// NewNativeLambda returns a nameless function.
func NewNativeLambda(fn NativeFnBody) Value {
	return NewNativeFunction("", fn)
}

// NewNativeFunctionAttr returns a new Attr with a named key and NativeFunction value.
func NewNativeFunctionAttr(name string, fn NativeFnBody) Attr {
	return NewAttr(name, NewNativeFunction(name, fn))
}

// Name returns a native function's name.
func (f *NativeFunction) Name() string {
	return f.name
}

// Hash computes a hash for a NativeFunction.
func (f *NativeFunction) Hash(seed uintptr) uintptr {
	return f.Hash128().Seeded(seed)
}

// Hash128 computes the 128-bit hash of a NativeFunction, by identity, which
// is how Equal compares them.
func (f *NativeFunction) Hash128() hash128.H128 {
	return hash128.Uintptr(uintptr(unsafe.Pointer(f)))
}

// Equal tests two Values for equality. Any other type returns false.
func (f *NativeFunction) Equal(i Value) bool {
	if g, ok := i.(*NativeFunction); ok {
		return f == g
	}
	return false
}

// String returns a string representation of the expression.
func (f *NativeFunction) String() string {
	return fu.String(f)
}

// Format formats the expression.
func (f *NativeFunction) Format(s fmt.State, verb rune) {
	fu.WriteString(s, f.name)
}

// Eval returns the Value
func (f *NativeFunction) Eval(ctx context.Context, local Scope) (Value, error) {
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
func (f *NativeFunction) Export(_ context.Context) interface{} {
	return f.fn
}

func (*NativeFunction) getSetBuilder() setBuilder {
	return newGenericTypeSetBuilder()
}

func (*NativeFunction) getBucket() fmt.Stringer {
	return genericType
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

func (*NativeFunction) Map(func(Value) (Value, error)) (Set, error) {
	panic("unimplemented")
}

func (*NativeFunction) Where(p func(v Value) (bool, error)) (Set, error) {
	panic("unimplemented")
}

// Call calls the NativeFunction with the given parameter.
func (f *NativeFunction) CallAll(ctx context.Context, arg Value, b SetBuilder) error {
	arg, err := Observe(arg)
	if err != nil {
		return err
	}
	v, err := f.fn(ctx, arg)
	if err != nil {
		return err
	}
	b.Add(v)
	return nil
}

func (*NativeFunction) unionSetSubsetBucket() string {
	// TODO: create its own subset bucket in unionset
	return genericType.String()
}

func (*NativeFunction) ArrayEnumerator() ValueEnumerator {
	panic("unimplemented")
}
