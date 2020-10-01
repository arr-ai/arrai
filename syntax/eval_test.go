package syntax

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/arr-ai/arrai/rel"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluateBundle(t *testing.T) {
	bundle, err := ioutil.ReadFile("../examples/os/echo.arraiz")
	require.NoError(t, err)

	out, err := EvaluateBundle(arraictx.InitRunCtx(context.Background()), bundle, "", "hello", "world")
	require.NoError(t, err)
	assert.Equal(t, "hello world", out.String())
}

func TestEvaluateBundle_Empty(t *testing.T) {
	bundle, err := ioutil.ReadFile("../examples/os/echo.arraiz")
	require.NoError(t, err)

	out, err := EvaluateBundle(arraictx.InitRunCtx(context.Background()), bundle, "")
	require.NoError(t, err)
	assert.Equal(t, rel.None, out)
}
