package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewArrayItemTuple(t *testing.T) {
	t.Parallel()
	type args struct {
		at   int
		item Value
	}
	tests := []struct {
		name string
		args args
		want ArrayItemTuple
	}{
		{name: "default", want: ArrayItemTuple{}},
		{name: "0@0", want: ArrayItemTuple{at: 0, item: NewNumber(0)}, args: args{at: 0, item: NewNumber(0)}},
		{name: "0@1", want: ArrayItemTuple{at: 1, item: NewNumber(0)}, args: args{at: 1, item: NewNumber(0)}},
		{name: "a@0", want: ArrayItemTuple{at: 0, item: None}, args: args{at: 0, item: None}},
		{name: "a@1", want: ArrayItemTuple{at: 1, item: None}, args: args{at: 1, item: None}},
		{name: "a@-1", want: ArrayItemTuple{at: -1, item: None}, args: args{at: -1, item: None}},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			assert.EqualValues(t, test.want, NewArrayItemTuple(test.args.at, test.args.item))
		})
	}
}

func Test_newArrayItemTupleFromTuple(t *testing.T) {
	t.Parallel()
	type args struct {
		t Tuple
	}
	tests := []struct {
		name string
		args args
		want ArrayItemTuple
		ok   bool
	}{
		{name: "0@0",
			want: NewArrayItemTuple(0, NewNumber(0)),
			ok:   true,
			args: args{t: NewTuple(NewAttr("@", NewNumber(0)), NewAttr(ArrayItemAttr, NewNumber(0)))},
		},
		{name: "{}@0",
			want: NewArrayItemTuple(0, None),
			ok:   true,
			args: args{t: NewTuple(NewAttr("@", NewNumber(0)), NewAttr(ArrayItemAttr, None))},
		},
		{name: "{}@42",
			want: NewArrayItemTuple(42, None),
			ok:   true,
			args: args{t: NewTuple(NewAttr("@", NewNumber(42)), NewAttr(ArrayItemAttr, None))},
		},
		{name: "no-@",
			want: ArrayItemTuple{},
			ok:   false,
			args: args{t: NewTuple(NewAttr("at", NewNumber(0)), NewAttr(ArrayItemAttr, None))},
		},
		{name: "no-ArrayItemAttr",
			want: ArrayItemTuple{},
			ok:   false,
			args: args{t: NewTuple(NewAttr("@", NewNumber(0)), NewAttr("item", NewNumber(0)))},
		},
	}
	for _, test := range tests { //nolint:dupl
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, ok := newArrayItemTupleFromTuple(test.args.t)
			got2 := maybeNewArrayItemTupleFromTuple(test.args.t)
			if test.ok {
				assert.True(t, ok)
				assert.EqualValues(t, test.want, got)
				assert.EqualValues(t, test.want, got2)
			} else {
				assert.False(t, ok)
				assert.True(t, test.args.t.Equal(got2))
			}
		})
	}
}

func TestArrayItemTuple_Equal(t *testing.T) {
	t.Parallel()
	values := []Value{
		None,
		NewNumber(42),
		NewArrayItemTuple(0, NewNumber(0)),
		NewArrayItemTuple(0, None),
		NewArrayItemTuple(42, NewNumber(0)),
		NewArrayItemTuple(42, None),
		NewArrayItemTuple(42, NewNumber(1)),
		NewArrayItemTuple(43, NewNumber(1)),
	}
	for i, x := range values {
		for j, y := range values {
			assert.Equal(t, i == j, x.Equal(y), "values[%d]=%v, values[%d]=%v", i, x, j, y)
		}
	}
}

func TestArrayItemTuple_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		tuple ArrayItemTuple
		want  string
	}{
		{name: "0@0", want: "(@: 0, @item: 0)", tuple: NewArrayItemTuple(0, NewNumber(0))},
		{name: "a@0", want: "(@: 0, @item: {})", tuple: NewArrayItemTuple(0, None)},
		{name: "0@42", want: "(@: 42, @item: 0)", tuple: NewArrayItemTuple(42, NewNumber(0))},
		{name: "a@42", want: "(@: 42, @item: {})", tuple: NewArrayItemTuple(42, None)},
		{name: "b@42", want: "(@: 42, @item: 1)", tuple: NewArrayItemTuple(42, NewNumber(1))},
		{name: "b@43", want: "(@: 43, @item: 1)", tuple: NewArrayItemTuple(43, NewNumber(1))},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.tuple.String())
		})
	}
}

func TestArrayItemTuple_Eval(t *testing.T) {
	t.Parallel()
	tuple := NewArrayItemTuple(42, None)
	value, err := tuple.Eval(EmptyScope)
	require.NoError(t, err)
	assert.Equal(t, tuple, value)
}

func TestArrayItemTuple_IsTrue(t *testing.T) {
	t.Parallel()
	values := []Value{
		NewArrayItemTuple(0, NewNumber(0)),
		NewArrayItemTuple(0, None),
		NewArrayItemTuple(42, NewNumber(0)),
		NewArrayItemTuple(42, NewNumber('😜')),
	}
	for i, x := range values {
		assert.True(t, x.IsTrue(), "values[%d]=%v", i, x)
	}
}

