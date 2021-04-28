package rel

import (
	"context"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetCall(t *testing.T) {
	t.Parallel()

	foo := func(at int, v Value) Tuple {
		return NewTuple(NewAttr("@", NewNumber(float64(at))), NewAttr("@foo", v))
	}

	set := MustNewSet(
		foo(1, NewNumber(42)),
		foo(1, NewNumber(24)),
	)
	ctx := arraictx.InitRunCtx(context.Background())
	result, err := SetCall(ctx, set, NewNumber(1))
	assert.Error(t, err, "%v", result)
	result, err = SetCall(ctx, set, NewNumber(0))
	assert.Error(t, err, "%v", result)

	set = MustNewSet(
		foo(1, NewNumber(42)),
		foo(2, NewNumber(24)),
	)

	result, err = SetCall(ctx, set, NewNumber(1))
	require.NoError(t, err)
	AssertEqualValues(t, result, NewNumber(42))

	result, err = SetCall(ctx, set, NewNumber(2))
	require.NoError(t, err)
	AssertEqualValues(t, result, NewNumber(24))
}

//nolint:structcheck
func TestNewValue(t *testing.T) {
	// Structs are serialized to tuples.
	type Foo struct {
		num int
		str string
		// Slices without the ordered tag are serialized to arrays.
		arr []int
		// Slices with the unordered tag are serialized to sets.
		set  []int         `unordered:"true"`
		iset []interface{} `unordered:"true"`
		none *Foo
		// All struct field names are serialized to start lowercase.
		CASE     int
		children []*Foo
		// Non-string maps are serialized to dictionaries.
		mixedMap  map[interface{}]interface{}
		stringMap map[string]interface{}
	}

	input := []*Foo{{
		num:  1,
		str:  "a",
		arr:  []int{2, 1},
		set:  []int{2, 1},
		iset: []interface{}{3},
		// Nil values are serialized to empty sets (None).
		none: nil,
		CASE: 0,
		// Unset fields of structs are serialized with default empty values.
		children:  []*Foo{{num: 2}},
		mixedMap:  map[interface{}]interface{}{1: 2, "k": nil},
		stringMap: map[string]interface{}{"a": 1},
	}}

	actual, err := NewValue(input)
	require.NoError(t, err)

	expected, err := NewSet(NewTuple(
		NewIntAttr("num", 1),
		NewStringAttr("str", []rune("a")),
		NewAttr("arr", NewArray(NewNumber(2), NewNumber(1))),
		NewAttr("set", MustNewSet(NewNumber(1), NewNumber(2))),
		NewAttr("iset", MustNewSet(NewNumber(3))),
		NewAttr("none", None),
		NewAttr("cASE", NewNumber(0)),
		NewAttr("mixedMap", MustNewDict(false,
			NewDictEntryTuple(NewNumber(1), NewNumber(2)),
			NewDictEntryTuple(NewString([]rune("k")), None),
		)),
		NewAttr("stringMap", NewTuple(NewIntAttr("a", 1))),
		NewAttr("children", NewArray(NewTuple(
			NewAttr("num", NewNumber(2)),
			NewAttr("str", None),
			NewAttr("arr", None),
			NewAttr("set", None),
			NewAttr("iset", None),
			NewAttr("none", None),
			NewAttr("cASE", NewNumber(0)),
			NewAttr("mixedMap", None),
			NewAttr("stringMap", NewTuple()),
			NewAttr("children", None),
		))),
	))
	require.NoError(t, err)

	AssertEqualValues(t, expected, actual)
}
