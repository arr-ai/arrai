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

// TestClosureEqualHash pins Closure's equality rule (#777): the same function
// node and equal captured bindings, where "captured" means the function's
// free identifiers, not the frame they live in. Hash must agree wherever
// Equal does.
func TestClosureEqualHash(t *testing.T) {
	t.Parallel()

	sc := *parser.NewScanner(`\a a + k`)
	// \a a + k, with k free.
	addK := NewFunction(sc, IdentPattern("a"), NewAddExpr(sc, NewIdentExpr(sc, "a"), NewIdentExpr(sc, "k"))).(*Function)
	k1 := NewClosure(EmptyScope.With("k", NewNumber(1)), addK)
	k1Again := NewClosure(EmptyScope.With("k", NewNumber(1)), addK)
	k2 := NewClosure(EmptyScope.With("k", NewNumber(2)), addK)

	assert.True(t, k1.Equal(k1), "f = f")
	assert.True(t, k1.Equal(k1Again), "same node, equal capture in a different frame")
	assert.Equal(t, k1.Hash(), k1Again.Hash())
	assert.False(t, k1.Equal(k2), "same node, different capture")
	assert.NotEqual(t, k1.Hash(), k2.Hash())
	assert.False(t, k1.Equal(NewNumber(1)))

	// A binding the body never mentions does not take part.
	k1Extra := NewClosure(EmptyScope.With("k", NewNumber(1)).With("unused", NewNumber(9)), addK)
	assert.True(t, k1.Equal(k1Extra))
	assert.Equal(t, k1.Hash(), k1Extra.Hash())

	// The free identifier bound in one scope but not the other.
	assert.False(t, k1.Equal(NewClosure(EmptyScope, addK)))

	// The same body text compiled twice is two nodes, hence two functions.
	addKCopy := NewFunction(sc, IdentPattern("a"), NewAddExpr(sc, NewIdentExpr(sc, "a"), NewIdentExpr(sc, "k")))
	assert.False(t, k1.Equal(NewClosure(EmptyScope.With("k", NewNumber(1)), addKCopy.(*Function))))

	// A set of closures differing only in a captured value keeps them all.
	set := MustNewSet(k1, k2, k1Again, NewClosure(EmptyScope.With("k", NewNumber(3)), addK))
	assert.Equal(t, 3, set.Count())
	assert.True(t, set.Has(k1Again))

	// Function-less closures (as some tests build) compare equal.
	assert.True(t, NewClosure(Scope{}, nil).Equal(NewClosure(Scope{}, nil)))
	assert.Equal(t, NewClosure(Scope{}, nil).Hash(), NewClosure(Scope{}, nil).Hash())
	assert.False(t, NewClosure(Scope{}, nil).Equal(k1))
}

// TestClosureEqualOpaqueBody checks the fallback for a body the walker cannot
// see into (opaqueExpr, from simplify_test.go): every visible binding takes
// part, so an otherwise-unused binding now distinguishes two closures.
func TestClosureEqualOpaqueBody(t *testing.T) {
	t.Parallel()

	sc := *parser.NewScanner(`\a ?`)
	f := NewFunction(sc, IdentPattern("a"), opaqueExpr{NewIdentExpr(sc, "k")}).(*Function)
	_, known := f.freeIdents()
	require.False(t, known)

	k1 := NewClosure(EmptyScope.With("k", NewNumber(1)), f)
	k1Again := NewClosure(EmptyScope.With("k", NewNumber(1)), f)
	k1Extra := NewClosure(EmptyScope.With("k", NewNumber(1)).With("unused", NewNumber(9)), f)
	k1ExtraFlat := NewClosure(EmptyScope.Update(k1Extra.scope), f)

	assert.True(t, k1.Equal(k1Again))
	assert.Equal(t, k1.Hash(), k1Again.Hash())
	assert.False(t, k1.Equal(k1Extra))
	assert.False(t, k1Extra.Equal(k1))
	// Same bindings in a differently shaped chain: equal, and hashed alike.
	assert.True(t, k1Extra.Equal(k1ExtraFlat))
	assert.Equal(t, k1Extra.Hash(), k1ExtraFlat.Hash())
}

// TestClosureEqualRecursiveBinding checks that closures capturing a `let rec`
// cell compare by the cell's value without looping on the self-reference,
// and that two evaluations of the same recursion with equal captures are
// equal while different captures are not.
func TestClosureEqualRecursiveBinding(t *testing.T) {
	t.Parallel()

	sc := *parser.NewScanner(`\n f(k)`)
	// \n f(k): f and k free; f will be a recursive cell holding this closure.
	body := NewCallExpr(sc, NewIdentExpr(sc, "f"), NewIdentExpr(sc, "k"))
	fn := NewFunction(sc, IdentPattern("n"), body).(*Function)
	tie := func(k Value) Closure {
		cell := &recCell{name: "f"}
		c := NewClosure(EmptyScope.With("k", k).With("f", cell), fn)
		cell.val, cell.set = c, true
		return c
	}
	a, b, c := tie(NewNumber(1)), tie(NewNumber(1)), tie(NewNumber(2))

	assert.True(t, a.Equal(a))
	assert.True(t, a.Equal(b))
	assert.Equal(t, a.Hash(), b.Hash())
	assert.False(t, a.Equal(c))
	assert.Equal(t, 2, MustNewSet(a, b, c).Count())

	// An unset cell (still being evaluated) equals only itself.
	unset := &recCell{name: "f"}
	assert.True(t, sameBinding(unset, unset))
	assert.False(t, sameBinding(unset, &recCell{name: "f"}))
}

// TestSameBindingComparesValues checks that the scope matchers (`[x, x]`
// patterns and friends) compare bindings by value rather than by printed
// form, which is what let two different closures of one function match.
func TestSameBindingComparesValues(t *testing.T) {
	t.Parallel()

	sc := *parser.NewScanner(`\a a + k`)
	addK := NewFunction(sc, IdentPattern("a"), NewAddExpr(sc, NewIdentExpr(sc, "a"), NewIdentExpr(sc, "k"))).(*Function)
	k1 := NewClosure(EmptyScope.With("k", NewNumber(1)), addK)
	k2 := NewClosure(EmptyScope.With("k", NewNumber(2)), addK)
	require.Equal(t, k1.String(), k2.String(), "the printed forms do not tell them apart")

	assert.True(t, sameBinding(k1, k1))
	assert.False(t, sameBinding(k1, k2))
	assert.True(t, sameBinding(NewNumber(1), NewNumber(1)))
	assert.False(t, sameBinding(NewNumber(1), NewNumber(2)))
	assert.True(t, sameBinding(nil, nil))
	assert.False(t, sameBinding(nil, NewNumber(1)))

	var b scopeBuilder
	require.NoError(t, b.add("x", k1))
	assert.NoError(t, b.add("x", k1))
	assert.ErrorIs(t, b.add("x", k2), errPatternMismatch)

	s := EmptyScope.With("x", k1)
	_, err := s.MatchedWith("x", k1)
	assert.NoError(t, err)
	_, err = s.MatchedWith("x", k2)
	assert.Error(t, err)
	_, err = s.MatchedUpdate(EmptyScope.With("x", k2))
	assert.Error(t, err)
}
