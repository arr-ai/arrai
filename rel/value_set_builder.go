package rel

import (
	"reflect"

	"github.com/arr-ai/frozen"
	"github.com/go-errors/errors"
)

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
	buckets map[interface{}][]Value
}

func NewSetBuilder() SetBuilder {
	return SetBuilder{buckets: map[interface{}][]Value{}}
}

func (b *SetBuilder) Add(v Value) {
	t := reflect.TypeOf(v)
	b.buckets[t] = append(b.buckets[t], v)
}

func (b *SetBuilder) Finish() (Set, error) {
	switch len(b.buckets) {
	case 0:
		return None, nil
	case 1:
		for typ, values := range b.buckets {
			switch typ {
			case stringCharTupleType:
				return asString(values...), nil
			case bytesByteTupleType:
				if b, is := asBytes(values...); is {
					return b, nil
				}
				return nil, errors.Errorf("unsupported byte array expr")
			case arrayItemTupleType:
				return asArray(values...), nil
			case dictEntryTupleType:
				tuples := make([]DictEntryTuple, 0, len(values))
				for _, value := range values {
					tuples = append(tuples, value.(DictEntryTuple))
				}
				return NewDict(true, tuples...)
			}
		}
	}
	sb := frozen.SetBuilder{}
	for _, values := range b.buckets {
		for _, value := range values {
			sb.Add(value)
		}
	}
	return GenericSet{sb.Finish()}, nil
}
