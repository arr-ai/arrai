//nolint: lll
package main

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/pkg/ctxrootcache"
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
				sentinelPath("/github.com/test/test"): "module github.com/test/test\n",
				"/github.com/test/test/test.arrai":    "1",
			},
			map[string]string{
				syntax.BundleConfig: config(
					"github.com/test/test",
					moduleFile("/github.com/test/test/test.arrai"),
				),
				sentinelFile("/github.com/test/test"):          "module github.com/test/test\n",
				moduleFile("/github.com/test/test/test.arrai"): "1",
			},
		},
		{
			"local dependencies with nested root", "/github.com/test/test/test.arrai",
			map[string]string{
				sentinelPath("/github.com/test/test"):               "module github.com/test/test\n",
				"/github.com/test/test/test.arrai":                  "//{./module/module2/module.arrai}",
				sentinelPath("/github.com/test/test/module/"):       "module github.com/test/test/module\n",
				"/github.com/test/test/module/1.arrai":              "1",
				"/github.com/test/test/module/module2/module.arrai": "//{/1.arrai}",
			},
			map[string]string{
				syntax.BundleConfig: config(
					"github.com/test/test",
					moduleFile("/github.com/test/test/test.arrai"),
				),
				sentinelFile("/github.com/test/test"):                           "module github.com/test/test\n",
				sentinelFile("/github.com/test/test/module/"):                   "module github.com/test/test/module\n",
				moduleFile("/github.com/test/test/test.arrai"):                  "//{./module/module2/module.arrai}",
				moduleFile("/github.com/test/test/module/module2/module.arrai"): "//{/1.arrai}",
				moduleFile("/github.com/test/test/module/1.arrai"):              "1",
			},
		},
		{
			"remote import", "/github.com/test/test/test.arrai",
			map[string]string{
				sentinelPath("/github.com/test/test"): "module github.com/test/test\n",
				"/github.com/test/test/test.arrai":    "//{https://raw.githubusercontent.com/arr-ai/arrai/v0.160.0/examples/import/bar.arrai}",
			},
			map[string]string{
				syntax.BundleConfig: config(
					"github.com/test/test",
					moduleFile("/github.com/test/test/test.arrai"),
				),
				moduleFile(
					"/github.com/test/test/test.arrai",
				): "//{https://raw.githubusercontent.com/arr-ai/arrai/v0.160.0/examples/import/bar.arrai}",
				moduleFile(
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
				noModuleFile("/test.arrai"): "1",
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			fs := ctxfs.CreateTestMemMapFs(t, c.files)
			ctx := ctxfs.SourceFsOnto(context.Background(), fs)
			ctx = ctxrootcache.WithRootCache(ctx)
			buf := &bytes.Buffer{}
			assert.NoError(t, bundleFiles(ctx, syntax.MustAbs(t, c.path), buf))
			ctxfs.ZipEqualToFiles(t, buf.Bytes(), c.expectedFiles)
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

func config(mainRoot, mainFile string) string {
	return fmt.Sprintf("(main_root: %q, main_file: %q)", mainRoot, mainFile)
}

func moduleFile(file string) string {
	return path.Join(syntax.ModuleDir, file)
}

func noModuleFile(file string) string {
	return path.Join(syntax.NoModuleDir, file)
}

func sentinelFile(file string) string {
	return moduleFile(sentinelPath(file))
}

func sentinelPath(file string) string {
	return path.Join(file, syntax.ModuleRootSentinel)
}
