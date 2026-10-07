package rel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// 🎯T45 (arr-ai/arrai#779): a native function is a Go function, so a plan
// node can carry only its name, and names are not unique. Lowering one must
// fail rather than write a node that decodes to whichever native registered
// last under that name.
func TestPlanLowerRejectsNativeFunction(t *testing.T) {
	t.Parallel()
	native := NewNativeFunction("decode", func(_ context.Context, v Value) (Value, error) { return v, nil })

	_, err := LowerPlan(native)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot bundle native function ⦑decode⦒")
	require.Contains(t, err.Error(), "#779")

	// The macro path wraps the value in a LiteralExpr; the error must surface
	// through that too, and through a containing tuple.
	_, err = LowerPlan(NewLiteralExpr(planSrc, native))
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot bundle native function")

	_, err = LowerPlan(NewTuple(NewAttr("f", native)))
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot bundle native function")
}

// A closure's one-element-set view used to send encodeValue into
// encodeSetEnum and back into itself until the stack overflowed. A closure
// now lowers to its function plus captured bindings and lifts back to a
// closure that evaluates identically.
func TestPlanClosureRoundtrip(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	// \n n + k, with k = 10 captured and an unused m to show it is left out.
	body := NewAddExpr(planSrc, NewIdentExpr(planSrc, "n"), NewIdentExpr(planSrc, "k"))
	fn := NewFunction(planSrc, NewIdentPattern("n"), body).(*Function)
	scope := EmptyScope.With("k", NewNumber(10)).With("m", NewNumber(99))
	closure := NewClosure(scope, fn)

	p, err := LowerPlan(NewLiteralExpr(planSrc, closure))
	require.NoError(t, err)
	lit := p.Root
	require.Equal(t, "lit", lit.K)
	require.Equal(t, planClosureKind, lit.Kids[0].K)
	require.Len(t, lit.Kids[0].Kids, 2, "fn + one capture (k); m is not free in the body")
	require.Equal(t, "k", lit.Kids[0].Kids[1].Attr)

	b, err := EncodePlan(p)
	require.NoError(t, err)
	p2, err := DecodePlan(b)
	require.NoError(t, err)
	v, err := p2.Eval(ctx, EmptyScope)
	require.NoError(t, err)
	got, ok := v.(Closure)
	require.True(t, ok, "%T", v)
	out, err := SetCall(ctx, got, NewNumber(5))
	require.NoError(t, err)
	require.True(t, out.Equal(NewNumber(15)), "%s", out)
}

// A captured native still fails, through the closure's captures.
func TestPlanClosureCapturingNativeIsRejected(t *testing.T) {
	t.Parallel()
	native := NewNativeFunction("decode", func(_ context.Context, v Value) (Value, error) { return v, nil })
	body := NewCallExpr(planSrc, NewIdentExpr(planSrc, "d"), NewIdentExpr(planSrc, "s"))
	fn := NewFunction(planSrc, NewIdentPattern("s"), body).(*Function)
	closure := NewClosure(EmptyScope.With("d", native), fn)
	_, err := LowerPlan(NewLiteralExpr(planSrc, closure))
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot bundle native function")
}
