package rel

import (
	"context"
	"errors"
	"testing"

	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFunctionValuesAreSingletonSets checks that every function kind behaves
// as the one-element set containing itself under the structural Set methods,
// rather than panicking (🎯T32).
func TestFunctionValuesAreSingletonSets(t *testing.T) {
	t.Parallel()

	sc := *parser.NewScanner(`\x x`)
	identity := func(_ context.Context, v Value) (Value, error) { return v, nil }
	closure := func() Set {
		return NewClosure(Scope{}, NewFunction(sc, IdentPattern("x"), NewNumber(1)).(*Function))
	}

	cases := []struct {
		name string
		f    Set // the function under test
		g    Set // a distinct function of the same kind
	}{
		{
			name: "NativeFunction",
			f:    NewNativeFunction("one", identity).(Set),
			g:    NewNativeFunction("two", identity).(Set),
		},
		{
			name: "Closure",
			f:    closure(),
			g:    closure(),
		},
		{
			name: "ExprClosure",
			f:    NewExprClosure(Scope{}, NewNumber(1)).(Set),
			g:    NewExprClosure(Scope{}, NewNumber(2)).(Set),
		},
	}

	elements := func(e ValueEnumerator) []Value {
		var out []Value
		for e.MoveNext() {
			out = append(out, e.Current())
		}
		return out
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			f, g := c.f, c.g
			five := NewNumber(5)

			assert.Equal(t, 1, f.Count())
			assert.True(t, f.IsTrue())

			assert.True(t, f.Has(f))
			assert.False(t, f.Has(g))
			assert.False(t, f.Has(five))

			got := elements(f.Enumerator())
			require.Len(t, got, 1)
			assert.True(t, got[0].Equal(f))
			got = elements(f.ArrayEnumerator())
			require.Len(t, got, 1)
			assert.True(t, got[0].Equal(f))

			assert.True(t, f.With(f).Equal(f))
			withFive := f.With(five)
			assert.Equal(t, 2, withFive.Count())
			assert.True(t, withFive.Has(f))
			assert.True(t, withFive.Has(five))
			withG := f.With(g)
			assert.Equal(t, 2, withG.Count())
			assert.True(t, withG.Has(g))

			assert.True(t, f.Without(f).Equal(None))
			assert.True(t, f.Without(five).Equal(f))
			assert.True(t, f.Without(g).Equal(f))

			mapped, err := f.Map(func(Value) (Value, error) { return five, nil })
			require.NoError(t, err)
			assert.True(t, mapped.Equal(MustNewSet(five)))
			mapped, err = f.Map(func(v Value) (Value, error) { return v, nil })
			require.NoError(t, err)
			assert.Equal(t, 1, mapped.Count())
			assert.True(t, mapped.Has(f))
			boom := errors.New("boom")
			_, err = f.Map(func(Value) (Value, error) { return nil, boom })
			assert.ErrorIs(t, err, boom)

			kept, err := f.Where(func(v Value) (bool, error) { return v.Equal(f), nil })
			require.NoError(t, err)
			assert.True(t, kept.Equal(f))
			dropped, err := f.Where(func(Value) (bool, error) { return false, nil })
			require.NoError(t, err)
			assert.True(t, dropped.Equal(None))
			_, err = f.Where(func(Value) (bool, error) { return false, boom })
			assert.ErrorIs(t, err, boom)
		})
	}
}

// TestExprClosureEqualHash pins ExprClosure's identity semantics: equal iff
// the same scope frame and the same comparable expression node, never
// panicking on uncomparable expression types.
func TestExprClosureEqualHash(t *testing.T) {
	t.Parallel()

	a := NewExprClosure(Scope{}, NewNumber(1))
	b := NewExprClosure(Scope{}, NewNumber(1))
	c := NewExprClosure(Scope{}, NewNumber(2))
	assert.True(t, a.Equal(b))
	assert.Equal(t, a.Hash(), b.Hash())
	assert.False(t, a.Equal(c))
	assert.False(t, a.Equal(NewNumber(1)))

	scoped := NewExprClosure(EmptyScope.With("x", NewNumber(1)), NewNumber(1))
	assert.False(t, a.Equal(scoped))

	// An Array value holds a slice, so it is not comparable; two closures
	// over such a node are unequal, but comparing them must not panic.
	arr := NewArray(NewNumber(1))
	assert.NotPanics(t, func() {
		assert.False(t, NewExprClosure(Scope{}, arr).Equal(NewExprClosure(Scope{}, arr)))
	})

	assert.True(t, NewExprClosure(Scope{}, nil).Equal(NewExprClosure(Scope{}, nil)))
	assert.False(t, NewExprClosure(Scope{}, nil).Equal(a))
}
