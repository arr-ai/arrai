package syntax

import (
	"io/ioutil"
	"testing"

	"github.com/arr-ai/arrai/rel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluateBundle(t *testing.T) {
	t.Parallel()

	bundle, err := ioutil.ReadFile("../examples/os/echo.arraiz")
	require.NoError(t, err)

	out, err := EvaluateBundle(bundle, "", "hello", "world")
	require.NoError(t, err)
	assert.Equal(t, "hello world", out.String())
}

func TestEvaluateBundle_NoArgs(t *testing.T) {
	t.Parallel()

	bundle, err := ioutil.ReadFile("../examples/os/echo.arraiz")
	require.NoError(t, err)

	out, err := EvaluateBundle(bundle, "")
	require.NoError(t, err)
	assert.Equal(t, rel.None, out)
}