func TestArrayItemTuple_Less(t *testing.T) {
	t.Parallel()
	assert.True(t, !NewArrayItemTuple(0, NewNumber(0)).Less(NewArrayItemTuple(0, NewNumber(0))))
	assert.True(t, NewArrayItemTuple(0, NewNumber(0)).Less(NewArrayItemTuple(0, None)))
	assert.True(t, NewArrayItemTuple(0, None).Less(NewArrayItemTuple(42, NewNumber(0))))
}

func TestArrayItemTuple_Negate(t *testing.T) {
	t.Parallel()
	assert.Equal(t, NewArrayItemTuple(0, NewNumber(0)), NewArrayItemTuple(0, NewNumber(0)).Negate())
	assert.Equal(t, NewArrayItemTuple(-1, None.Negate()), NewArrayItemTuple(1, None).Negate())
}

func TestArrayItemTuple_Export(t *testing.T) {
	t.Parallel()
	assert.Equal(t, map[string]interface{}{"@": 1, "@item": []interface{}{}}, NewArrayItemTuple(1, None).Export())
}

func TestArrayItemTuple_Count(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 2, NewArrayItemTuple(0, NewNumber(0)).Count())
	assert.Equal(t, 2, NewArrayItemTuple(1, None).Count())
}

func TestArrayItemTuple_Get(t *testing.T) {
	t.Parallel()

	assertGet := func(tuple Tuple, attr string, value Value) {
		v, has := tuple.Get(attr)
		if assert.True(t, has) {
			assert.Equal(t, value, v)
		}
	}
	assertGet(NewArrayItemTuple(1, None), "@", NewNumber(1))
	assertGet(NewArrayItemTuple(1, None), "@item", None)

	assertNotGet := func(tuple Tuple, attr string) bool {
		v, has := tuple.Get(attr)
		return assert.False(t, has, "%q => %v", attr, v)
	}
	assertNotGet(NewArrayItemTuple(1, None), "@@")
}

func TestArrayItemTuple_MustGet(t *testing.T) {
	t.Parallel()
	assert.Equal(t, NewNumber(1), NewArrayItemTuple(1, None).MustGet("@"))
	assert.Equal(t, None, NewArrayItemTuple(1, None).MustGet("@item"))
	assert.Panics(t, func() { NewArrayItemTuple(1, None).MustGet("item") })
}

func TestArrayItemTuple_With(t *testing.T) {
	t.Parallel()
	assert.True(t, NewArrayItemTuple(1, None).With("@", NewNumber(42)).Equal(NewArrayItemTuple(42, None)))
	assert.True(t, NewArrayItemTuple(1, None).With("@item", NewNumber(2)).Equal(NewArrayItemTuple(1, NewNumber(2))))
	assert.True(t, NewArrayItemTuple(1, None).With("x", NewNumber(2)).Equal(NewTuple(
		NewAttr("@", NewNumber(1)),
		NewAttr("@item", None),
		NewAttr("x", NewNumber(2)),
	)))
}

func TestArrayItemTuple_Without(t *testing.T) {
	t.Parallel()
	assert.True(t, NewArrayItemTuple(1, None).Without("@").Equal(NewTuple(NewAttr("@item", None))))
	assert.True(t, NewArrayItemTuple(1, None).Without("@item").Equal(NewTuple(NewAttr("@", NewNumber(1)))))
	assert.True(t, NewArrayItemTuple(1, None).Without("x").Equal(NewArrayItemTuple(1, None)))
}

func TestArrayItemTuple_Map(t *testing.T) {
	t.Parallel()
	assert.True(t, NewArrayItemTuple(1, NewNumber(10)).Map(func(v Value) Value {
		return NewNumber(v.(Number).Float64() + 1)
	}).Equal(NewArrayItemTuple(2, NewNumber(11))))
}

func TestArrayItemTuple_HasName(t *testing.T) {
	t.Parallel()
	assert.True(t, NewArrayItemTuple(1, None).HasName("@"))
	assert.True(t, NewArrayItemTuple(1, None).HasName("@item"))
	assert.False(t, NewArrayItemTuple(1, None).HasName("item"))
}

func TestArrayItemTuple_Attributes(t *testing.T) {
	t.Parallel()
	assert.Equal(t,
		map[string]Value{"@": NewNumber(1), "@item": None},
		NewArrayItemTuple(1, None).Attributes())
}

func TestArrayItemTuple_Names(t *testing.T) {
	t.Parallel()
	assert.True(t, NewArrayItemTuple(1, None).Names().Equal(NewNames("@", "@item")))
}

func TestArrayItemTuple_Project(t *testing.T) {
	t.Parallel()
	assert.True(t, NewArrayItemTuple(1, None).Project(NewNames("@")).Equal(NewTuple(NewAttr("@", NewNumber(1)))))
	assert.True(t, NewArrayItemTuple(1, None).Project(NewNames("@item")).Equal(NewTuple(NewAttr("@item", None))))
	assert.Nil(t, NewArrayItemTuple(1, None).Project(NewNames("x")))
}

func TestArrayItemTuple_Enumerator(t *testing.T) {
	t.Parallel()
	attrs := map[string]int{}
	for e := NewArrayItemTuple(1, NewNumber(2)).Enumerator(); e.MoveNext(); {
		name, value := e.Current()
		v, has := attrs[name]
		require.False(t, has, "%q => %v %v", name, v, value)
		attrs[name] = int(value.(Number).Float64())
	}
	assert.Equal(t, map[string]int{"@": 1, "@item": 2}, attrs)
}
