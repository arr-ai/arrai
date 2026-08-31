package rel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A failed match on the fast path must allocate nothing: cond tries arms by
// matching and discards the misses, so the miss itself has to be free. The
// user-facing message still exists — behind scopeBuilder.explain.
func TestPatternMissAllocatesNothing(t *testing.T) {
	ctx := context.Background()
	pat, err := NewTuplePattern(
		NewTuplePatternAttr("a", NewFallbackPattern(NewIdentPattern("x"), nil)),
	)
	require.NoError(t, err)
	var notATuple Value = NewNumber(42)
	wrongAttr := NewTuple(NewAttr("b", NewNumber(1)))

	// The builder is reused across attempts, as cond's driver reuses one
	// across arms.
	var b scopeBuilder
	allocs := testing.AllocsPerRun(200, func() {
		b.reset()
		if _, err := pat.Bind(ctx, EmptyScope, notATuple, &b); err != errPatternMismatch {
			t.Fatal(err)
		}
		b.reset()
		if _, err := pat.Bind(ctx, EmptyScope, wrongAttr, &b); err != errPatternMismatch {
			t.Fatal(err)
		}
	})
	assert.Zero(t, allocs, "a discarded miss must not allocate")

	// The explain path still yields the full message.
	err = explainBind(ctx, pat, EmptyScope, wrongAttr)
	require.EqualError(t, err, "couldn't find a in tuple (b: 1)")
}
