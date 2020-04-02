package rel //nolint:dupl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBytesByteTuple(t *testing.T) {
	t.Parallel()
	type args struct {
		at      int
		byteval byte
	}
	tests := []struct {
		name string
		args args
		want BytesByteTuple
	}{
		{name: "default", want: BytesByteTuple{}},
		{name: "0@0", want: BytesByteTuple{}, args: args{at: 0, byteval: 0}},
		{name: "0@1", want: BytesByteTuple{at: 1, byteval: 0}, args: args{at: 1, byteval: 0}},
		{name: "a@0", want: BytesByteTuple{at: 0, byteval: 'a'}, args: args{at: 0, byteval: 'a'}},
		{name: "a@1", want: BytesByteTuple{at: 1, byteval: 'a'}, args: args{at: 1, byteval: 'a'}},
		{name: "a@-1", want: BytesByteTuple{at: -1, byteval: 'a'}, args: args{at: -1, byteval: 'a'}},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			assert.EqualValues(t, test.want, NewBytesByteTuple(test.args.at, test.args.byteval))
		})
	}
}

func TestNewBytesByteTupleFromTuple(t *testing.T) {
	t.Parallel()
	type args struct {
		t Tuple
	}
	tests := []struct {
		name string
		args args
		want BytesByteTuple
		ok   bool
	}{
		{
			name: "0@0",
			want: NewBytesByteTuple(0, 0),
			ok:   true,
			args: args{t: NewTuple(NewIntAttr("@", 0), NewIntAttr(BytesByteAttr, 0))},
		},
		{
			name: "a@0",
			want: NewBytesByteTuple(0, 'a'),
			ok:   true,
			args: args{t: NewTuple(NewIntAttr("@", 0), NewIntAttr(BytesByteAttr, 'a'))},
		},
		{
			name: "a@42",
			want: NewBytesByteTuple(42, 'a'),
			ok:   true,
			args: args{t: NewTuple(NewIntAttr("@", 42), NewIntAttr(BytesByteAttr, 'a'))},
		},
		{
			name: "no-@",
			want: BytesByteTuple{},
			ok:   false,
			args: args{t: NewTuple(NewIntAttr("at", 0), NewIntAttr(BytesByteAttr, 'a'))},
		},
		{
			name: "no-ByteAttr",
			want: BytesByteTuple{},
			ok:   false,
			args: args{t: NewTuple(NewIntAttr("@", 0), NewIntAttr("byte", 'a'))},
		},
	}
	for _, test := range tests { //nolint:dupl
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, ok := newBytesByteTupleFromTuple(test.args.t)
			got2 := maybeNewBytesByteTupleFromTuple(test.args.t)
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

func TestBytesByteTuple_Equal(t *testing.T) {
	t.Parallel()
	values := []Value{
		None,
		NewNumber(42),
		NewBytesByteTuple(0, 0),
		NewBytesByteTuple(0, 'a'),
		NewBytesByteTuple(42, 0),
		NewBytesByteTuple(42, 'a'),
		NewBytesByteTuple(42, 'b'),
		NewBytesByteTuple(43, 'b'),
	}
	for i, x := range values {
		for j, y := range values {
			assert.Equal(t, i == j, x.Equal(y), "values[%d]=%v, values[%d]=%v", i, x, j, y)
		}
	}
}

func TestBytesByteTuple_String(t *testing.T) { //nolint:dupl
	t.Parallel()
	tests := []struct {
		name  string
		tuple BytesByteTuple
		want  string
	}{
		{name: "0@0", want: "(@: 0, @byte: 0)", tuple: NewBytesByteTuple(0, 0)},
		{name: "a@0", want: "(@: 0, @byte: 97)", tuple: NewBytesByteTuple(0, 'a')},
		{name: "0@42", want: "(@: 42, @byte: 0)", tuple: NewBytesByteTuple(42, 0)},
		{name: "a@42", want: "(@: 42, @byte: 97)", tuple: NewBytesByteTuple(42, 'a')},
		{name: "b@42", want: "(@: 42, @byte: 98)", tuple: NewBytesByteTuple(42, 'b')},
		{name: "b@43", want: "(@: 43, @byte: 98)", tuple: NewBytesByteTuple(43, 'b')},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.tuple.String())
		})
	}
}

func TestBytesByteTuple_Eval(t *testing.T) {
	t.Parallel()
	tuple := NewBytesByteTuple(42, 'a')
	value, err := tuple.Eval(EmptyScope)
	require.NoError(t, err)
	assert.Equal(t, tuple, value)
}

func TestBytesByteTuple_IsTrue(t *testing.T) {
	t.Parallel()
	values := []Value{
		NewBytesByteTuple(0, 0),
		NewBytesByteTuple(0, 'a'),
		NewBytesByteTuple(42, 0),
		NewBytesByteTuple(42, 'z'),
	}
	for i, x := range values {
		assert.True(t, x.IsTrue(), "values[%d]=%v", i, x)
	}
}

func TestBytesByteTuple_Less(t *testing.T) {
	t.Parallel()
	assert.True(t, !NewBytesByteTuple(0, 0).Less(NewBytesByteTuple(0, 0)))
	assert.True(t, NewBytesByteTuple(0, 0).Less(NewBytesByteTuple(0, 'a')))
	assert.True(t, NewBytesByteTuple(0, 'a').Less(NewBytesByteTuple(42, 0)))
}

