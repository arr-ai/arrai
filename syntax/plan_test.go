package syntax

import (
	"context"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/rel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func evalPlan(t *testing.T, code string) rel.Value {
	t.Helper()
	ctx := arraictx.InitRunCtx(context.Background())
	expr, err := Compile(ctx, "", code)
	require.NoError(t, err)
	p, err := rel.LowerPlan(expr)
	require.NoError(t, err, "lower %s", code)
	b, err := rel.EncodePlan(p)
	require.NoError(t, err)
	p2, err := rel.DecodePlan(b)
	require.NoError(t, err)
	v, err := p2.Eval(ctx, rel.EmptyScope)
	require.NoError(t, err)
	return v
}

// A $-string decoded from a plan must carry a real source scanner (🎯T25
// regression): liftXstr used to leave it at the zero value, so formatting an
// error raised inside the interpolation panicked in wbnf's Scanner.Contains
// instead of returning a message.
func TestPlanRoundtripXstrErrorHasSource(t *testing.T) {
	t.Parallel()
	code := `$"${ {0: 1, 1: 1}(2) }"`
	ctx := arraictx.InitRunCtx(context.Background())

	expr, err := Compile(ctx, "", code)
	require.NoError(t, err)
	p, err := rel.LowerPlan(expr)
	require.NoError(t, err)
	b, err := rel.EncodePlan(p)
	require.NoError(t, err)
	p2, err := rel.DecodePlan(b)
	require.NoError(t, err)
	_, planErr := p2.Eval(ctx, rel.EmptyScope)
	require.Error(t, planErr)

	// The plan loses real source text and position (by design: that's the
	// whole point of skipping wbnf parse), so only the underlying message is
	// checked, not an exact match against a source-parsed error.
	var msg string
	assert.NotPanics(t, func() { msg = planErr.Error() })
	assert.Contains(t, msg, "Call: no return values for input 2 from set {0: 1, 1: 1}")
}

func TestPlanRoundtripEqualsEval(t *testing.T) {
	t.Parallel()
	cases := []string{
		`1 + 2`,
		`{1, 2, 3} count`,
		`{|k, a| (1, 10), (2, 20)} where .k = 1`,
		`[0, 1, 2] >> . + 1 >> . * 2`,
		`let r = {|k, a| (1, 10), (2, 20)}; let p = r => (k: .a); [p where .k = 10, r where .k = 10]`,
		`//str.lower("ABC")`,
		`1 if 2 > 0 else 3`,
		`(a: 1, b: 2).a`,
		`let f = \x x + 1; f(41)`,
		`$"hello ${1+2}"`,
		`[1, 2, 3](1)`,
		`( '': 'App143' ).'' rank (:.@)`,
		`(x: ( '': 'hi' )).x?.''?:'no'`,
	}
	ctx := arraictx.InitRunCtx(context.Background())
	for _, code := range cases {
		code := code
		t.Run(code, func(t *testing.T) {
			t.Parallel()
			direct, err := EvaluateExpr(ctx, "", code)
			require.NoError(t, err)
			got := evalPlan(t, code)
			assert.True(t, got.Equal(direct), "plan=%s eval=%s", got, direct)
		})
	}
}
