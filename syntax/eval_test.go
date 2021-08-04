package syntax

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/pkg/importcache"
	"github.com/arr-ai/arrai/rel"
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

// This test ensures that import cache gets reset every time EvaluateExpr is called.
func TestImportCacheResetsAfterEveryCompilation(t *testing.T) {
	t.Parallel()

	fs := ctxfs.CreateTestMemMapFs(t, map[string]string{
		"a.arrai": "1",
	})

	ctx := ctxfs.SourceFsOnto(context.Background(), fs)

	source := "//{./a.arrai}"
	val, err := EvaluateExpr(ctx, "", source)
	require.NoError(t, err)
	rel.RequireEqualValues(t, rel.NewNumber(float64(1)), val)

	// change a.arrai to "2".
	file, err := fs.OpenFile("a.arrai", os.O_RDWR, os.ModeAppend)
	require.NoError(t, err)
	_, err = file.Write([]byte("2"))
	require.NoError(t, err)
	file.Close()

	val, err = EvaluateExpr(ctx, "", source)
	require.NoError(t, err)
	rel.AssertEqualValues(t, rel.NewNumber(float64(2)), val)
}

func TestImportCacheUsesExistingCache(t *testing.T) {
	t.Parallel()

	fs := ctxfs.CreateTestMemMapFs(t, map[string]string{
		"a.arrai": "1",
	})

	ctx := ctxfs.SourceFsOnto(context.Background(), fs)
	ctx = importcache.WithNewImportCache(ctx)
	_, err := importcache.GetOrAddFromCache(
		ctx, "a.arrai",
		func() (rel.Expr, error) { return rel.NewNumber(float64(2)), nil },
	)
	require.NoError(t, err)
	val, err := EvaluateExpr(ctx, "", "//{./a.arrai}")
	require.NoError(t, err)
	rel.AssertEqualValues(t, rel.NewNumber(float64(2)), val)
}
