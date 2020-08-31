package rel //nolint:dupl

import (
	"context"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStringCharTuple(t *testing.T) {
	t.Parallel()
	type args struct {
		at   int
		char rune
	}
	tests := []struct {
		name string
		args args
		want StringCharTuple
	}{
		{name: "default", want: StringCharTuple{}},
		{name: "0@0", want: StringCharTuple{}, args: args{at: 0, char: 0}},
		{name: "0@1", want: StringCharTuple{at: 1, char: 0}, args: args{at: 1, char: 0}},
		{name: "a@0", want: StringCharTuple{at: 0, char: 'a'}, args: args{at: 0, char: 'a'}},
		{name: "a@1", want: StringCharTuple{at: 1, char: 'a'}, args: args{at: 1, char: 'a'}},
		{name: "a@-1", want: StringCharTuple{at: -1, char: 'a'}, args: args{at: -1, char: 'a'}},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			assert.EqualValues(t, test.want, NewStringCharTuple(test.args.at, test.args.char))
		})
	}
}

func Test_newCharTupleFromTuple(t *testing.T) {
	t.Parallel()
	type args struct {
		t Tuple
	}
	tests := []struct {
		name string
		args args
		want StringCharTuple
		ok   bool
	}{
		{
			name: "0@0",
			want: NewStringCharTuple(0, 0),
			ok:   true,
			args: args{t: NewTuple(NewIntAttr("@", 0), NewIntAttr(StringCharAttr, 0))},
		},
		{
			name: "a@0",
			want: NewStringCharTuple(0, 'a'),
			ok:   true,
			args: args{t: NewTuple(NewIntAttr("@", 0), NewIntAttr(StringCharAttr, 'a'))},
		},
		{
			name: "a@42",
			want: NewStringCharTuple(42, 'a'),
			ok:   true,
			args: args{t: NewTuple(NewIntAttr("@", 42), NewIntAttr(StringCharAttr, 'a'))},
		},
		{
			name: "no-@",
			want: StringCharTuple{},
			ok:   false,
			args: args{t: NewTuple(NewIntAttr("at", 0), NewIntAttr(StringCharAttr, 'a'))},
		},
		{
			name: "no-CharAttr",
			want: StringCharTuple{},
			ok:   false,
			args: args{t: NewTuple(NewIntAttr("@", 0), NewIntAttr("char", 'a'))},
		},
	}
	for _, test := range tests { //nolint:dupl
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, ok := newCharTupleFromTuple(test.args.t)
			got2 := maybeNewCharTupleFromTuple(test.args.t)
			if test.ok {
				assert.True(t, ok)
				AssertEqualValues(t, test.want, got)
				AssertEqualValues(t, test.want, got2)
			} else {
				assert.False(t, ok)
				AssertEqualValues(t, test.args.t, got2)
			}
		})
	}
}

func TestStringCharTuple_Equal(t *testing.T) {
	t.Parallel()
	values := []Value{
		None,
		NewNumber(42),
		NewStringCharTuple(0, 0),
		NewStringCharTuple(0, 'a'),
		NewStringCharTuple(42, 0),
		NewStringCharTuple(42, 'a'),
		NewStringCharTuple(42, 'b'),
		NewStringCharTuple(43, 'b'),
	}
	for i, x := range values {
		for j, y := range values {
			assert.Equal(t, i == j, x.Equal(y), "values[%d]=%v, values[%d]=%v", i, x, j, y)
		}
	}
}

func TestStringCharTuple_String(t *testing.T) { //nolint:dupl
	t.Parallel()
	tests := []struct {
		name  string
		tuple StringCharTuple
		want  string
	}{
		{name: "0@0", want: "(@: 0, @char: 0)", tuple: NewStringCharTuple(0, 0)},
		{name: "a@0", want: "(@: 0, @char: 97)", tuple: NewStringCharTuple(0, 'a')},
		{name: "0@42", want: "(@: 42, @char: 0)", tuple: NewStringCharTuple(42, 0)},
		{name: "a@42", want: "(@: 42, @char: 97)", tuple: NewStringCharTuple(42, 'a')},
		{name: "b@42", want: "(@: 42, @char: 98)", tuple: NewStringCharTuple(42, 'b')},
		{name: "b@43", want: "(@: 43, @char: 98)", tuple: NewStringCharTuple(43, 'b')},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.tuple.String())
		})
	}
}

func TestStringCharTuple_Eval(t *testing.T) {
	t.Parallel()
	tuple := NewStringCharTuple(42, 'a')
	value, err := tuple.Eval(arraictx.InitRunCtx(context.Background()), EmptyScope)
	require.NoError(t, err)
	assert.Equal(t, tuple, value)
}

func TestStringCharTuple_IsTrue(t *testing.T) {
	t.Parallel()
	values := []Value{
		NewStringCharTuple(0, 0),
		NewStringCharTuple(0, 'a'),
		NewStringCharTuple(42, 0),
		NewStringCharTuple(42, '😜'),
	}
	for i, x := range values {
		assert.True(t, x.IsTrue(), "values[%d]=%v", i, x)
	}
}

func TestStringCharTuple_Less(t *testing.T) {
	t.Parallel()
	assert.True(t, !NewStringCharTuple(0, 0).Less(NewStringCharTuple(0, 0)))
	assert.True(t, NewStringCharTuple(0, 0).Less(NewStringCharTuple(0, 'a')))
	assert.True(t, NewStringCharTuple(0, 'a').Less(NewStringCharTuple(42, 0)))
}

