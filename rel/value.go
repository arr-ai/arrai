package rel

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/marcelocantos/frozen"
)

// Expr represents an arr.ai expression.
type Expr interface {
	// Require a String() method.
	fmt.Stringer

	// Eval evaluates the expr in a given scope.
	Eval(local, global *Scope) (Value, error)
}

// Value represents any arr.ai value.
type Value interface {
	frozen.Key

	// Values are Exprs.
	Expr

	// Kind returns a number that is unique for each major kind of Value.
	Kind() int

	// Bool returns true iff the Value is non-zero or non-empty.
	Bool() bool

	// Less return true iff the Value is less than v. Number < Tuple < Set.
	Less(v Value) bool

	// Negate returns the negation of the Value.
	// - For numbers, this is the arithmetic negation of the value.
	// - For the tuple {(negateTag): x}, it returns x.
	// - For all other values, x, it returns {(negateTag): x}.
	Negate() Value

	// Export converts the Value to a natural Go value.
	Export() interface{}
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
	HasName(name string) bool
	Names() Names
	Enumerator() AttrEnumerator

	// Transform
	With(name string, value Value) Tuple
	Without(name string) Tuple
	Project(names Names) Tuple
}

// ValueEnumerator enumerates Values.
type ValueEnumerator interface {
	MoveNext() bool
	Current() Value
}

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
	Map(func(Value) Value) Set
	Where(func(Value) bool) Set
}

// NewValue constructs a new value from a Go value.
func NewValue(v interface{}) (Value, error) {
	switch x := v.(type) {
	case Value:
		return x, nil
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
