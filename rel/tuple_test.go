package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTupleNewTuple(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, EmptyTuple)
	assert.NotNil(t, Attr{"a", NewNumber(42)})
	assert.NotNil(t,
		Attr{"a", NewNumber(42)},
		Attr{"a", NewNumber(42)},
	)
}

func TestTupleNewTupleFromMap(t *testing.T) {
	t.Parallel()
	tuple, err := NewTupleFromMap(map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, tuple)

	tuple, err = NewTupleFromMap(map[string]interface{}{"a": 42})
	assert.NoError(t, err)
	assert.NotNil(t, tuple)

	tuple, err = NewTupleFromMap(map[string]interface{}{"a": 42, "b": 43})
	assert.NoError(t, err)
	assert.NotNil(t, tuple)
}

func TestTupleHash(t *testing.T) {
	t.Parallel()
	a := EmptyTuple
	b := NewTuple(Attr{"a", NewNumber(42)})
	c := NewTuple(Attr{"a", NewNumber(42)})
	d := NewTuple(Attr{"A", NewNumber(42)})
	e := NewTuple(Attr{"a", NewNumber(4321)})
	f := NewTuple(
		Attr{"a", NewNumber(4321)},
		Attr{"b", NewNumber(321)},
	)
	g := NewTuple(
		Attr{"b", NewNumber(321)},
		Attr{"a", NewNumber(4321)},
	)

	assert.Equal(t, b.Hash(0), c.Hash(0), "should hash the same")
	assert.Equal(t, f.Hash(0), g.Hash(0), "should hash the same")

	allTuples := []Tuple{a, b, c, d, e, f, g}
	for _, x := range allTuples {
		for i := uintptr(0); i < 10; i++ {
			assert.NotEqual(t, 0, x.Hash(i), "shouldn't hash to zero")
		}
	}

	distinctTuples := []Tuple{a, b, d, e, f}
	for _, x := range distinctTuples {
		for _, y := range distinctTuples {
			for i := uintptr(0); i < 10; i++ {
				if x == y {
					assert.Equal(t, x.Hash(i), y.Hash(i), "should hash stably")
					hx, hy := x.Hash(i), y.Hash(i+1)
					assert.NotEqual(t, hx, hy,
						"%s and %s should hash differently for different "+
							"seeds, not %d and %d",
						x, y, hx, hy)
				} else {
					assert.NotEqual(t, x.Hash(i), y.Hash(i),
						"should hash differently")
				}
			}
		}
	}
}

func TestTupleEqual(t *testing.T) {
	t.Parallel()
	a := EmptyTuple
	b := NewTuple(Attr{"a", NewNumber(42)})
	c := NewTuple(Attr{"a", NewNumber(42)})
	d := NewTuple(Attr{"A", NewNumber(42)})
	e := NewTuple(Attr{"a", NewNumber(4321)})
	f := NewTuple(
		Attr{"a", NewNumber(4321)},
		Attr{"b", NewNumber(321)},
	)
	g := NewTuple(
		Attr{"b", NewNumber(321)},
		Attr{"a", NewNumber(4321)},
	)

	assert.True(t, b.Equal(c))
	assert.True(t, c.Equal(b))
	assert.True(t, f.Equal(g))
	assert.True(t, g.Equal(f))

	distinctTuples := []Tuple{a, b, d, e, f}
	for _, x := range distinctTuples {
		for _, y := range distinctTuples {
			if x == y {
				assert.True(t, x.Equal(y))
				assert.True(t, x.Equal(y))
			} else {
				assert.False(t, x.Equal(y))
				assert.False(t, y.Equal(x))
			}
		}
	}
}

func TestTupleBool(t *testing.T) {
	t.Parallel()
	assert.False(t, EmptyTuple.Bool())
	assert.True(t, NewTuple(Attr{"a", NewNumber(42)}).Bool())
}

func TestTupleLess(t *testing.T) {
	t.Parallel()
	attr := func(name string, n float64) Attr {
		return Attr{name, NewNumber(n)}
	}

	a := []Tuple{
		EmptyTuple,
		NewTuple(attr("a", 41)),
		NewTuple(attr("b", 42), attr("a", 41)),
		NewTuple(attr("a", 42)),
		NewTuple(attr("a", 43)),
		NewTuple(attr("b", 42)),
	}

	for i, x := range a {
		for j, y := range a {
			assert.Equal(t, i < j, x.Less(y), "a[%d] < a[%d]", i, j)
		}
	}
}

func TestTupleExport(t *testing.T) {
	t.Parallel()
	scenario := func(m map[string]interface{}) {
		v, err := NewTupleFromMap(m)
		if assert.NoError(t, err) {
			assert.Equal(t, m, v.Export())
		}
	}
	scenario(map[string]interface{}{})
	scenario(map[string]interface{}{"a": 42.0})
	scenario(map[string]interface{}{"b": 42.0})
	scenario(map[string]interface{}{"a": 5432.0})
	scenario(map[string]interface{}{"a": 4321.0, "b": 321.0})
}

func TestTupleString(t *testing.T) {
	t.Parallel()
	scenario := func(repr string, attrs ...Attr) {
		tuple := NewTuple(attrs...)
		if assert.Equal(
			t, len(attrs), tuple.Count(), "%v", tuple.Export(),
		) {
			assert.Equal(t, repr, tuple.String(), "%v", tuple)
		}
	}
	scenario("()")
	scenario("(a: 42)", Attr{"a", NewNumber(42)})
	scenario("(b: 42)", Attr{"b", NewNumber(42)})
	scenario("(a: 5432)", Attr{"a", NewNumber(5432)})
	scenario("(a: 4321, b: 321)",
		Attr{"a", NewNumber(4321)},
		Attr{"b", NewNumber(321)},
	)
	scenario("(a: 4321, b: 321)",
		Attr{"b", NewNumber(321)},
		Attr{"a", NewNumber(4321)},
	)
}

