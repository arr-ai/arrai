package shell

import (
	"fmt"
	"testing"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
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

	sh := newShellInstance(newLineCollector(), rel.EmptyScope)
	assert.NoError(t, tryRunCommand(`/set a = 1 + 2`, sh))
	assert.NoError(t, tryRunCommand(`/unset a`, sh))

	assert.EqualError(t, tryRunCommand("/hi", sh), "command hi not found")
	assert.EqualError(t, tryRunCommand("random", sh), "random is not a command")
}

func TestSetCmd(t *testing.T) {
	t.Parallel()

	set := &setCmd{}
	assert.Equal(t, "set", set.name())

	sh := newShellInstance(newLineCollector(), rel.EmptyScope)

	setCmdAssertion(t, "a", `{"hello": 123}`, set, sh)
	setCmdAssertion(t, "a123", "123", set, sh)

	errMsg := `/set command error, usage: /set <name> = <expr>`
	assert.EqualError(t, set.process("a 1+2", sh), errMsg)
	assert.EqualError(t, set.process("= 1+2", sh), errMsg)
}

func setCmdAssertion(t *testing.T, name, exprToSet string, set *setCmd, sh *shellInstance) {
	assert.NoError(t, set.process(fmt.Sprintf("%s = %s", name, exprToSet), sh))
	actualVal, err := syntax.EvalWithScope("", exprToSet, sh.scope)
	require.NoError(t, err)
	expr, exists := sh.scope.Get(name)
	assert.True(t, exists)
	rel.AssertExprsEvalToSameValue(t, expr, actualVal)
}

func TestUnsetCmd(t *testing.T) {
	t.Parallel()

	unset := &unsetCmd{}
	assert.Equal(t, "unset", unset.name())

	sh := newShellInstance(newLineCollector(), rel.EmptyScope)
	require.True(t, sh.scope.Count() == 0)
	assert.NoError(t, unset.process("x", sh))

	sh.scope = sh.scope.With("a", rel.NewNumber(123))
	require.NotPanics(t, func() {
		sh.scope.MustGet("a")
	})

	assert.NoError(t, unset.process("a", sh))
	_, exists := sh.scope.Get("a")
	assert.False(t, exists)

	assert.EqualError(t, unset.process("123", sh), "/unset command error, usage: /unset <name>")
}
