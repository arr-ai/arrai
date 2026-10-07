// nolint: lll
package bundle

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
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

func TestBundleZipIsByteIdenticalOnRepeat(t *testing.T) {
	t.Parallel()
	files := map[string]string{
		SentinelPath("/github.com/test/test"):               "module github.com/test/test\n",
		"/github.com/test/test/test.arrai":                  "//{./module/module2/module.arrai}",
		SentinelPath("/github.com/test/test/module/"):       "module github.com/test/test/module\n",
		"/github.com/test/test/module/1.arrai":              "1",
		"/github.com/test/test/module/module2/module.arrai": "//{/1.arrai}",
	}
	path := syntax.MustAbs(t, "/github.com/test/test/test.arrai")
	a := MustCreateTestBundleFromMap(t, files, path)
	b := MustCreateTestBundleFromMap(t, files, path)
	require.Equal(t, a, b)
}

func TestBundlePlanImportPathsArePortable(t *testing.T) {
	t.Parallel()
	files := map[string]string{
		SentinelPath("/github.com/test/test"):               "module github.com/test/test\n",
		"/github.com/test/test/test.arrai":                  "//{./module/module2/module.arrai}",
		SentinelPath("/github.com/test/test/module/"):       "module github.com/test/test/module\n",
		"/github.com/test/test/module/1.arrai":              "1",
		"/github.com/test/test/module/module2/module.arrai": "//{/1.arrai}",
	}
	zipBytes := MustCreateTestBundleFromMap(t, files, syntax.MustAbs(t, "/github.com/test/test/test.arrai"))
	ctx, err := syntax.WithBundleRun(arraictx.InitRunCtx(context.Background()), zipBytes)
	require.NoError(t, err)
	p, err := syntax.LoadCompiledPlan(ctx)
	require.NoError(t, err)
	require.NotNil(t, p)
	paths := collectImportPaths(p.Root)
	require.NotEmpty(t, paths, "nested local imports must appear in the compiled plan")
	var sawModule bool
	for _, path := range paths {
		if strings.HasPrefix(path, syntax.ModuleDir+"/") || strings.HasPrefix(path, syntax.NoModuleDir+"/") {
			sawModule = sawModule || strings.HasPrefix(path, syntax.ModuleDir+"/")
			continue
		}
		require.False(t, filepath.IsAbs(path), "host-absolute import path %q leaked into plan.bin", path)
	}
	require.True(t, sawModule, "relative import //{./...} must lower to a /module/... path, got %q", paths)
}

func collectImportPaths(n rel.PlanNode) []string {
	var out []string
	var walk func(rel.PlanNode)
	walk = func(n rel.PlanNode) {
		if n.K == "import" {
			out = append(out, n.Str)
		}
		for _, k := range n.Kids {
			walk(k)
		}
	}
	walk(n)
	return out
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
	names := make([]string, 0, len(zr.File))
	for _, f := range zr.File {
		names = append(names, f.Name)
		if f.Name == "plan.bin" || f.Name == "/plan.bin" {
			hasPlan = true
			break
		}
	}
	require.True(t, hasPlan, "bundle zip must contain plan.bin; got %q", names)
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

// 🎯T45 (arr-ai/arrai#779): a macro's @transform result is evaluated at parse
// time and embedded in the plan as a literal. A native function cannot be
// embedded (its name does not identify it), so bundling must fail with a
// clear error rather than write a bundle that runs some other native.
func TestBundleMacroReturningNativeFails(t *testing.T) {
	t.Parallel()
	const src = `let g = {://grammar.lang.wbnf: doc -> /{x}; :};
let mk = (@grammar: g, @transform: (doc: \ast //encoding.json.decode));
let f = {:mk:x:};
f('{"a":1}')`
	ctx := arraictx.InitRunCtx(context.Background())
	path := filepath.Join(t.TempDir(), "m.arrai")
	require.NoError(t, os.WriteFile(path, []byte(src), 0o644))
	var buf bytes.Buffer
	err := BundledScripts(ctx, path, &buf)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot bundle native function ⦑decode⦒")
	assert.Contains(t, err.Error(), "arr-ai/arrai#779")
}

// Transforms that return a closure or a plain value still bundle, and the
// bundle (run from plan.bin, not re-parsed source) evaluates to the same
// value as the source. The closure cases cover: calling a native through a
// closure (the #779 workaround), two macros whose closures call natives that
// share a name (json and csv `decode`, the exact collision behind #779), and a
// closure capturing a let-bound value.
func TestBundleMacroReturningClosureOrValueRunsIdentically(t *testing.T) {
	t.Parallel()
	const grammar = `let g = {://grammar.lang.wbnf: doc -> /{x}; :};
`
	cases := map[string]string{
		"closure calling native": grammar + `
let mk = (@grammar: g, @transform: (doc: \ast \s //encoding.json.decode(s)));
let f = {:mk:x:};
f('{"a":1}')`,
		"two closures calling same-named natives": grammar + `
let mkj = (@grammar: g, @transform: (doc: \ast \s //encoding.json.decode(s)));
let mkc = (@grammar: g, @transform: (doc: \ast \s //encoding.csv.decode(s)));
let fj = {:mkj:x:};
let fc = {:mkc:x:};
[fj('{"a":1}'), fc('p,q')]`,
		"closure capturing a let binding": grammar + `
let k = 10;
let mk = (@grammar: g, @transform: (doc: \ast \n n + k));
let f = {:mk:x:};
f(5)`,
		"plain value": grammar + `
let mk = (@grammar: g, @transform: (doc: \ast (hello: 42, items: [1, 2, 3])));
let v = {:mk:x:};
v.items(1)`,
	}
	for name, src := range cases {
		src := src
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := arraictx.InitRunCtx(context.Background())
			path := filepath.Join(t.TempDir(), "m.arrai")
			require.NoError(t, os.WriteFile(path, []byte(src), 0o644))

			want, err := syntax.EvaluateExpr(ctx, path, src)
			require.NoError(t, err)

			var buf bytes.Buffer
			require.NoError(t, BundledScripts(ctx, path, &buf))
			runCtx, err := syntax.WithBundleRun(ctx, buf.Bytes())
			require.NoError(t, err)
			p, err := syntax.LoadCompiledPlan(runCtx)
			require.NoError(t, err)
			require.NotNil(t, p, "bundle must carry plan.bin; a source fallback would hide the bug")

			got, err := syntax.EvaluateBundleCtx(ctx, buf.Bytes())
			require.NoError(t, err)
			assert.True(t, want.Equal(got), "source=%s bundle=%s", want, got)
		})
	}
}