func TestStringCharTuple_Negate(t *testing.T) {
	t.Parallel()
	assert.Equal(t, NewStringCharTuple(0, 0), NewStringCharTuple(0, 0).Negate())
	assert.Equal(t, NewStringCharTuple(-1, -'a'), NewStringCharTuple(1, 'a').Negate())
}

func TestStringCharTuple_Export(t *testing.T) {
	t.Parallel()
	assert.Equal(t,
		map[string]interface{}{"@": 1, "@char": 'a'},
		NewStringCharTuple(1, 'a').Export(arraictx.InitRunCtx(context.Background())),
	)
}

func TestStringCharTuple_Count(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 2, NewStringCharTuple(0, 0).Count())
	assert.Equal(t, 2, NewStringCharTuple(1, 'a').Count())
}

func TestStringCharTuple_Get(t *testing.T) { //nolint:dupl
	t.Parallel()

	assertGet := func(tuple Tuple, attr string, value Value) {
		v, has := tuple.Get(attr)
		if assert.True(t, has) {
			assert.Equal(t, value, v)
		}
	}
	assertGet(NewStringCharTuple(1, 'a'), "@", NewNumber(1))
	assertGet(NewStringCharTuple(1, 'a'), "@char", NewNumber('a'))

	assertNotGet := func(tuple Tuple, attr string) bool {
		v, has := tuple.Get(attr)
		return assert.False(t, has, "%q => %v", attr, v)
	}
	assertNotGet(NewStringCharTuple(1, 'a'), "@@")
}

func TestStringCharTuple_MustGet(t *testing.T) {
	t.Parallel()
	assert.Equal(t, NewNumber(1), NewStringCharTuple(1, 'a').MustGet("@"))
	assert.Equal(t, NewNumber('a'), NewStringCharTuple(1, 'a').MustGet("@char"))
	assert.Panics(t, func() { NewStringCharTuple(1, 'a').MustGet("char") })
}

func TestStringCharTuple_With(t *testing.T) {
	t.Parallel()
	AssertEqualValues(t, NewStringCharTuple(42, 'a'), NewStringCharTuple(1, 'a').With("@", NewNumber(42)))
	AssertEqualValues(t, NewStringCharTuple(1, 'b'), NewStringCharTuple(1, 'a').With("@char", NewNumber('b')))
	AssertEqualValues(t, NewTuple(
		NewIntAttr("@", 1),
		NewIntAttr("@char", 'a'),
		NewIntAttr("x", 'b'),
	), NewStringCharTuple(1, 'a').With("x", NewNumber('b')))
}

func TestStringCharTuple_Without(t *testing.T) {
	t.Parallel()
	AssertEqualValues(t, NewTuple(NewIntAttr("@char", 'a')), NewStringCharTuple(1, 'a').Without("@"))
	AssertEqualValues(t, NewTuple(NewIntAttr("@", 1)), NewStringCharTuple(1, 'a').Without("@char"))
	AssertEqualValues(t, NewStringCharTuple(1, 'a'), NewStringCharTuple(1, 'a').Without("x"))
}

func TestStringCharTuple_Map(t *testing.T) {
	t.Parallel()
	m, err := NewStringCharTuple(1, 'a').Map(func(v Value) (Value, error) {
		return NewNumber(v.(Number).Float64() + 1), nil
	})
	require.NoError(t, err)
	AssertEqualValues(t,
		NewStringCharTuple(2, 'b'),
		m,
	)
}

func TestStringCharTuple_HasName(t *testing.T) {
	t.Parallel()
	assert.True(t, NewStringCharTuple(1, 'a').HasName("@"))
	assert.True(t, NewStringCharTuple(1, 'a').HasName("@char"))
	assert.False(t, NewStringCharTuple(1, 'a').HasName("@item"))
}

func TestStringCharTuple_Attributes(t *testing.T) {
	t.Parallel()

	tuple := NewStringCharTuple(1, 'a')
	assert.Equal(t, NewNumber(1), tuple.MustGet("@"))
	assert.Equal(t, NewNumber('a'), tuple.MustGet("@char"))
}

func TestStringCharTuple_Names(t *testing.T) {
	t.Parallel()
	assert.True(t, NewStringCharTuple(1, 'a').Names().Equal(NewNames("@", "@char")))
}

func TestStringCharTuple_Project(t *testing.T) {
	t.Parallel()
	AssertEqualValues(t,
		NewTuple(NewIntAttr("@", 1)),
		NewStringCharTuple(1, 'a').Project(NewNames("@")))
	AssertEqualValues(t,
		NewTuple(NewIntAttr("@char", 'a')),
		NewStringCharTuple(1, 'a').Project(NewNames("@char")))
	assert.Nil(t, NewStringCharTuple(1, 'a').Project(NewNames("x")))
}

func TestStringCharTuple_Enumerator(t *testing.T) {
	t.Parallel()
	attrs := map[string]int{}
	for e := NewStringCharTuple(1, 'a').Enumerator(); e.MoveNext(); {
		name, value := e.Current()
		v, has := attrs[name]
		require.False(t, has, "%q => %v %v", name, v, value)
		attrs[name] = int(value.(Number).Float64())
	}
	assert.Equal(t, map[string]int{"@": 1, "@char": 'a'}, attrs)
}
