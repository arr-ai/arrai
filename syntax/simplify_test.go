package syntax

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/rel"
)

// countingScope binds `tick` to a native function that counts its calls,
// so a test can observe how many times a let's rhs actually ran.
func countingScope(calls *atomic.Int64) rel.Scope {
	return rel.EmptyScope.With("tick", rel.NewNativeFunction("tick", func(context.Context, rel.Value) (rel.Value, error) {
		return rel.NewNumber(float64(calls.Add(1))), nil
	}))
}

func evalCounting(t *testing.T, src string) (rel.Value, int64) {
	t.Helper()
	var calls atomic.Int64
	ctx := arraictx.InitRunCtx(context.Background())
	expr, err := Compile(ctx, NoPath, src)
	require.NoError(t, err)
	v, err := expr.Eval(ctx, countingScope(&calls))
	require.NoError(t, err)
	return v, calls.Load()
}

// TestSimplifyDoesNotDuplicateWork is the standing oracle for 🎯T28: a let
// whose value is used inside a mapped body is evaluated once, before and
// after simplification, whatever the simplifier does to the tree.
func TestSimplifyDoesNotDuplicateWork(t *testing.T) {
	t.Parallel()

	for _, c := range []struct{ src, want string }{
		{`let x = tick(0); {1, 2, 3} => x + .`, `{2, 3, 4}`},
		{`let x = tick(0); {1, 2, 3} where . > x`, `{2, 3}`},
		{`let x = tick(0); [1, 2, 3] >> . + x`, `[2, 3, 4]`},
		{`let x = tick(0); let f = \y y + x; f(1) + f(2)`, `5`},
		{`let x = tick(0); x + x`, `2`},
		{`let x = tick(0); x + 1`, `2`},
		{`let x = tick(0); cond {x > 0: x, _: 0}`, `1`},
	} {
		v, calls := evalCounting(t, c.src)
		require.Equal(t, int64(1), calls, "%s evaluated its rhs %d times", c.src, calls)
		want, err := EvaluateExpr(arraictx.InitRunCtx(context.Background()), NoPath, c.want)
		require.NoError(t, err)
		require.True(t, want.Equal(v), "%s = %v, want %s", c.src, v, c.want)
	}
}

// TestSimplifyFoldsIntoRewrites checks the payoff: a where behind a let
// binding is recognised the same way as one written inline.
func TestSimplifyFoldsIntoRewrites(t *testing.T) {
	t.Parallel()
	if !rel.SimplifyEnabled() {
		t.Skip("simplifier is off in the slowpath build")
	}

	ctx := arraictx.InitRunCtx(context.Background())
	folded, err := Compile(ctx, NoPath, `let r = {|a, b| (1, 2), (3, 4)}; r where .a = 1`)
	require.NoError(t, err)
	inline, err := Compile(ctx, NoPath, `{|a, b| (1, 2), (3, 4)} where .a = 1`)
	require.NoError(t, err)
	require.Equal(t, inline.String(), folded.String())
}

// TestSimplifyKeepsEffectOrder: an rhs with a side effect stays where it
// was written, so its output precedes the body's.
func TestSimplifyKeepsEffectOrder(t *testing.T) {
	t.Parallel()

	ctx := arraictx.InitRunCtx(context.Background())
	e, err := Compile(ctx, NoPath, `let x = //log.print(1); x`)
	require.NoError(t, err)
	_, isLet := e.(*rel.ArrowExpr)
	require.True(t, isLet, "%v", e)

	// A call-free, effect-free rhs does fold, so the guard is what kept it.
	e, err = Compile(ctx, NoPath, `let x = "A" ++ "b"; x`)
	require.NoError(t, err)
	_, isLet = e.(*rel.ArrowExpr)
	require.Equal(t, !rel.SimplifyEnabled(), isLet, "%v", e)
}
