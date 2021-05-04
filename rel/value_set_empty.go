package rel

import (
	"context"
	"fmt"
	"reflect"

	"github.com/arr-ai/wbnf/parser"
)

type EmptySet struct{}

var _ Set = EmptySet{}

var emptySetKind = registerKind(199, reflect.TypeOf(EmptySet{}))

func (e EmptySet) Kind() int {
	return emptySetKind
}

func (e EmptySet) IsTrue() bool {
	return false
}

func (e EmptySet) Less(v Value) bool {
	if e == v {
		return false
	}
	switch v.(type) {
	case Number, Tuple:
		return false
	}
	return true
}

func (e EmptySet) Negate() Value {
	return NewTuple(NewAttr(negateTag, EmptySet{}))
}

func (e EmptySet) Export(context.Context) interface{} {
	return nil
}

func (EmptySet) getSetBuilder() setBuilder {
	return newGenericTypeSetBuilder()
}

func (EmptySet) getBucket() fmt.Stringer {
	return genericType
}

func (e EmptySet) Equal(i interface{}) bool {
	_, is := i.(EmptySet)
	return is
}

func (e EmptySet) Hash(seed uintptr) uintptr {
	return seed
}

func (e EmptySet) Eval(ctx context.Context, local Scope) (Value, error) {
	return e, nil
}

func (e EmptySet) Source() parser.Scanner {
	return *parser.NewScanner("")
}

func (e EmptySet) String() string {
	return "{}"
}

func (e EmptySet) Count() int {
	return 0
}

func (e EmptySet) Has(Value) bool {
	return false
}

func (emptyEnumerator) Offset() int {
	panic("wtf")
}

func (e EmptySet) Enumerator() ValueEnumerator {
	return emptyEnumerator{}
}

// func (e EmptySet) ArrayEnumerator() ValueEnumerator {
func (e EmptySet) ArrayEnumerator() ValueEnumerator {
	return emptyEnumerator{}
}

func (e EmptySet) With(v Value) Set {
	return MustNewSet(v)
}

func (e EmptySet) Without(Value) Set {
	return e
}

func (e EmptySet) Map(func(Value) (Value, error)) (Set, error) {
	return e, nil
}

func (e EmptySet) Where(func(Value) (bool, error)) (Set, error) {
	return e, nil
}

func (e EmptySet) CallAll(context.Context, Value, SetBuilder) error {
	return nil
}

func (EmptySet) unionSetSubsetBucket() string {
	return genericType.String()
}
