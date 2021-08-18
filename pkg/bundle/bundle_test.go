//nolint: lll
package bundle

import (
	"testing"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/syntax"
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
