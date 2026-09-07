package rel

import (
	"context"
	"testing"

	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScopeShadowingAndOrder(t *testing.T) {
	t.Parallel()

	s := EmptyScope.With("a", NewNumber(1)).With("b", NewNumber(2)).With("a", NewNumber(3))
	v, has := s.Get("a")
	require.True(t, has)
	assert.Equal(t, NewNumber(3), v)
	assert.Equal(t, []string{"b", "a"}, s.Names(), "oldest first, shadowed entry not repeated")
	assert.Equal(t, 2, s.Count())
	assert.Equal(t, "{a: 3, b: 2}", s.String())

	_, has = s.With("_", NewNumber(9)).Get("_")
	assert.False(t, has, "_ never binds")
}

func TestScopeUpdateAndWithout(t *testing.T) {
	t.Parallel()

	base := EmptyScope.With("x", NewNumber(1)).With("y", NewNumber(2))
	bindings := EmptyScope.With("y", NewNumber(20)).With("z", NewNumber(30))
	u := base.Update(bindings)
	assert.Equal(t, "{x: 1, y: 20, z: 30}", u.String())
	assert.Equal(t, "{x: 1, y: 2}", base.String(), "base unchanged")
	assert.Equal(t, u.String(), base.Update(EmptyScope).Update(bindings).String())
	assert.Equal(t, bindings.String(), EmptyScope.Update(bindings).String())

	w := u.Without("y")
	assert.Equal(t, "{x: 1, z: 30}", w.String())
	_, has := w.Get("y")
	assert.False(t, has)
	assert.Equal(t, EmptyScope.String(), u.Without("x", "y", "z").String())

	_, err := base.MatchedUpdate(EmptyScope.With("x", NewNumber(1)))
	assert.NoError(t, err, "same value may be rebound")
	_, err = base.MatchedUpdate(EmptyScope.With("x", NewNumber(2)))
	assert.Error(t, err, "different value may not")
}

func TestIdentExprCacheSurvivesShapeChanges(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ident := NewIdentExpr(*parser.NewScanner("v"), "v")
	eval := func(s Scope) Value {
		val, err := ident.Eval(ctx, s)
		require.NoError(t, err)
		return val
	}

	// Same shape twice: second read is served from the cache.
	shape1 := func(n float64) Scope {
		return EmptyScope.With("a", NewNumber(0)).With("v", NewNumber(n)).With("b", NewNumber(0))
	}
	assert.Equal(t, NewNumber(1), eval(shape1(1)))
	assert.Equal(t, NewNumber(2), eval(shape1(2)))

	// Different shapes: the cached address is wrong and must be ignored.
	assert.Equal(t, NewNumber(3), eval(EmptyScope.With("v", NewNumber(3))))
	deep := EmptyScope.With("v", NewNumber(4))
	for i := 0; i < 20; i++ {
		deep = deep.With("pad", NewNumber(float64(i)))
	}
	assert.Equal(t, NewNumber(4), eval(deep))
	// A frame whose slot holds a different name at the cached position.
	assert.Equal(t, NewNumber(5), eval(EmptyScope.Update(
		EmptyScope.With("q", NewNumber(0)).With("r", NewNumber(0)).With("v", NewNumber(5)))))
	// Shadowing wins even when the cache points at the older binding.
	assert.Equal(t, NewNumber(6), eval(EmptyScope.With("v", NewNumber(5)).With("v", NewNumber(6))))

	_, err := ident.Eval(ctx, EmptyScope.With("w", NewNumber(0)))
	assert.ErrorContains(t, err, `name "v" not found in {w}`)
}
