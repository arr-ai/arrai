package rel

import (
	"context"
	"reflect"
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

//nolint:structcheck,unused,maligned
func TestNewValue(t *testing.T) {
	// Structs are serialized to tuples.
	type Foo struct {
		bool bool
		num  int
		str  string
		// Slices without the ordered tag are serialized to arrays.
		arr []int
		// Slices with the unordered tag are serialized to sets.
		set       []int                  `arrai:",unordered"`
		iset      []interface{}          `arrai:",unordered"`
		omitStr   string                 `arrai:",omitempty"`
		omitInt   int                    `arrai:",omitempty"`
		omitBool  bool                   `arrai:",omitempty"`
		omitSlice []interface{}          `arrai:",omitempty"`
		omitMap   map[string]interface{} `arrai:",omitempty"`
		omitChild *Foo                   `arrai:",omitempty"`
		zeroChild *Foo                   `arrai:",zeroempty"`
		named     string                 `arrai:"newName"`
		none      *Foo
		// All struct field names are serialized to start lowercase.
		CASE     int
		children []*Foo
		// Non-string maps are serialized to dictionaries.
		mixedMap  map[interface{}]interface{}
		stringMap map[string]interface{}
	}

	input := []*Foo{{
		bool: true,
		num:  1,
		str:  "a",
		arr:  []int{2, 1},
		set:  []int{2, 1},
		iset: []interface{}{3},

		// Keep if not empty. This will be kept here, and omitted in children.
		omitStr: "keep",
		// Both empty children and pointers to empty children are considered empty for omission.
		omitInt:   0,
		omitBool:  false,
		omitSlice: []interface{}{},
		omitMap:   map[string]interface{}{},
		omitChild: &Foo{},
		// Empty child structs will be coalesced to empty tuples.
		zeroChild: &Foo{},
		// named will be renamed to newName.
		named: "",
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

	expected := NewArray(NewTuple(
		NewBoolAttr("bool", true),
		NewIntAttr("num", 1),
		NewStringAttr("str", []rune("a")),
		NewAttr("arr", NewArray(NewNumber(2), NewNumber(1))),
		NewAttr("set", MustNewSet(NewNumber(1), NewNumber(2))),
		NewAttr("iset", MustNewSet(NewNumber(3))),
		NewAttr("omitStr", NewString([]rune("keep"))),
		NewAttr("zeroChild", NewTuple()),
		NewAttr("newName", None),
		NewAttr("none", None),
		NewAttr("cASE", NewNumber(0)),
		NewAttr("mixedMap", MustNewDict(false,
			NewDictEntryTuple(NewNumber(1), NewNumber(2)),
			NewDictEntryTuple(NewString([]rune("k")), None),
		)),
		NewAttr("stringMap", NewTuple(NewIntAttr("a", 1))),
		NewAttr("children", NewArray(NewTuple(
			NewBoolAttr("bool", false),
			NewAttr("num", NewNumber(2)),
			NewAttr("str", None),
			NewAttr("arr", None),
			NewAttr("set", None),
			NewAttr("iset", None),
			NewAttr("zeroChild", NewTuple()),
			NewAttr("newName", None),
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

func TestReflectToArray(t *testing.T) {
	t.Parallel()

	arr, err := reflectToArray(reflect.ValueOf([]int{3, 1, 2}))
	require.NoError(t, err)

	assert.Equal(t, NewArray(NewNumber(3), NewNumber(1), NewNumber(2)), arr)
}

func TestArraiTags(t *testing.T) {
	t.Parallel()

	type Foo struct {
		X int `arrai:"foo , omitempty, baz"`
	}
	f := reflect.ValueOf(Foo{}).Type().Field(0)

	assert.Equal(t, []string{"foo", "omitempty", "baz"}, arraiTags(f))
}

func TestArraiTagMap(t *testing.T) {
	t.Parallel()

	type Foo struct {
		X int `arrai:"foo , omitempty, baz"`
	}
	f := reflect.ValueOf(Foo{}).Type().Field(0)

	assert.Equal(t, map[string]bool{"foo": true, "omitempty": true, "baz": true}, arraiTagMap(f))
}
