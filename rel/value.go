package rel

import (
	"context"
	"fmt"

	"github.com/arr-ai/frozen"
	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

// Expr represents an arr.ai expression.
type Expr interface {
	// Require a String() method.
	fmt.Stringer

	// Eval evaluates the expr in a given scope.
	Eval(ctx context.Context, local Scope) (Value, error)

	Source() parser.Scanner
}

// Value represents any arr.ai value.
type Value interface {
	frozen.Key

	// Values are Exprs.
	Expr

	// Kind returns a number that is unique for each major kind of Value.
	Kind() int

	// IsTrue returns true iff the Value is non-zero or non-empty.
	IsTrue() bool

	// Less return true iff the Value is less than v. Number < Tuple < Set.
	Less(v Value) bool

	// Negate returns the negation of the Value.
	// - For numbers, this is the arithmetic negation of the value.
	// - For the tuple {(negateTag): x}, it returns x.
	// - For all other values, x, it returns {(negateTag): x}.
	Negate() Value

	// Export converts the Value to a natural Go value.
	Export(context.Context) interface{}
}

// intfValueLess supports
func intfValueLess(a, b interface{}) bool {
	return a.(Value).Less(b.(Value))
}

func exprIsValue(expr Expr) (Value, bool) {
	switch expr := expr.(type) {
	case Value:
		return expr, true
	case LiteralExpr:
		return expr.literal, true
	}
	return nil, false
}

// Attr is a name/Value pair used to construct a Tuple.
type Attr struct {
	Name  string
	Value Value
}

// AttrEnumerator enumerates Values.
type AttrEnumerator interface {
	MoveNext() bool
	Current() (string, Value)
}

// Tuple is a mapping from names to Values.
type Tuple interface {
	Value

	// Access
	Count() int
	Get(name string) (Value, bool)
	MustGet(name string) Value
	HasName(name string) bool
	Names() Names
	Enumerator() AttrEnumerator

	// Transform
	With(name string, value Value) Tuple
	Without(name string) Tuple
	Map(func(Value) (Value, error)) (Tuple, error)
	Project(names Names) Tuple
}

// TupleProjectAllBut returns the projection of t over all of its attributes except
// those specified in names.
func TupleProjectAllBut(t Tuple, names Names) Tuple {
	return t.Project(t.Names().Minus(names))
}

// ValueEnumerator enumerates Values.
type ValueEnumerator interface {
	MoveNext() bool
	Current() Value
}

// OffsetValueEnumerator defines an enumerator that can report the offset of
// each element it enumerates.
type OffsetValueEnumerator interface {
	ValueEnumerator
	Offset() int
}

// Less defines a comparator that returns true iff a < b.
type Less func(a, b Value) bool

// Set represents a Set of Values.
type Set interface {
	Value

	// Access
	Count() int
	Has(Value) bool
	Enumerator() ValueEnumerator

	// Transform
	With(Value) Set
	Without(Value) Set
	Map(func(Value) (Value, error)) (Set, error)
	Where(func(Value) (bool, error)) (Set, error)
	CallAll(context.Context, Value) (Set, error)

	ArrayEnumerator() (OffsetValueEnumerator, bool)
}

// NoReturnError is an error signififying that there was no return value.
type NoReturnError struct {
	input Value
	s     Set
}

func (n NoReturnError) Error() string {
	return fmt.Sprintf("Call: no return values for input %v from set %v", n.input, n.s)
}

// SetCall is a convenience wrapper to call a set and return the result or an
// error if there isn't exactly one result.
func SetCall(ctx context.Context, s Set, arg Value) (Value, error) {
	all, err := s.CallAll(ctx, arg)
	if err != nil {
		return nil, err
	}
	e := all.Enumerator()
	if !e.MoveNext() {
		return nil, NoReturnError{input: arg, s: s}
	}
	result := e.Current()
	if e.MoveNext() {
		return nil, fmt.Errorf("call: too many return values from set %v: %v", s, all)
	}
	return result, nil
}

func mustCallAll(ctx context.Context, s Set, v Value) Value {
	result, err := s.CallAll(ctx, v)
	if err != nil {
		panic(err)
	}
	return result
}

// NewValue constructs a new value from a Go value.
func NewValue(v interface{}) (Value, error) {
	switch x := v.(type) {
	case Value:
		return x, nil
	case bool:
		return NewBool(x), nil
	case uint:
		return NewNumber(float64(x)), nil
	case uint8:
		return NewNumber(float64(x)), nil
	case uint16:
		return NewNumber(float64(x)), nil
	case uintptr:
		return NewNumber(float64(x)), nil
	case uint64:
		return NewNumber(float64(x)), nil
	case int:
		return NewNumber(float64(x)), nil
	case int8:
		return NewNumber(float64(x)), nil
	case int16:
		return NewNumber(float64(x)), nil
	case int32:
		return NewNumber(float64(x)), nil
	case int64:
		return NewNumber(float64(x)), nil
	case float32:
		return NewNumber(float64(x)), nil
	case float64:
		return NewNumber(x), nil
	case string:
		return NewString([]rune(x)), nil
	case []rune:
		return NewString(x), nil
	case []byte:
		return NewBytes(x), nil
	case map[string]interface{}:
		return NewTupleFromMap(x)
	case []interface{}:
		return NewSetFrom(x...)
	default:
		return nil, errors.Errorf("%v (%[1]T) not convertible to Value", v)
	}
}

// AttrEnumeratorToSlice transcribes its Attrs in a slice.
func AttrEnumeratorToSlice(e AttrEnumerator) []Attr {
	attrs := []Attr{}
	for e.MoveNext() {
		name, value := e.Current()
		attrs = append(attrs, Attr{name, value})
	}
	return attrs
}

// AttrEnumeratorToMap transcribes its Attrs in a map.
func AttrEnumeratorToMap(e AttrEnumerator) map[string]Value {
	attrs := map[string]Value{}
	for e.MoveNext() {
		name, value := e.Current()
		attrs[name] = value
	}
	return attrs
}

// ValueEnumeratorToSlice transcribes its Values in a slice.
func ValueEnumeratorToSlice(e ValueEnumerator) []Value {
	values := []Value{}
	for e.MoveNext() {
		values = append(values, e.Current())
	}
	return values
}
