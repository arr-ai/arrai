package rel

import (
	"fmt"
	"reflect"

	"github.com/arr-ai/frozen"
)

// TODO: Replace genericSetBuilder with custom setBuilders.
var (
	genericSetFinish = func(values ...Value) (Set, error) {
		sb := frozen.SetBuilder{}
		for _, v := range values {
			sb.Add(v)
		}
		return GenericSet{sb.Finish()}, nil
	}

	stringFinish = func(values ...Value) (Set, error) {
		return asString(values...), nil
	}

	bytesFinish = func(values ...Value) (Set, error) {
		return asBytes(values...), nil
	}

	arrayFinish = func(values ...Value) (Set, error) {
		return asArray(values...), nil
	}

	dictFinish = func(values ...Value) (Set, error) {
		v := make([]DictEntryTuple, 0, len(values))
		for _, value := range values {
			v = append(v, value.(DictEntryTuple))
		}
		return NewDict(true, v...)
	}
)

type setBuilder interface {
	Add(v Value)
	Finish() (Set, error)
}

type genericSetBuilder struct {
	values []Value
	finish func(values ...Value) (Set, error)
}

func (b *genericSetBuilder) Add(v Value) {
	b.values = append(b.values, v)
}

func (b *genericSetBuilder) Finish() (Set, error) {
	return b.finish(b.values...)
}

// MustNewSet constructs a genericSet from a set of Values, or panics if construction fails.
func MustNewSet(values ...Value) Set {
	s, err := NewSet(values...)
	if err != nil {
		panic(err)
	}
	return s
}

// NewSet constructs a genericSet from a set of Values.
func NewSet(values ...Value) (Set, error) {
	b := NewSetBuilder()
	for _, v := range values {
		b.Add(v)
	}
	return b.Finish()
}

// NewSetFrom constructs a genericSet from interfaces.
func NewSetFrom(intfs ...interface{}) (Set, error) {
	b := NewSetBuilder()
	for _, intf := range intfs {
		value, err := NewValue(intf)
		if err != nil {
			return nil, err
		}
		b.Add(value)
	}
	return b.Finish()
}

type SetBuilder struct {
	buckets map[fmt.Stringer]setBuilder
}

func NewSetBuilder() SetBuilder {
	return SetBuilder{buckets: map[fmt.Stringer]setBuilder{}}
}

type generic struct{}

var genericType = reflect.TypeOf(generic{})

func newGenericTypeSetBuilder() setBuilder {
	return &genericSetBuilder{values: []Value{}, finish: genericSetFinish}
}

func (b *SetBuilder) Add(v Value) {
	builder, has := b.buckets[v.getBucket()]
	if !has {
		builder = v.getSetBuilder()
		b.buckets[v.getBucket()] = builder
	}
	builder.Add(v)
}

func (b *SetBuilder) Finish() (Set, error) {
	switch len(b.buckets) {
	case 0:
		return None, nil
	case 1:
		for _, values := range b.buckets {
			return values.Finish()
		}
	}

	var mb frozen.StringMapBuilder
	for k, v := range b.buckets {
		set, err := v.Finish()
		if err != nil {
			return nil, err
		}
		mb.Put(k.String(), set)
	}
	return UnionSet{mb.Finish()}, nil
}