func TestTupleCount(t *testing.T) {
	t.Parallel()
	attrs := []Attr{
		{"a", NewNumber(42)},
		{"b", NewNumber(43)},
		{"c", NewNumber(44)},
	}
	assert.EqualValues(t, 0, EmptyTuple.Count())
	assert.EqualValues(t, 1, NewTuple(attrs[:1]...).Count())
	assert.EqualValues(t, 1, NewTuple(attrs[1:2]...).Count())
	assert.EqualValues(t, 1, NewTuple(attrs[2:]...).Count())
	assert.EqualValues(t, 2, NewTuple(attrs[:2]...).Count())
	assert.EqualValues(t, 2, NewTuple(attrs[1:]...).Count())
	assert.EqualValues(t, 3, NewTuple(attrs...).Count())
}

func TestTupleGet(t *testing.T) {
	t.Parallel()
	tuple := NewTuple(
		Attr{"a", NewNumber(42)},
		Attr{"b", NewNumber(43)},
	)
	a, found := tuple.Get("a")
	if assert.True(t, found) {
		assert.Equal(t, 42.0, a.Export())
	}
	b, found := tuple.Get("b")
	if assert.True(t, found) {
		assert.Equal(t, 43.0, b.Export())
	}
	c, found := tuple.Get("c")
	if found {
		assert.Fail(t, "found non-existent key", "%s", c)
	}
}

func TestTupleWith(t *testing.T) {
	t.Parallel()
	n42 := NewNumber(42)
	n43 := NewNumber(43)
	n44 := NewNumber(44)

	tuple := EmptyTuple

	assertAttr := func(name string, value Value) {
		got, found := tuple.Get(name)
		if assert.True(t, found) {
			assert.Equal(t, value, got)
		}
	}

	tuple = tuple.With("a", n42)
	assertAttr("a", n42)

	tuple = tuple.With("a", n43)
	assertAttr("a", n43)

	tuple = tuple.With("b", n44)
	assertAttr("a", n43)
	assertAttr("b", n44)

	tuple = tuple.With("c", n44)
	assertAttr("a", n43)
	assertAttr("b", n44)
	assertAttr("c", n44)
}

func TestTupleWithout(t *testing.T) {
	t.Parallel()
	n42 := NewNumber(42)
	n43 := NewNumber(43)
	n44 := NewNumber(44)

	tuple := NewTuple(
		Attr{"a", n42},
		Attr{"b", n43},
		Attr{"c", n44},
	)

	assertAttr := func(name string, value Value) {
		got, found := tuple.Get(name)
		if assert.True(t, found) {
			assert.Equal(t, value, got)
		}
	}

	assertNoAttr := func(name string) {
		got, found := tuple.Get(name)
		if found {
			assert.Fail(t, "attribute not removed", "%s -> %s", name, got)
		}
	}

	assertAttr("a", n42)
	assertAttr("b", n43)
	assertAttr("c", n44)

	tuple = tuple.Without("d")
	tuple = tuple.Without("c")
	assertAttr("a", n42)
	assertAttr("b", n43)
	assertNoAttr("c")

	tuple = tuple.Without("a")
	assertNoAttr("a")
	assertAttr("b", n43)
	assertNoAttr("c")

	tuple = tuple.Without("c")

	tuple = tuple.Without("b")
	assertNoAttr("a")
	assertNoAttr("b")
	assertNoAttr("c")

	tuple = tuple.Without("a")

	tuple = tuple.Without("b")
}

func TestTupleHasName(t *testing.T) {
	t.Parallel()
	tuple := NewTuple(
		Attr{"a", NewNumber(42)},
		Attr{"b", NewNumber(43)},
	)
	assert.True(t, tuple.HasName("a"))
	assert.True(t, tuple.HasName("b"))
	assert.False(t, tuple.HasName("c"))
}

func TestAttrEnumeratorToMap(t *testing.T) {
	t.Parallel()
	n42 := NewNumber(42)
	n43 := NewNumber(43)
	n44 := NewNumber(44)
	n45 := NewNumber(45)

	tuple := NewTuple(
		Attr{"a", n42},
		Attr{"b", n43},
		Attr{"c", n44},
		Attr{"d", n45},
	)

	m := AttrEnumeratorToMap(tuple.Enumerator())
	assert.Equal(t,
		map[string]Value{"a": n42, "b": n43, "c": n44, "d": n45},
		m,
	)
}

func TestTupleEnumerator(t *testing.T) {
	t.Parallel()
	tuple := NewTuple(
		Attr{"a", NewNumber(42)},
		Attr{"b", NewNumber(43)},
		Attr{"c", NewNumber(44)},
		Attr{"d", NewNumber(45)},
		Attr{"e", NewNumber(46)},
	)

	m := map[string]interface{}{}
	for e := tuple.Enumerator(); e.MoveNext(); {
		name, value := e.Current()
		m[name] = value.Export()
	}
	assert.Equal(t,
		map[string]interface{}{
			"a": 42.0, "b": 43.0, "c": 44.0, "d": 45.0, "e": 46.0,
		},
		m,
	)
}
