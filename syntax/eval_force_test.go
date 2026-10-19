package syntax

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/rel"
)

// TestEvaluateExprReturnsConcreteArray is the arr-ai/arrai#789 repro: a >>
// result, at top level or nested in a returned tuple, must reach Go callers
// as rel.Array, not a suspended pipeline.
func TestEvaluateExprReturnsConcreteArray(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	v, err := EvaluateExpr(ctx, ".", `['a', 'b'] >> .`)
	require.NoError(t, err)
	a, ok := v.(rel.Array)
	require.True(t, ok, "got %T", v)
	assert.Len(t, a.Values(), 2)

	v, err = EvaluateExpr(ctx, ".", `(appName: [(s: 'a'), (s: 'b')] >> .s)`)
	require.NoError(t, err)
	tup, ok := v.(rel.Tuple)
	require.True(t, ok, "got %T", v)
	attr, _ := tup.Get("appName")
	_, ok = attr.(rel.Array)
	assert.True(t, ok, "appName: got %T", attr)
}

func TestEvaluateExprForcesNestedPipelines(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	for _, src := range []string{
		`[1, 2] >> . + 1 >> . * 2`,
		`[1, 2] >>> \i \x i + x`,
		`let f = \x x >> . + 1; f([1, 2])`,
		`[[1] >> ., 2\[3] >> .]`,
	} {
		v, err := EvaluateExpr(ctx, ".", src)
		require.NoError(t, err, src)
		a, ok := v.(rel.Array)
		require.True(t, ok, "%s: got %T", src, v)
		for _, item := range a.Values() {
			if s, isSet := item.(rel.Set); isSet && s.Count() > 0 {
				_, ok := item.(rel.Array)
				assert.True(t, ok, "%s: element %T", src, item)
			}
		}
	}

	v, err := EvaluateExpr(ctx, ".", `{'a': 1} >> . + 1`)
	require.NoError(t, err)
	_, ok := v.(rel.Dict)
	assert.True(t, ok, "dict >>: got %T", v)

	v, err = EvaluateExpr(ctx, ".", `{'k': [1] >> ., 'd': {'x': 1} >> .}`)
	require.NoError(t, err)
	d, ok := v.(rel.Dict)
	require.True(t, ok, "got %T", v)
	for e := d.DictEnumerator(); e.MoveNext(); {
		k, item := e.Current()
		switch k.String() {
		case "k":
			_, ok = item.(rel.Array)
		case "d":
			_, ok = item.(rel.Dict)
		}
		assert.True(t, ok, "%v: got %T", k, item)
	}

	v, err = EvaluateExpr(ctx, ".", `{[1] >> .}`)
	require.NoError(t, err)
	s, ok := v.(rel.Set)
	require.True(t, ok)
	for e := s.Enumerator(); e.MoveNext(); {
		_, ok := e.Current().(rel.Array)
		assert.True(t, ok, "set element: got %T", e.Current())
	}
}

// TestEvaluateExprSurfacesDeferredPipelineError: a failing pure >> stage
// must fail Eval, as it did before pipelines were suspended, rather than
// returning a value whose Export panics.
func TestEvaluateExprSurfacesDeferredPipelineError(t *testing.T) {
	t.Parallel()
	_, err := EvaluateExpr(context.Background(), ".", `[1, 'a'] >> . + 1`)
	require.Error(t, err)
	_, err = EvaluateExpr(context.Background(), ".", `(x: [1, 'a'] >> . + 1)`)
	require.Error(t, err)
}

// TestCompileEvalForces covers embedders that Compile once and Eval per
// scope (review of #790): the nested attribute and the deferred error.
func TestCompileEvalForces(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	expr, err := Compile(ctx, ".", `(appName: x >> .s)`)
	require.NoError(t, err)
	x, err := EvaluateExpr(ctx, ".", `[(s: 'a'), (s: 'b')]`)
	require.NoError(t, err)
	v, err := expr.Eval(ctx, rel.Scope{}.With("x", x))
	require.NoError(t, err)
	attr, _ := v.(rel.Tuple).Get("appName")
	_, ok := attr.(rel.Array)
	assert.True(t, ok, "appName: got %T", attr)
	assert.Equal(t, rel.Unforced(expr).String(), expr.String(), "String delegates")
	assert.NotEqual(t, expr, rel.Unforced(expr), "Compile wraps a non-constant expr")

	bad, err := Compile(ctx, ".", `[1, 'a'] >> . + 1`)
	require.NoError(t, err)
	_, err = bad.Eval(ctx, rel.EmptyScope)
	require.Error(t, err)

	bad, err = Compile(ctx, ".", `\x (r: x >> . + 1)`)
	require.NoError(t, err)
	f, err := bad.Eval(ctx, rel.EmptyScope)
	require.NoError(t, err)
	_, err = rel.CallForced(ctx, f, rel.NewArray(rel.NewNumber(1), rel.NewString([]rune("a"))))
	require.Error(t, err)

	v, err = MustCompile(ctx, ".", `[1, 2] >> . + 1`).Eval(ctx, rel.EmptyScope)
	require.NoError(t, err)
	_, ok = v.(rel.Array)
	assert.True(t, ok, "MustCompile: got %T", v)
}

// TestCallForcedFromGo: a closure returned to Go and called with
// rel.CallForced yields concrete values; rel.Call stays lazy, as it is the
// in-language call.
func TestCallForcedFromGo(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	f, err := EvaluateExpr(ctx, ".", `\xs (appName: xs >> .s)`)
	require.NoError(t, err)
	arg, err := EvaluateExpr(ctx, ".", `[(s: 'a'), (s: 'b')]`)
	require.NoError(t, err)

	v, err := rel.CallForced(ctx, f, arg)
	require.NoError(t, err)
	attr, _ := v.(rel.Tuple).Get("appName")
	_, ok := attr.(rel.Array)
	assert.True(t, ok, "CallForced: got %T", attr)
}

// TestCompiledPlanEvalForces: a lowered, encoded and decoded plan (the
// compiled-bundle path) also returns concrete values from Plan.Eval.
func TestCompiledPlanEvalForces(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	expr, err := Compile(ctx, ".", `(appName: [(s: 'a'), (s: 'b')] >> .s)`)
	require.NoError(t, err)
	p, err := rel.LowerPlan(expr)
	require.NoError(t, err)
	b, err := rel.EncodePlan(p)
	require.NoError(t, err)
	p, err = rel.DecodePlan(b)
	require.NoError(t, err)
	v, err := p.Eval(ctx, rel.EmptyScope)
	require.NoError(t, err)
	attr, _ := v.(rel.Tuple).Get("appName")
	_, ok := attr.(rel.Array)
	assert.True(t, ok, "plan: got %T", attr)
}
