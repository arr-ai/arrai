package rel

import (
	"context"
	"fmt"
	"reflect"

	"github.com/arr-ai/frozen"
	"github.com/arr-ai/hash"
	"github.com/arr-ai/wbnf/parser"
)

// TrueSet is a special set that represents True/{()}.
type TrueSet struct{}

var _ Set = TrueSet{}

var trueSetKind = registerKind(199, reflect.TypeOf(TrueSet{}))

func (TrueSet) Kind() int {
	return trueSetKind
}

func (TrueSet) IsTrue() bool {
	return true
}

func (t TrueSet) Less(v Value) bool {
	switch v.(type) {
	case TrueSet, Number, Tuple, EmptySet:
		return false
	}
	return true
}

func (t TrueSet) Negate() Value {
	return NewTuple(NewAttr(negateTag, t))
}

func (TrueSet) Export(context.Context) interface{} {
	return true
}

func (TrueSet) getSetBuilder() setBuilder {
	return newGenericTypeSetBuilder()
}

func (TrueSet) getBucket() fmt.Stringer {
	return genericType
}

func (TrueSet) Equal(i interface{}) bool {
	_, is := i.(TrueSet)
	return is
}

func (TrueSet) Hash(seed uintptr) uintptr {
	return seed ^ hash.Interface(EmptyTuple, 0)
}

func (t TrueSet) Eval(ctx context.Context, local Scope) (Value, error) {
	return t, nil
}

func (TrueSet) Source() parser.Scanner {
	return *parser.NewScanner("")
}

func (TrueSet) String() string {
	return sTrue
}

func (TrueSet) Count() int {
	return 1
}

func (TrueSet) Has(v Value) bool {
	return v.Equal(EmptyTuple)
}

func (TrueSet) Enumerator() ValueEnumerator {
	return &genericSetEnumerator{frozen.NewSet(EmptyTuple).Range()}
}

func (t TrueSet) ArrayEnumerator() ValueEnumerator {
	return t.Enumerator()
}

func (t TrueSet) With(v Value) Set {
	if v.Equal(EmptyTuple) {
		return t
	}
	return MustNewSet(EmptyTuple, v)
}

func (t TrueSet) Without(v Value) Set {
	if v.Equal(EmptyTuple) {
		return None
	}
	return t
}

func (TrueSet) Map(f func(Value) (Value, error)) (Set, error) {
	v, err := f(EmptyTuple)
	if err != nil {
		return nil, err
	}
	return NewSet(v)
}

func (t TrueSet) Where(p func(Value) (bool, error)) (Set, error) {
	ok, err := p(EmptyTuple)
	if err != nil {
		return nil, err
	}
	if !ok {
		return None, nil
	}
	return t, nil
}

func (TrueSet) CallAll(context.Context, Value, SetBuilder) error {
	// TODO: what to do here
	return nil
}

func (TrueSet) unionSetSubsetBucket() string {
	return genericType.String()
}
