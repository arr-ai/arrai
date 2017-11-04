package tests

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/rel"
)

// TestNewSet tests rel.NewSet.
func TestNewSet(t *testing.T) {
	assert.NotNil(t, rel.NewSet())
	assert.NotNil(t, rel.NewSet(rel.NewNumber(42)))
}

// TestNewSetFrom tests rel.TestNewSetFrom.
func TestNewSetFrom(t *testing.T) {
	a, err := rel.NewSetFrom()
	if assert.NoError(t, err) {
		assert.EqualValues(t, 0, a.Count())
		assert.Equal(t, rel.NewSet(), a)
	}
	a, err = rel.NewSetFrom(42)
	if assert.NoError(t, err) {
		assert.EqualValues(t, 1, a.Count())
		assert.Equal(t, rel.NewSet(rel.NewNumber(42)), a)
	}
}

// TestNewSetFromWithBadInput tests rel.TestNewSetFrom(somethingAwful).
func TestNewSetFromWithBadInput(t *testing.T) {
	_, err := rel.NewSetFrom(func() {})
	assert.Error(t, err)
}

// TestSetHash tests rel.Set.Hash.
func TestSetHash(t *testing.T) {
	a := intSet()
	b := intSet(42)
	c := intSet(42)
	d := intSet(4321)
	e := intSet(4321, 321)
	f := intSet(321, 4321)

	b.Hash(0)
	assert.Equal(t, b.Hash(0), c.Hash(0))
	assert.Equal(t, e.Hash(0), f.Hash(0))

	allSets := []rel.Set{a, b, c, d, e, f}
	for _, x := range allSets {
		for i := uint32(0); i < 10; i++ {
			assert.NotEqual(t, 0, x.Hash(i))
		}
	}

	distinctSets := []rel.Set{a, b, d, e}
	for _, x := range distinctSets {
		for _, y := range distinctSets {
			for i := uint32(0); i < 10; i++ {
				if x == y {
					assert.Equal(t, x.Hash(i), y.Hash(i),
						"%s.Hash(%d) != %s.Hash(%[2]d)", x, i, y)
					assert.NotEqual(t, x.Hash(i), y.Hash(i+1),
						"%s.Hash(%d) == %s.Hash(%[2]d+1)", x, i, y)
				} else {
					assert.NotEqual(t, x.Hash(i), y.Hash(i),
						"%s.Hash(%s) == %s.Hash(%[2]s)", x, i, y)
				}
			}
		}
	}
}

