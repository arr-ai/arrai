package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/rel"
)

// TestTupleNewTuple tests rel.NewTuple.
func TestTupleNewTuple(t *testing.T) {
	assert.NotNil(t, rel.EmptyTuple)
	assert.NotNil(t, rel.Attr{"a", rel.NewNumber(42)})
	assert.NotNil(t,
		rel.Attr{"a", rel.NewNumber(42)},
		rel.Attr{"a", rel.NewNumber(42)},
	)
}

// TestTupleNewTupleFromMap tests rel.NewTupleFromMap.
func TestTupleNewTupleFromMap(t *testing.T) {
	tuple, err := rel.NewTupleFromMap(map[string]interface{}{})
	assert.NoError(t, err)
	assert.NotNil(t, tuple)

	tuple, err = rel.NewTupleFromMap(map[string]interface{}{"a": 42})
	assert.NoError(t, err)
	assert.NotNil(t, tuple)

	tuple, err = rel.NewTupleFromMap(map[string]interface{}{"a": 42, "b": 43})
	assert.NoError(t, err)
	assert.NotNil(t, tuple)
}

// TestTupleHash tests rel.Tuple.Hash.
func TestTupleHash(t *testing.T) {
	a := rel.EmptyTuple
	b := rel.NewTuple(rel.Attr{"a", rel.NewNumber(42)})
	c := rel.NewTuple(rel.Attr{"a", rel.NewNumber(42)})
	d := rel.NewTuple(rel.Attr{"A", rel.NewNumber(42)})
	e := rel.NewTuple(rel.Attr{"a", rel.NewNumber(4321)})
	f := rel.NewTuple(
		rel.Attr{"a", rel.NewNumber(4321)},
		rel.Attr{"b", rel.NewNumber(321)},
	)
	g := rel.NewTuple(
		rel.Attr{"b", rel.NewNumber(321)},
		rel.Attr{"a", rel.NewNumber(4321)},
	)

	assert.Equal(t, b.Hash(0), c.Hash(0), "should hash the same")
	assert.Equal(t, f.Hash(0), g.Hash(0), "should hash the same")

	allTuples := []rel.Tuple{a, b, c, d, e, f, g}
	for _, x := range allTuples {
		for i := uint32(0); i < 10; i++ {
			assert.NotEqual(t, 0, x.Hash(i), "shouldn't hash to zero")
		}
	}

	distinctTuples := []rel.Tuple{a, b, d, e, f}
	for _, x := range distinctTuples {
		for _, y := range distinctTuples {
			for i := uint32(0); i < 10; i++ {
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

// TestTupleEqual tests rel.Tuple.Equal.
func TestTupleEqual(t *testing.T) {
	a := rel.EmptyTuple
	b := rel.NewTuple(rel.Attr{"a", rel.NewNumber(42)})
	c := rel.NewTuple(rel.Attr{"a", rel.NewNumber(42)})
	d := rel.NewTuple(rel.Attr{"A", rel.NewNumber(42)})
	e := rel.NewTuple(rel.Attr{"a", rel.NewNumber(4321)})
	f := rel.NewTuple(
		rel.Attr{"a", rel.NewNumber(4321)},
		rel.Attr{"b", rel.NewNumber(321)},
	)
	g := rel.NewTuple(
		rel.Attr{"b", rel.NewNumber(321)},
		rel.Attr{"a", rel.NewNumber(4321)},
	)

	assert.True(t, b.Equal(c))
	assert.True(t, c.Equal(b))
	assert.True(t, f.Equal(g))
	assert.True(t, g.Equal(f))

	distinctTuples := []rel.Tuple{a, b, d, e, f}
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

// TestTupleBool tests rel.Tuple.Bool.
func TestTupleBool(t *testing.T) {
	assert.False(t, rel.EmptyTuple.Bool())
	assert.True(t, rel.NewTuple(rel.Attr{"a", rel.NewNumber(42)}).Bool())
}

// TestTupleLess tests rel.Tuple.Less.
func TestTupleLess(t *testing.T) {
	attr := func(name string, n float64) rel.Attr {
		return rel.Attr{name, rel.NewNumber(n)}
	}

	a := []rel.Tuple{
		rel.EmptyTuple,
		rel.NewTuple(attr("a", 41)),
		rel.NewTuple(attr("b", 42), attr("a", 41)),
		rel.NewTuple(attr("a", 42)),
		rel.NewTuple(attr("a", 43)),
		rel.NewTuple(attr("b", 42)),
	}

	for i, x := range a {
		for j, y := range a {
			assert.Equal(t, i < j, x.Less(y), "a[%d] < a[%d]", i, j)
		}
	}
}

// TestTupleExport tests rel.Tuple.Export.
func TestTupleExport(t *testing.T) {
	scenario := func(m map[string]interface{}) {
		v, err := rel.NewTupleFromMap(m)
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

// TestTupleString tests rel.Tuple.String.
func TestTupleString(t *testing.T) {
	scenario := func(repr string, attrs ...rel.Attr) {
		tuple := rel.NewTuple(attrs...)
		if assert.Equal(
			t, uint64(len(attrs)), tuple.Count(), "%v", tuple.Export(),
		) {
			assert.Equal(t, repr, tuple.String(), "%v", tuple)
		}
	}
	scenario("{}")
	scenario("{a: 42}", rel.Attr{"a", rel.NewNumber(42)})
	scenario("{a: 42}", rel.Attr{"a", rel.NewNumber(42)})
	scenario("{b: 42}", rel.Attr{"b", rel.NewNumber(42)})
	scenario("{a: 5432}", rel.Attr{"a", rel.NewNumber(5432)})
	scenario("{a: 4321, b: 321}",
		rel.Attr{"a", rel.NewNumber(4321)},
		rel.Attr{"b", rel.NewNumber(321)},
	)
	scenario("{a: 4321, b: 321}",
		rel.Attr{"b", rel.NewNumber(321)},
		rel.Attr{"a", rel.NewNumber(4321)},
	)
}

// TestTupleCount tests rel.Tuple.Count.
func TestTupleCount(t *testing.T) {
	attrs := []rel.Attr{
		{"a", rel.NewNumber(42)},
		{"b", rel.NewNumber(43)},
		{"c", rel.NewNumber(44)},
	}
	assert.EqualValues(t, 0, rel.EmptyTuple.Count())
	assert.EqualValues(t, 1, rel.NewTuple(attrs[:1]...).Count())
	assert.EqualValues(t, 1, rel.NewTuple(attrs[1:2]...).Count())
	assert.EqualValues(t, 1, rel.NewTuple(attrs[2:]...).Count())
	assert.EqualValues(t, 2, rel.NewTuple(attrs[:2]...).Count())
	assert.EqualValues(t, 2, rel.NewTuple(attrs[1:]...).Count())
	assert.EqualValues(t, 3, rel.NewTuple(attrs...).Count())
}

// TestTupleGet tests rel.Tuple.Get.
func TestTupleGet(t *testing.T) {
	tuple := rel.NewTuple(
		rel.Attr{"a", rel.NewNumber(42)},
		rel.Attr{"b", rel.NewNumber(43)},
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

// TestTupleWith tests rel.Tuple.With.
func TestTupleWith(t *testing.T) {
	n42 := rel.NewNumber(42)
	n43 := rel.NewNumber(43)
	n44 := rel.NewNumber(44)

	tuple := rel.EmptyTuple

	assertAttr := func(name string, value rel.Value) {
		got, found := tuple.Get(name)
		if assert.True(t, found) {
			assert.Equal(t, value, got)
		}
	}

	tuple, added := tuple.With("a", n42)
	require.True(t, added)
	assertAttr("a", n42)

	tuple, added = tuple.With("a", n43)
	assert.False(t, added)
	assertAttr("a", n43)

	tuple, added = tuple.With("b", n44)
	assertAttr("a", n43)
	assertAttr("b", n44)

	tuple, added = tuple.With("c", n44)
	assertAttr("a", n43)
	assertAttr("b", n44)
	assertAttr("c", n44)
}

// TestTupleWithout tests rel.Tuple.Without.
func TestTupleWithout(t *testing.T) {
	n42 := rel.NewNumber(42)
	n43 := rel.NewNumber(43)
	n44 := rel.NewNumber(44)

	tuple := rel.NewTuple(
		rel.Attr{"a", n42},
		rel.Attr{"b", n43},
		rel.Attr{"c", n44},
	)

	assertAttr := func(name string, value rel.Value) {
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

	tuple, removed := tuple.Without("d")
	require.False(t, removed)
	require.False(t, removed)

	tuple, removed = tuple.Without("c")
	require.True(t, removed)
	assertAttr("a", n42)
	assertAttr("b", n43)
	assertNoAttr("c")

	tuple, removed = tuple.Without("a")
	assertNoAttr("a")
	assertAttr("b", n43)
	assertNoAttr("c")

	tuple, removed = tuple.Without("c")
	assert.False(t, removed)

	tuple, removed = tuple.Without("b")
	assertNoAttr("a")
	assertNoAttr("b")
	assertNoAttr("c")

	tuple, removed = tuple.Without("a")
	assert.False(t, removed)

	tuple, removed = tuple.Without("b")
	assert.False(t, removed)
}

// TestTupleHasName tests rel.Tuple.HasName.
func TestTupleHasName(t *testing.T) {
	tuple := rel.NewTuple(
		rel.Attr{"a", rel.NewNumber(42)},
		rel.Attr{"b", rel.NewNumber(43)},
	)
	assert.True(t, tuple.HasName("a"))
	assert.True(t, tuple.HasName("b"))
	assert.False(t, tuple.HasName("c"))
}

// TestAttrEnumeratorToMap test rel.Tuple.Attributes.
func TestAttrEnumeratorToMap(t *testing.T) {
	n42 := rel.NewNumber(42)
	n43 := rel.NewNumber(43)
	n44 := rel.NewNumber(44)
	n45 := rel.NewNumber(45)

	tuple := rel.NewTuple(
		rel.Attr{"a", n42},
		rel.Attr{"b", n43},
		rel.Attr{"c", n44},
		rel.Attr{"d", n45},
	)

	m := rel.AttrEnumeratorToMap(tuple.Enumerator())
	assert.Equal(t,
		map[string]rel.Value{"a": n42, "b": n43, "c": n44, "d": n45},
		m,
	)
}

// TestTupleEnumerator test rel.Tuple.Enumerator.
func TestTupleEnumerator(t *testing.T) {
	tuple := rel.NewTuple(
		rel.Attr{"a", rel.NewNumber(42)},
		rel.Attr{"b", rel.NewNumber(43)},
		rel.Attr{"c", rel.NewNumber(44)},
		rel.Attr{"d", rel.NewNumber(45)},
		rel.Attr{"e", rel.NewNumber(46)},
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
