// nolint: lll
package bundle

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type bundleTestCase struct {
	name, path           string
	files, expectedFiles map[string]string
}

func TestBundleFiles(t *testing.T) {
	t.Parallel()
	cases := []bundleTestCase{
		{
			"local dependencies", "/github.com/test/test/test.arrai",
			map[string]string{
				SentinelPath("/github.com/test/test"): "module github.com/test/test\n",
				"/github.com/test/test/test.arrai":    "1",
			},
			map[string]string{
				syntax.BundleConfig: ConfigFile(
					"github.com/test/test",
					ModuleFile("/github.com/test/test/test.arrai"),
				),
				SentinelFile("/github.com/test/test"):          "module github.com/test/test\n",
				ModuleFile("/github.com/test/test/test.arrai"): "1",
			},
		},
		{
			"local dependencies with nested root", "/github.com/test/test/test.arrai",
			map[string]string{
				SentinelPath("/github.com/test/test"):               "module github.com/test/test\n",
				"/github.com/test/test/test.arrai":                  "//{./module/module2/module.arrai}",
				SentinelPath("/github.com/test/test/module/"):       "module github.com/test/test/module\n",
				"/github.com/test/test/module/1.arrai":              "1",
				"/github.com/test/test/module/module2/module.arrai": "//{/1.arrai}",
			},
			map[string]string{
				syntax.BundleConfig: ConfigFile(
					"github.com/test/test",
					ModuleFile("/github.com/test/test/test.arrai"),
				),
				SentinelFile("/github.com/test/test"):                           "module github.com/test/test\n",
				SentinelFile("/github.com/test/test/module/"):                   "module github.com/test/test/module\n",
				ModuleFile("/github.com/test/test/test.arrai"):                  "//{./module/module2/module.arrai}",
				ModuleFile("/github.com/test/test/module/module2/module.arrai"): "//{/1.arrai}",
				ModuleFile("/github.com/test/test/module/1.arrai"):              "1",
			},
		},
		{
			"remote import", "/github.com/test/test/test.arrai",
			map[string]string{
				SentinelPath("/github.com/test/test"): "module github.com/test/test\n",
				"/github.com/test/test/test.arrai":    "//{https://raw.githubusercontent.com/arr-ai/arrai/v0.160.0/examples/import/bar.arrai}",
			},
			map[string]string{
				syntax.BundleConfig: ConfigFile(
					"github.com/test/test",
					ModuleFile("/github.com/test/test/test.arrai"),
				),
				ModuleFile(
					"/github.com/test/test/test.arrai",
				): "//{https://raw.githubusercontent.com/arr-ai/arrai/v0.160.0/examples/import/bar.arrai}",
				ModuleFile(
					"raw.githubusercontent.com/arr-ai/arrai/v0.160.0/examples/import/bar.arrai",
				): "1\n",
			},
		},
		{
			"no root", "/github.com/test/test/test.arrai",
			map[string]string{
				"/github.com/test/test/test.arrai": "1",
			},
			map[string]string{
				NoModuleFile("/test.arrai"): "1",
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			result := MustCreateTestBundleFromMap(t, c.files, syntax.MustAbs(t, c.path))
			ctxfs.ZipEqualToFiles(t, result, c.expectedFiles)
		})
	}
}

func TestBundleCompiledPlanRunsWithoutParse(t *testing.T) {
	t.Parallel()
	ctx := arraictx.InitRunCtx(context.Background())
	path := filepath.Join(t.TempDir(), "add.arrai")
	require.NoError(t, os.WriteFile(path, []byte("1 + 2"), 0o644))
	var buf bytes.Buffer
	require.NoError(t, BundledScripts(ctx, path, &buf))
	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	require.NoError(t, err)
	var hasPlan bool
	for _, f := range zr.File {
		if f.Name == "plan.bin" || f.Name == "/plan.bin" {
			hasPlan = true
			break
		}
	}
	require.True(t, hasPlan, "bundle zip must contain plan.bin")
	runCtx, err := syntax.WithBundleRun(ctx, buf.Bytes())
	require.NoError(t, err)
	p, err := syntax.LoadCompiledPlan(runCtx)
	require.NoError(t, err)
	require.NotNil(t, p, "LoadCompiledPlan must find /plan.bin")
	v, err := syntax.EvaluateBundleCtx(ctx, buf.Bytes())
	require.NoError(t, err)
	assert.True(t, v.Equal(rel.NewNumber(3)), "%s", v)
}

// FIXME: test github module import, only works locally, unable to locate cached module in CI
// func TestDeepModuleImports(t *testing.T) {
// 	t.Parallel()

// 	layerFS := ctxfs.CreateTestMemMapFs(t, map[string]string{
// 		sentinelPath("/github.com/test/test"): "module github.com/test/test\n",
// 		"/github.com/test/test/test.arrai":    "//{github.com/arr-ai/arrai/examples/comb_import}",
// 	})
// 	fs := afero.NewCopyOnWriteFs(afero.NewOsFs(), layerFS)
// 	ctx := ctxfs.SourceFsOnto(context.Background(), fs)
// 	ctx = ctxrootcache.WithRootCache(ctx)
// 	buf := &bytes.Buffer{}
// 	assert.NoError(t, bundleFiles(ctx, syntax.MustAbs(t, "/github.com/test/test/test.arrai"), buf))
// 	ctxfs.ZipEqualToFiles(t, buf.Bytes(), map[string]string{
// 		syntax.BundleConfig: config(
// 			"github.com/test/test",
// 			moduleFile("/github.com/test/test/test.arrai"),
// 		),
// 		moduleFile("/github.com/test/test/test.arrai"): "//{github.com/arr-ai/arrai/examples/comb_import}",
// 		moduleFile(
// 			"/github.com/arr-ai/arrai/examples/comb_import.arrai",
// 		): "//{./module_import} + //{/examples/import/relative_import.arrai}\n",
// 		moduleFile(
// 			"/github.com/arr-ai/arrai/examples/relative_import.arrai",
// 		): "//{./bar}\n",
// 		moduleFile(
// 			"/github.com/arr-ai/arrai/examples/bar.arrai",
// 		): "1\n",
// 		moduleFile(
// 			"/github.com/arr-ai/arrai/examples/module_import.arrai",
// 		): "//{/examples/import/bar}\n",
// 	})
// }
