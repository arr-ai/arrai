package rel

import (
	"context"
	"reflect"

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

func (ExprClosure) Count() int {
	return 1
}

func (ExprClosure) Has(Value) bool {
	panic("unimplemented")
}

func (ExprClosure) Enumerator() ValueEnumerator {
	panic("unimplemented")
}

func (c ExprClosure) With(Value) Set {
	panic("unimplemented")
}

func (ExprClosure) Without(Value) Set {
	panic("unimplemented")
}

func (ExprClosure) Map(func(Value) (Value, error)) (Set, error) {
	panic("unimplemented")
}

func (ExprClosure) Where(p func(v Value) (bool, error)) (Set, error) {
	panic("unimplemented")
}

func (c ExprClosure) CallAll(_ context.Context, arg Value) (Set, error) {
	//TODO: CallAll
	panic("unimplemented")
}

func (ExprClosure) ArrayEnumerator() (OffsetValueEnumerator, bool) {
	panic("unimplemented")
}
