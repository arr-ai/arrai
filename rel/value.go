package rel

import (
	"fmt"
	"github.com/arr-ai/frozen"
	"github.com/arr-ai/wbnf/parser"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"reflect"
)

// Expr represents an arr.ai expression.
type Expr interface {
	// Require a String() method.
	fmt.Stringer

	// Eval evaluates the expr in a given scope.
	Eval(local Scope) (Value, error)

	Source() parser.Scanner
}

// Value represents any arr.ai value.
type Value interface {
	frozen.Key

	// Values are Exprs.
	Expr

	// Kind returns a number that is unique for each major kind of Value.
	Kind() int

	// Bool returns true iff the Value is non-zero or non-empty.
	IsTrue() bool

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
	Map(func(Value) Value) Tuple
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

type OffsetValueEnumerator interface {
	ValueEnumerator
	Offset() int
}

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
	Map(func(Value) Value) Set
	Where(func(Value) (bool, error)) (Set, error)
	CallAll(Value) (Set, error)

	ArrayEnumerator() (OffsetValueEnumerator, bool)
}

type NoReturnError struct {
	input Value
	s     Set
}

func (n NoReturnError) Error() string {
	return fmt.Sprintf("Call: no return values for input %v from set %v", n.input, n.s)
}

func SetCall(s Set, arg Value) (Value, error) {
	result, err := s.CallAll(arg)
	if err != nil {
		return nil, err
	}
	if !result.IsTrue() {
		return nil, NoReturnError{input: arg, s: s}
	}
	for i, e := 1, result.Enumerator(); e.MoveNext(); i++ {
		if i > 1 {
			return nil, fmt.Errorf("call: too many return values from set %v: %v", s, result)
		}
	}
	return SetAny(result), nil
}

func MustCallAll(s Set, v Value) Value {
	result, err := s.CallAll(v)
	if err != nil {
		panic(err)
	}
	return result
}

func SetAny(s Set) Value {
	for e := s.Enumerator(); e.MoveNext(); {
		return e.Current()
	}
	panic("SetAny: set is empty")
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
		v, err := newReflectValueValue(reflect.ValueOf(x), make(visited))
		if err != nil {
			return nil, err
		} else {
			return v, nil
		}
	}
}

// visited tracks pointers of ancestors of a reflected value to detect cycles.
type visited map[uintptr]bool

func newReflectValueValue(v reflect.Value, p visited) (Value, error) {
	if v.Kind() == reflect.Ptr {
		if _, cycle := p[v.Pointer()]; cycle {
			logrus.Tracef("cycle in reflection (path len %d)", len(p))
			return nil, nil
		}
		p[v.Pointer()] = true
		defer func() { delete(p, v.Pointer()) }()
	}

	switch v.Kind() {
	case reflect.Bool:
		return NewBool(v.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr:
		return NewNumber(float64(v.Int())), nil
	case reflect.Float32, reflect.Float64:
		return NewNumber(v.Float()), nil
	case reflect.Complex64, reflect.Complex128:
		return nil, fmt.Errorf("%T not yet supported as arr.ai value", v.Interface())
	case reflect.Array, reflect.Slice:
		n := v.Len()
		items := make([]Value, 0, n)
		for i := 0; i < n; i++ {
			item, err := newReflectValueValue(v.Index(i), p)
			if err != nil {
				return nil, err
			} else if item != nil {
				items = append(items, item)
			}
		}
		return NewArray(items...), nil
	case reflect.Func:
		return nil, fmt.Errorf("function wrapping not yet supported")
	case reflect.Interface, reflect.Ptr:
		if !v.IsValid() || !v.Elem().IsValid() {
			return None, nil
		}
		return newReflectValueValue(v.Elem(), p)
	case reflect.Map:
		kvs := make([]DictEntryTuple, 0, v.Len())
		for r := v.MapRange(); r.Next(); {
			at, err := newReflectValueValue(r.Key(), p)
			if err != nil {
				return nil, err
			}
			value, err := newReflectValueValue(r.Value(), p)
			if err != nil {
				return nil, err
			} else if value != nil {
				kvs = append(kvs, NewDictEntryTuple(at, value))
			}
		}
		return NewDict(false, kvs...), nil
	case reflect.String:
		return NewString([]rune(v.String())), nil
	case reflect.Struct:
		n := v.NumField()
		attrs := make([]Attr, 0, n+1)
		attrs = append(attrs, NewStringAttr("@type", []rune(v.Type().Name())))
		t := v.Type()
		for i := 0; i < n; i++ {
			var name string
			if tag, has := t.Field(i).Tag.Lookup("arrai"); has {
				name = tag
			} else {
				name = t.Field(i).Name
			}
			value, err := newReflectValueValue(v.Field(i), p)
			if err != nil {
				return nil, err
			} else if value != nil {
				attrs = append(attrs, NewAttr(name, value))
			}
		}
		return NewTuple(attrs...), nil
	}
	return nil, errors.Errorf("%v (%[1]T) not convertible to Value", v.Interface())
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
