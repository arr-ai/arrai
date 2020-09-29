package rel

import (
	"context"
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
func NewClosure(scope Scope, f *Function) Closure {
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
	return c.f.String()
}

// Eval returns the Value
func (c Closure) Eval(ctx context.Context, local Scope) (Value, error) {
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
func (c Closure) Export(ctx context.Context) interface{} {
	if c.f.Arg() == "-" {
		result, err := SetCall(ctx, c, None)
		if err != nil {
			panic(err)
		}
		return result.Export(ctx)
	}
	return func(arg Value) Value {
		result, err := SetCall(ctx, c, None)
		if err != nil {
			panic(err)
		}
		return result
	}
}

func (c Closure) Count() int {
	panic("unimplemented")
}

func (c Closure) Has(v Value) bool {
	panic("unimplemented")
}

func (c Closure) Enumerator() ValueEnumerator {
	panic("unimplemented")
}

func (c Closure) With(v Value) Set {
	panic("unimplemented")
}

func (c Closure) Without(v Value) Set {
	panic("unimplemented")
}

func (c Closure) Map(f func(v Value) (Value, error)) (Set, error) {
	panic("unimplemented")
}

func (c Closure) Where(p func(v Value) (bool, error)) (Set, error) {
	panic("unimplemented")
}

//FIXME: context not used properly
func (c Closure) CallAll(ctx context.Context, arg Value) (Set, error) {
	niladic := c.f.Arg() == "-"
	noArg := arg == nil
	if niladic != noArg {
		panic(errors.Errorf(
			"nullary-vs-unary function arg mismatch (%s vs %s)", c.f.Arg(), arg))
	}
	if niladic {
		val, err := c.f.body.Eval(ctx, c.scope)
		if err != nil {
			return nil, err
		}
		return NewSet(val)
	}
	ctx, scope, err := c.f.arg.Bind(ctx, c.scope, arg)
	if err != nil {
		return nil, err
	}
	val, err := c.f.body.Eval(ctx, c.scope.Update(scope))
	if err != nil {
		return nil, err
	}
	return NewSet(val)
}

func (c Closure) ArrayEnumerator() (OffsetValueEnumerator, bool) {
	panic("unimplemented")
}