func TestBytesByteTuple_Negate(t *testing.T) {
	t.Parallel()
	assert.Equal(t, NewBytesByteTuple(0, 0), NewBytesByteTuple(0, 0).Negate())
}

func TestBytesByteTuple_Export(t *testing.T) {
	t.Parallel()
	assert.Equal(t, map[string]interface{}{"@": 1, "@byte": uint8('a')}, NewBytesByteTuple(1, 'a').Export())
}

func TestBytesByteTuple_Count(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 2, NewBytesByteTuple(0, 0).Count())
	assert.Equal(t, 2, NewBytesByteTuple(1, 'a').Count())
}

func TestBytesByteTuple_Get(t *testing.T) { //nolint:dupl
	t.Parallel()

	assertGet := func(tuple Tuple, attr string, value Value) {
		v, has := tuple.Get(attr)
		if assert.True(t, has) {
			assert.Equal(t, value, v)
		}
	}
	assertGet(NewBytesByteTuple(1, 'a'), "@", NewNumber(1))
	assertGet(NewBytesByteTuple(1, 'a'), "@byte", NewNumber('a'))

	assertNotGet := func(tuple Tuple, attr string) bool {
		v, has := tuple.Get(attr)
		return assert.False(t, has, "%q => %v", attr, v)
	}
	assertNotGet(NewBytesByteTuple(1, 'a'), "@@")
}

func TestBytesByteTuple_MustGet(t *testing.T) {
	t.Parallel()
	assert.Equal(t, NewNumber(1), NewBytesByteTuple(1, 'a').MustGet("@"))
	assert.Equal(t, NewNumber('a'), NewBytesByteTuple(1, 'a').MustGet("@byte"))
	assert.Panics(t, func() { NewBytesByteTuple(1, 'a').MustGet("byteval") })
}

func TestBytesByteTuple_With(t *testing.T) {
	t.Parallel()
	AssertEqualValues(t, NewBytesByteTuple(42, 'a'), NewBytesByteTuple(1, 'a').With("@", NewNumber(42)))
	AssertEqualValues(t, NewBytesByteTuple(1, 'b'), NewBytesByteTuple(1, 'a').With("@byte", NewNumber('b')))
	AssertEqualValues(t, NewTuple(
		NewIntAttr("@", 1),
		NewIntAttr("@byte", 'a'),
		NewIntAttr("x", 'b'),
	), NewBytesByteTuple(1, 'a').With("x", NewNumber('b')))
}

func TestBytesByteTuple_Without(t *testing.T) {
	t.Parallel()
	AssertEqualValues(t, NewTuple(NewIntAttr("@byte", 'a')), NewBytesByteTuple(1, 'a').Without("@"))
	AssertEqualValues(t, NewTuple(NewIntAttr("@", 1)), NewBytesByteTuple(1, 'a').Without("@byte"))
	AssertEqualValues(t, NewBytesByteTuple(1, 'a'), NewBytesByteTuple(1, 'a').Without("x"))
}

func TestBytesByteTuple_Map(t *testing.T) {
	t.Parallel()
	AssertEqualValues(t,
		NewBytesByteTuple(2, 'b'),
		NewBytesByteTuple(1, 'a').Map(func(v Value) Value {
			return NewNumber(v.(Number).Float64() + 1)
		}),
	)
}

func TestBytesByteTuple_HasName(t *testing.T) {
	t.Parallel()
	assert.True(t, NewBytesByteTuple(1, 'a').HasName("@"))
	assert.True(t, NewBytesByteTuple(1, 'a').HasName("@byte"))
	assert.False(t, NewBytesByteTuple(1, 'a').HasName("@item"))
}

func TestBytesByteTuple_Attributes(t *testing.T) {
	t.Parallel()
	assert.Equal(t,
		map[string]Value{"@": NewNumber(1), "@byte": NewNumber('a')},
		NewBytesByteTuple(1, 'a').Attributes())
}

func TestBytesByteTuple_Names(t *testing.T) {
	t.Parallel()
	assert.True(t, NewBytesByteTuple(1, 'a').Names().Equal(NewNames("@", "@byte")))
}

func TestBytesByteTuple_Project(t *testing.T) {
	t.Parallel()
	AssertEqualValues(t,
		NewTuple(NewIntAttr("@", 1)),
		NewBytesByteTuple(1, 'a').Project(NewNames("@")))
	AssertEqualValues(t,
		NewTuple(NewIntAttr("@byte", 'a')),
		NewBytesByteTuple(1, 'a').Project(NewNames("@byte")))
	assert.Nil(t, NewBytesByteTuple(1, 'a').Project(NewNames("x")))
}

func TestBytesByteTuple_Enumerator(t *testing.T) {
	t.Parallel()
	attrs := map[string]int{}
	for e := NewBytesByteTuple(1, 'a').Enumerator(); e.MoveNext(); {
		name, value := e.Current()
		v, has := attrs[name]
		require.False(t, has, "%q => %v %v", name, v, value)
		attrs[name] = int(value.(Number).Float64())
	}
	assert.Equal(t, map[string]int{"@": 1, "@byte": 'a'}, attrs)
}