// TestSetEqual tests rel.Set.Equal.
func TestSetEqual(t *testing.T) {
	a := intSet()
	b := intSet(42)
	c := intSet(42)
	d := intSet(4321)
	e := intSet(4321, 321)
	f := intSet(321, 4321)

	assert.True(t, b.Equal(c))
	assert.True(t, c.Equal(b))
	assert.True(t, e.Equal(f))
	assert.True(t, f.Equal(e))

	distinctSets := []rel.Set{a, b, d, e}
	for _, x := range distinctSets {
		for _, y := range distinctSets {
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

// Float64InterfaceList represents an []interface{} of float64 for sort.Sort().
type Float64InterfaceList []interface{}

func (vl Float64InterfaceList) Len() int {
	return len(([]interface{})(vl))
}

func (vl Float64InterfaceList) Less(i, j int) bool {
	values := ([]interface{})(vl)
	return values[i].(float64) < values[j].(float64)
}

func (vl Float64InterfaceList) Swap(i, j int) {
	values := ([]interface{})(vl)
	values[i], values[j] = values[j], values[i]
}

// TestSetBool tests rel.Set.Bool.
func TestSetBool(t *testing.T) {
	assert.False(t, intSet().Bool())
	assert.True(t, intSet(42).Bool())
}

// TestSetLess tests rel.Set.Less.
func TestSetLess(t *testing.T) {
	a := []rel.Set{
		intSet(),
		intSet(41),
		intSet(42, 41),
		intSet(42),
		intSet(43),
	}

	for i, x := range a {
		for j, y := range a {
			assert.Equal(t, i < j, x.Less(y), "a[%d] < a[%d]", i, j)
		}
	}
}

// TestSetExport tests rel.Set.Export.
func TestSetExport(t *testing.T) {
	scenario := func(intfs ...interface{}) {
		if intfs == nil {
			intfs = []interface{}{}
		}
		v, err := rel.NewSetFrom(intfs...)
		if assert.NoError(t, err) {
			a := v.Export().([]interface{})
			sort.Sort(Float64InterfaceList(a))
			assert.Equal(t, intfs, a)
		}
	}
	scenario()
	scenario(42.0)
	scenario(42.0)
	scenario(5432.0)
	scenario(321.0, 4321.0)
}

// TestSetString tests rel.Set.String.
func TestSetString(t *testing.T) {
	scenario := func(repr string, values ...interface{}) {
		set := intSet(values...)
		if assert.Equal(
			t, uint64(len(values)), set.Count(), "%v", set.Export(),
		) {
			assert.Equal(t, repr, set.String(), "%v", set)
		}
	}
	scenario(`none`)
	scenario(`{|42|}`, 42)
	scenario(`{|5432|}`, 5432)
	scenario(`{|321, 4321|}`, 4321, 321)
	scenario(`{|321, 4321|}`, 321, 4321)
}

// TestSetCount tests rel.Set.Count.
func TestSetCount(t *testing.T) {
	attrs := []interface{}{42, 43, 44}
	assert.EqualValues(t, 0, intSet().Count())
	assert.EqualValues(t, 1, intSet(attrs[0:1]...).Count())
	assert.EqualValues(t, 1, intSet(attrs[1:2]...).Count())
	assert.EqualValues(t, 1, intSet(attrs[2:3]...).Count())
	assert.EqualValues(t, 2, intSet(attrs[0:2]...).Count())
	assert.EqualValues(t, 2, intSet(attrs[1:3]...).Count())
	assert.EqualValues(t, 3, intSet(attrs[0:3]...).Count())
}

// TestSetHas tests rel.Set.Has.
func TestSetHas(t *testing.T) {
	set := intSet(42, 43)
	assert.True(t, set.Has(rel.NewNumber(42)))
	assert.True(t, set.Has(rel.NewNumber(43)))
	assert.False(t, set.Has(rel.NewNumber(44)))
}

// TestSetWith tests rel.Set.With.
func TestSetWith(t *testing.T) {
	n42 := rel.NewNumber(42)
	n43 := rel.NewNumber(43)
	n44 := rel.NewNumber(44)

	set := intSet()

	assertHas := func(value rel.Value) {
		assert.True(t, set.Has(value))
	}

	set = set.With(n42)
	assertHas(n42)

	set = set.With(n43)
	assertHas(n42)
	assertHas(n43)

	set = set.With(n44)
	assertHas(n42)
	assertHas(n43)
	assertHas(n44)
}

// TestSetWithout tests rel.Set.Without.
func TestSetWithout(t *testing.T) {
	n42 := rel.NewNumber(42)
	n43 := rel.NewNumber(43)
	n44 := rel.NewNumber(44)

	set := rel.NewSet(n42, n43, n44)
	assert.EqualValues(t, 3, set.Count())

	assert.True(t, set.Has(n42))
	assert.True(t, set.Has(n43))
	assert.True(t, set.Has(n44))

	set = set.Without(n44)
	assert.EqualValues(t, 2, set.Count())
	assert.True(t, set.Has(n42))
	assert.True(t, set.Has(n43))
	assert.False(t, set.Has(n44))

	set = set.Without(n42)
	assert.EqualValues(t, 1, set.Count())
	assert.False(t, set.Has(n42))
	assert.True(t, set.Has(n43))
	assert.False(t, set.Has(n44))

	set = set.Without(n44)
	assert.EqualValues(t, 1, set.Count())

	set = set.Without(n43)
	assert.EqualValues(t, 0, set.Count())
	assert.False(t, set.Has(n42))
	assert.False(t, set.Has(n43))
	assert.False(t, set.Has(n44))

	set = set.Without(n42)
	assert.EqualValues(t, 0, set.Count())

	set = set.Without(n43)
	assert.EqualValues(t, 0, set.Count())
}

// TestSetWalk test rel.Set.Enumerator().
func TestSetWalk(t *testing.T) {
	set := intSet(42, 43, 44, 45, 46)

	a := []interface{}{}
	for e := set.Enumerator(); e.MoveNext(); {
		a = append(a, e.Current().Export())
	}
	sort.Sort(Float64InterfaceList(a))
	assert.Equal(t, []interface{}{42.0, 43.0, 44.0, 45.0, 46.0}, a)
}
