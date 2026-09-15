package rel

import (
	"context"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/arr-ai/arrai/pkg/fu"

	"github.com/arr-ai/wbnf/parser"
)

// ExprClosure represents the closure of an expression over a scope.
type ExprClosure struct {
	scope Scope
	e     Expr
}

// NewExprClosure returns a new ExprClosure.
func NewExprClosure(scope Scope, e Expr) Value {
	return ExprClosure{scope: scope, e: e}
}

// Hash computes a hash for a ExprClosure. It folds in only the captured
// scope's frame, which every closure equal under EqualExprClosure shares.
func (c ExprClosure) Hash() uintptr {
	return mix(exprClosureSalt, uintptr(unsafe.Pointer(c.scope.f)))
}

// Equal tests two Values for equality. Any other type returns false.
func (c ExprClosure) Equal(i Value) bool {
	if d, ok := i.(ExprClosure); ok {
		return c.EqualExprClosure(d)
	}
	return false
}

// EqualExprClosure tests two ExprClosures for equality. As with Function,
// equality of the wrapped expressions is undecidable in general, so two
// ExprClosures are equal iff they capture the same scope frame and wrap the
// same expression node. Expression types that are not comparable (those
// holding slices or maps) are never equal, and never panic here.
func (c ExprClosure) EqualExprClosure(d ExprClosure) bool {
	if c.scope.f != d.scope.f {
		return false
	}
	if c.e == nil || d.e == nil {
		return c.e == nil && d.e == nil
	}
	ce, de := reflect.ValueOf(c.e), reflect.ValueOf(d.e)
	return ce.Type() == de.Type() && ce.Comparable() && ce.Equal(de)
}

// String returns a string representation of the expression.
func (c ExprClosure) String() string {
	return fu.String(c)
}

// Format formats the expression.
func (c ExprClosure) Format(f fmt.State, verb rune) {
	fu.WriteString(f, "◖")
	fu.FRepr(f, c.e)
	fu.WriteString(f, "◗")
}

// Eval returns the Value
func (c ExprClosure) Eval(ctx context.Context, _ Scope) (Value, error) {
	return c.e.Eval(ctx, c.scope)
}

// Source returns a scanner locating the ExprClosure's source code.
func (c ExprClosure) Source() parser.Scanner {
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
func (c ExprClosure) Export(ctx context.Context) interface{} {
	return func(v Value) Value {
		result, err := SetCall(ctx, c, v)
		if err != nil {
			panic(err)
		}
		return result
	}
}

func (ExprClosure) getSetBuilder() setBuilder {
	return newGenericTypeSetBuilder()
}

func (ExprClosure) getBucket() fmt.Stringer {
	return genericType
}

func (ExprClosure) Count() int {
	return 1
}

// An ExprClosure is the one-element set {c}; see value_set_singleton.go.

func (c ExprClosure) Has(v Value) bool {
	return c.Equal(v)
}

func (c ExprClosure) Enumerator() ValueEnumerator {
	return &singletonEnumerator{v: c}
}

func (c ExprClosure) With(v Value) Set {
	return singletonWith(c, v)
}

func (c ExprClosure) Without(v Value) Set {
	return singletonWithout(c, v)
}

func (c ExprClosure) Map(f func(Value) (Value, error)) (Set, error) {
	return singletonMap(c, f)
}

func (c ExprClosure) Where(p func(v Value) (bool, error)) (Set, error) {
	return singletonWhere(c, p)
}

func (c ExprClosure) CallAll(_ context.Context, arg Value, b SetBuilder) error {
	panic("unimplemented")
}

func (ExprClosure) unionSetSubsetBucket() string {
	// TODO: create its own subset bucket in union set
	return genericType.String()
}

func (c ExprClosure) ArrayEnumerator() ValueEnumerator {
	return &singletonEnumerator{v: c}
}
