package perf

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/importcache"
	"github.com/arr-ai/arrai/syntax"
)

func TestCompileOnlyScratch(t *testing.T) {
	dir, _ := filepath.Abs("reconstruct")
	script := filepath.Join(dir, "vendor", "run.arrai")
	source, _ := os.ReadFile(script)
	wd, _ := os.Getwd()
	require.NoError(t, os.Chdir(filepath.Join(dir, "vendor")))
	defer os.Chdir(wd)
	for i := 0; i < 5; i++ {
		ctx := importcache.WithNewImportCache(
			arraictx.WithArgs(arraictx.InitRunCtx(context.Background()), script))
		_, err := syntax.Compile(ctx, script, string(source))
		require.NoError(t, err)
	}
}
