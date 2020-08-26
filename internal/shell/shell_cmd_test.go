package shell

import (
	"context"
	"fmt"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsCommand(t *testing.T) {
	t.Parallel()

	assert.True(t, isCommand("/hi"))
	assert.False(t, isCommand("//seq.join"))
}

func TestTryRunCommand(t *testing.T) {
	t.Parallel()

	sh := newShellInstance(newLineCollector(), []rel.ContextErr{})
	ctx := arraictx.InitRunCtx(context.Background())
	assert.NoError(t, tryRunCommand(ctx, `/set a = 1 + 2`, sh))
	assert.NoError(t, tryRunCommand(ctx, `/set $a = 1 + 2`, sh))
	assert.NoError(t, tryRunCommand(ctx, `/set . = 1 + 2`, sh))
	assert.NoError(t, tryRunCommand(ctx, `/set _a = 1 + 2`, sh))
	assert.NoError(t, tryRunCommand(ctx, `/unset a`, sh))
	assert.NoError(t, tryRunCommand(ctx, `/unset $a`, sh))
	assert.NoError(t, tryRunCommand(ctx, `/unset .`, sh))
	assert.NoError(t, tryRunCommand(ctx, `/unset _a`, sh))

	assert.EqualError(t, tryRunCommand(ctx, "/hi", sh), "command hi not found")
	assert.EqualError(t, tryRunCommand(ctx, "random", sh), "random is not a command")
}

func TestSetCmd(t *testing.T) {
	t.Parallel()

	set := &setCmd{}
	assert.Equal(t, []string{"set"}, set.names())

	sh := newShellInstance(newLineCollector(), []rel.ContextErr{})

	setCmdAssertion(t, "a", `{"hello": 123}`, set, sh)
	setCmdAssertion(t, "a123", "123", set, sh)

	errMsg := `/set command error, usage: /set <name> = <expr>`
	ctx := arraictx.InitRunCtx(context.Background())
	assert.EqualError(t, set.process(ctx, "a 1+2", sh), errMsg)
	assert.EqualError(t, set.process(ctx, "= 1+2", sh), errMsg)
}

func setCmdAssertion(t *testing.T, name, exprToSet string, set *setCmd, sh *shellInstance) {
	ctx := arraictx.InitRunCtx(context.Background())
	assert.NoError(t, set.process(ctx, fmt.Sprintf("%s = %s", name, exprToSet), sh))
	actualVal, err := syntax.EvalWithScope(ctx, "", exprToSet, sh.scope)
	require.NoError(t, err)
	expr, exists := sh.scope.Get(name)
	assert.True(t, exists)
	rel.AssertExprsEvalToSameValue(t, expr, actualVal)
}

func TestUnsetCmd(t *testing.T) {
	t.Parallel()

	unset := &unsetCmd{}
	assert.Equal(t, []string{"unset"}, unset.names())

	sh := newShellInstance(newLineCollector(), []rel.ContextErr{})
	ctx := arraictx.InitRunCtx(context.Background())

	require.True(t, sh.scope.Count() == syntax.StdScope().Count())
	assert.NoError(t, unset.process(ctx, "x", sh))

	sh.scope = sh.scope.With("a", rel.NewNumber(123))
	require.NotPanics(t, func() {
		sh.scope.MustGet("a")
	})

	assert.NoError(t, unset.process(ctx, "a", sh))
	_, exists := sh.scope.Get("a")
	assert.False(t, exists)

	assert.EqualError(t, unset.process(ctx, "123", sh), "/unset command error, usage: /unset <name>")
}

func TestUpFrameCmd(t *testing.T) {
	t.Parallel()

	up := &upFrameCmd{}
	assert.Equal(t, []string{"up", "u"}, up.names())

	sh := newShellInstance(newLineCollector(), []rel.ContextErr{})
	ctx := arraictx.InitRunCtx(context.Background())

	assert.EqualError(t, up.process(ctx, "", sh), "frame index out of range, frame length: 0")
	assertEqualScope(t, removeArraiInfo(syntax.StdScope()), removeArraiInfo(sh.scope))

	ctxErrs := createContextErrs()
	sh = newShellInstance(newLineCollector(), ctxErrs)
	assertCurrScope(t, sh, 2, ctxErrs)
	require.NoError(t, up.process(ctx, "", sh))
	assertCurrScope(t, sh, 1, ctxErrs)
	require.NoError(t, up.process(ctx, "", sh))
	assertCurrScope(t, sh, 0, ctxErrs)
	assert.EqualError(t, up.process(ctx, "", sh), "frame index out of range, frame length: 3")
}

func TestDownFrameCmd(t *testing.T) {
	t.Parallel()

	down := &downFrameCmd{}
	assert.Equal(t, []string{"down", "d"}, down.names())

	sh := newShellInstance(newLineCollector(), []rel.ContextErr{})
	ctx := arraictx.InitRunCtx(context.Background())

	assert.EqualError(t, down.process(ctx, "", sh), "frame index out of range, frame length: 0")
	assertEqualScope(t, removeArraiInfo(syntax.StdScope()), removeArraiInfo(sh.scope))

	ctxErrs := createContextErrs()
	sh = newShellInstance(newLineCollector(), ctxErrs)
	sh.currentFrameIndex = 0
	sh.scope = syntax.StdScope().Update(ctxErrs[0].GetScope())

	assertCurrScope(t, sh, 0, ctxErrs)
	require.NoError(t, down.process(ctx, "", sh))
	assertCurrScope(t, sh, 1, ctxErrs)
	require.NoError(t, down.process(ctx, "", sh))
	assertCurrScope(t, sh, 2, ctxErrs)
	assert.EqualError(t, down.process(ctx, "", sh), "frame index out of range, frame length: 3")
}

func assertCurrScope(t *testing.T, sh *shellInstance, index int, frames []rel.ContextErr) {
	assert.Equal(t, index, sh.currentFrameIndex)
	assertEqualScope(t, removeArraiInfo(syntax.StdScope().Update(frames[index].GetScope())), removeArraiInfo(sh.scope))
}

func createContextErrs() []rel.ContextErr {
	baseErr := rel.NewContextErr(
		fmt.Errorf("random"),
		*parser.NewScannerAt("random", 5, 1),
		rel.EmptyScope.With("random", rel.NewNumber(1)),
	)
	nextErr := rel.NewContextErr(
		baseErr,
		*parser.NewScannerAt("random", 3, 1),
		rel.EmptyScope.With("random", rel.NewNumber(2)),
	)
	nextErr2 := rel.NewContextErr(
		nextErr,
		*parser.NewScannerAt("random", 1, 1),
		rel.EmptyScope.With("random", rel.NewNumber(3)),
	)
	return nextErr2.GetImportantFrames()
}

func assertEqualScope(t *testing.T, expected, actual rel.Scope) {
	assert.Equal(t, expected.Count(), actual.Count())
	for e := expected.Enumerator(); e.MoveNext(); {
		name, v1 := e.Current()
		v2, exists := actual.Get(name)
		assert.True(t, exists)
		assert.True(t, v1.(rel.Value).Equal(v2))
	}
}

// removeArraiInfo removes `//arrai.info` from the scope as it can't be constructed
// in this test due to its lack of package `main`.
func removeArraiInfo(scope rel.Scope) rel.Scope {
	root, _ := scope.Get("//")
	rootTuple, _ := root.(rel.Tuple)

	arrai, _ := rootTuple.Get("arrai")
	arraiTuple, _ := arrai.(rel.Tuple)

	return scope.With("//", rootTuple.With("arrai", arraiTuple.Without("info")))
}
