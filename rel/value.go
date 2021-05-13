package rel

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/iancoleman/strcase"

	"github.com/arr-ai/frozen"
	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

const (
	arraiTag  = "arrai"
	unordered = "unordered"
	omitempty = "omitempty"
	zeroempty = "zeroempty"
)

// Expr represents an arr.ai expression.
type Expr interface {
	// All exprs can be serialized to strings with the String() method.
	fmt.Stringer

	// Eval evaluates the expr in a given scope.
	Eval(ctx context.Context, local Scope) (Value, error)

	// Source returns the Scanner that locates the expression in a source file.
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

	// functions for building sets
	getSetBuilder() setBuilder
	getBucket() fmt.Stringer
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

// Less defines a comparator that returns true iff a < b.
type Less func(a, b Value) bool

// Set represents a Set of Values.
type Set interface {
	Value

	// Access
	Count() int
	Has(Value) bool
	Enumerator() ValueEnumerator
	ArrayEnumerator() ValueEnumerator // iterates in ascending order.

	// Transform
	With(Value) Set
	Without(Value) Set
	Map(func(Value) (Value, error)) (Set, error)
	Where(func(Value) (bool, error)) (Set, error)
	CallAll(context.Context, Value, SetBuilder) error

	unionSetSubsetBucket() string
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
	b := NewSetBuilder()
	err := s.CallAll(ctx, arg, b)
	if err != nil {
		return nil, err
	}
	all, err := b.Finish()
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
	b := NewSetBuilder()
	err := s.CallAll(ctx, v, b)
	if err != nil {
		panic(err)
	}
	result, err := b.Finish()
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
	default:
		// Fall back on reflection for custom types.
		return reflectNewValue(reflect.ValueOf(x))
	}
}

// reflectNewValue uses reflection to inspect the type of x and unpack its values.
func reflectNewValue(x reflect.Value) (Value, error) {
	if !x.IsValid() {
		return None, nil
	}
	t := x.Type()
	switch t.Kind() {
	case reflect.Ptr:
		if x.IsNil() {
			return None, nil
		}
		return NewValue(x.Elem().Interface())
	case reflect.Array, reflect.Slice:
		return reflectToArray(x)
	case reflect.Map:
		entries := make([]DictEntryTuple, 0, x.Len())
		for _, k := range x.MapKeys() {
			v := x.MapIndex(k)
			kv, err := NewValue(k.Interface())
			if err != nil {
				return nil, err
			}
			vv, err := NewValue(v.Interface())
			if err != nil {
				return nil, err
			}
			entries = append(entries, NewDictEntryTuple(kv, vv))
		}
		return NewDict(false, entries...)
	case reflect.Struct:
		s := map[string]interface{}{}

		// Ensure x is accessible.
		xv := reflect.New(t).Elem()
		xv.Set(x)

		for i := 0; i < t.NumField(); i++ {
			tf := t.Field(i)
			// Ensure each field of x is accessible.
			f := xv.Field(i)
			f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()

			var v Value
			var err error
			tags := arraiTagMap(tf)
			if _, ok := tags[omitempty]; ok && isEmpty(f) {
				continue
			} else if tags[zeroempty] && isEmpty(f) {
				v = zeroValue(f.Type())
			} else {
				switch f.Type().Kind() {
				case reflect.Array, reflect.Slice:
					v, err = reflectToValues(f, tags[unordered])
				default:
					v, err = NewValue(f.Interface())
				}
			}
			if err != nil {
				return nil, err
			}
			// Lowercase the first character in case it's uppercase only for Go exporting.
			// TODO: Handle a name tag to override behaviour.
			s[fieldName(tf)] = v
		}
		return NewTupleFromMap(s)
	default:
		return nil, errors.Errorf("%v (%[1]T) not convertible to Value", x)
	}
}

// zeroValue returns the arr.ai value that best represents an empty value of a Go value's type.
func zeroValue(v reflect.Type) Value {
	switch v.Kind() {
	case reflect.Ptr:
		return zeroValue(v.Elem())
	case reflect.Struct, reflect.Map:
		return NewTuple()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return NewNumber(0)
	}
	return None
}

// isEmpty returns true if a value is considered empty for serialization purposes (i.e. should be
// omitted if the omitempty tag is specified).
//
// This is true for zero values, empty collections, structs with only empty fields, and pointers to
// any of those.
func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr:
		return v.IsNil() || isEmpty(v.Elem())
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !isEmpty(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		return v.IsZero()
	}
}

// fieldName returns the name that should be used for a given field.
//
// It will return the name if specified in a tag (e.g. `arrai:"name"`), or else default to the name
// of the field in the struct in lowerCamelCase.
func fieldName(sf reflect.StructField) string {
	if name := arraiTags(sf)[0]; name != "" {
		return name
	}
	return strcase.ToLowerCamel(sf.Name)
}

// arraiTags returns the values associated with the `arrai:` tag, separated on commas and trimmed.
func arraiTags(sf reflect.StructField) []string {
	tags := strings.Split(sf.Tag.Get(arraiTag), ",")
	for i, t := range tags {
		tags[i] = strings.TrimSpace(t)
	}
	return tags
}

// arraiTags returns a map with a key for each values associated with the `arrai:` tag.
func arraiTagMap(sf reflect.StructField) map[string]bool {
	tags := map[string]bool{}
	for _, t := range arraiTags(sf) {
		tags[t] = true
	}
	return tags
}

// reflectToValues assumed x is a slice or array, and returns x serialized to a collection of Values.
//
// If ordered is true, the result will be an Array. If false, it will be a Set.
// If x is not a slice or array, reflectToValues will panic.
func reflectToValues(x reflect.Value, unordered bool) (Value, error) {
	vs := make([]Value, 0, x.Len())
	for i := 0; i < x.Len(); i++ {
		v, err := NewValue(x.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	if unordered {
		return NewSet(vs...)
	}
	return NewArray(vs...), nil
}

func reflectToArray(x reflect.Value) (Value, error) {
	return reflectToValues(x, false)
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
