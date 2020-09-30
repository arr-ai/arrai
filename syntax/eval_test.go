package syntax

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluateBundle(t *testing.T) {
	bundle, err := ioutil.ReadFile("../examples/os/echo.arraiz")
	require.NoError(t, err)

	tools.Arguments = nil
	out, err := EvaluateBundle(arraictx.InitRunCtx(context.Background()), bundle, "", "hello", "world")
	require.NoError(t, err)
	assert.Equal(t, "hello world", out.String())

	tools.Arguments = nil
}

// TODO: Make tools.Arguments (and StdScope generally) slightly more mutable: new per eval?
//func TestEvaluateBundle_Empty(t *testing.T) {
//	bundle, err := ioutil.ReadFile("../examples/os/echo.arraiz")
//	require.NoError(t, err)
//
//	tools.Arguments = nil
//	out, err := EvaluateBundle(arraictx.InitRunCtx(context.Background()), bundle, "")
//	require.NoError(t, err)
//	assert.Equal(t, rel.None, out)
//
//	tools.Arguments = nil
//}
