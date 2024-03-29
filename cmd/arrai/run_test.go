//nolint:lll,dupl
package main

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"path/filepath"
	s "sync"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/bundle"
	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/pkg/ctxrootcache"
	"github.com/arr-ai/arrai/syntax"
)

var (
	loadImportFs s.Once
	importFs     afero.Fs
)

func getImportFs(t *testing.T) afero.Fs {
	//
	// path1
	// ├── ModuleRootSentinel
	// └── path2
	//     ├── path3
	//     │   └── import_from_root.arrai
	//     └── path4
	//         ├── ModuleRootSentinel
	//         ├── 1.arrai
	//         ├── import_from_same_dir_root.arrai
	//         └── path5
	//             ├── 2.arrai
	//             └── import_from_nested_root.arrai
	//
	loadImportFs.Do(func() {
		importFs = afero.NewMemMapFs()
		require.NoError(t, importFs.MkdirAll("/path1/path2/path3", os.ModeDir))
		require.NoError(t, importFs.MkdirAll("/path1/path2/path4", os.ModeDir))
		files := []struct {
			fileName, expr string
		}{
			{
				filepath.Join("/path1", syntax.ModuleRootSentinel),
				"module github.com/test/path1\n",
			},
			{
				filepath.Join("/path1/path2/path4", syntax.ModuleRootSentinel),
				"module github.com/test/path1/path2/path4\n",
			},
			{"/path1/path2/path3/import_from_root.arrai", "//{/path2/path4/path5/2}"},
			{"/path1/path2/path4/1.arrai", "1"},
			{"/path1/path2/path4/import_from_same_dir_root.arrai", "//{/1}"},
			{"/path1/path2/path4/path5/2.arrai", "2"},
			{"/path1/path2/path4/path5/import_from_nested_root.arrai", "//{/1}"},
		}
		for _, af := range files {
			f, err := importFs.Create(syntax.MustAbs(t, af.fileName))
			require.NoError(t, err)
			defer f.Close()
			mustWrite(t, f, []byte(af.expr))
		}
	})
	return importFs
}

func TestRunBundle(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name, filePath, result string
	}{
		{"single file", "/path1/path2/path4/1.arrai", "1"},
		{"single root", "/path1/path2/path4/path5/import_from_nested_root.arrai", "1"},
		{"nested roots", "/path1/path2/path3/import_from_root.arrai", "2"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			ctx := ctxfs.SourceFsOnto(context.Background(), getImportFs(t))
			ctx = ctxrootcache.WithRootCache(ctx)
			zipped := &bytes.Buffer{}
			require.NoError(t, bundle.BundledScripts(ctx, syntax.MustAbs(t, c.filePath), zipped))
			out := &bytes.Buffer{}
			assert.NoError(t, runBundled(ctx, zipped.Bytes(), out, ""))
			assert.Equal(t, c.result+"\n", out.String())
		})
	}
}

func TestRunBundleOsArgs(t *testing.T) {
	t.Parallel()

	files := map[string]string{
		bundle.SentinelFile("github.com/args"):          "module github.com/args\n",
		bundle.ModuleFile("github.com/args/args.arrai"): "//os.args",
		syntax.BundleConfig:                             bundle.ConfigFile("github.com/args", bundle.ModuleFile("github.com/args/args.arrai")),
	}

	buf := createBundle(t, files)

	ctx := arraictx.InitRunCtx(context.Background())
	ctx = arraictx.WithArgs(ctx, "1", "2", "3")
	actual := &bytes.Buffer{}
	assert.NoError(t, runBundled(ctx, buf, actual, ""))
	assert.Equal(t, "['1', '2', '3']\n", actual.String())
}

func TestRunBundleWithHttp(t *testing.T) {
	t.Parallel()

	files := map[string]string{
		bundle.SentinelFile("github.com/test1"):          "module github.com/test1\n",
		bundle.ModuleFile("github.com/test1/test.arrai"): "//{https://raw.githubusercontent.com/arr-ai/arrai/v0.160.0/examples/import/bar.arrai}",
		bundle.ModuleFile(
			"raw.githubusercontent.com/arr-ai/arrai/v0.160.0/examples/import/bar.arrai",
		): "1",
		syntax.BundleConfig: bundle.ConfigFile("github.com/test1", bundle.ModuleFile("github.com/test1/test.arrai")),
	}

	buf := createBundle(t, files)

	ctx := arraictx.InitRunCtx(context.Background())
	actual := &bytes.Buffer{}
	assert.NoError(t, runBundled(ctx, buf, actual, ""))
	assert.Equal(t, "1\n", actual.String())
}

func TestRunBundleWithGithubImport(t *testing.T) {
	t.Parallel()

	files := map[string]string{
		bundle.SentinelFile("github.com/test1"):          "module github.com/test1\n",
		bundle.ModuleFile("github.com/test1/test.arrai"): "//{github.com/test2/test.arrai}",
		bundle.ModuleFile("github.com/test2/test.arrai"): "1",
		syntax.BundleConfig:                              bundle.ConfigFile("github.com/test1", bundle.ModuleFile("github.com/test1/test.arrai")),
	}

	buf := createBundle(t, files)

	ctx := arraictx.InitRunCtx(context.Background())
	actual := &bytes.Buffer{}
	assert.NoError(t, runBundled(ctx, buf, actual, ""))
	assert.Equal(t, "1\n", actual.String())
}

func createBundle(t *testing.T, files map[string]string) []byte {
	buf := &bytes.Buffer{}
	w := zip.NewWriter(buf)
	for name, content := range files {
		f, err := w.Create(name)
		require.NoError(t, err)
		_, err = f.Write([]byte(content))
		require.NoError(t, err)
	}
	w.Close()
	return buf.Bytes()
}

func TestRunBundleWithoutRoot(t *testing.T) {
	t.Parallel()
	fs := afero.NewMemMapFs()
	path := "very/deep/dir/test.arrai"
	f, err := fs.Create(path)
	require.NoError(t, err)
	expectedOut := "1\n"
	_, err = f.Write([]byte(expectedOut))
	require.NoError(t, err)
	f.Close()

	ctx := ctxfs.SourceFsOnto(context.Background(), fs)
	ctx = ctxrootcache.WithRootCache(ctx)
	zipped := &bytes.Buffer{}
	require.NoError(t, bundle.BundledScripts(ctx, path, zipped))
	out := &bytes.Buffer{}
	assert.NoError(t, runBundled(ctx, zipped.Bytes(), out, ""))
	assert.Equal(t, expectedOut, out.String())
}

func TestModuleImportRoot(t *testing.T) {
	t.Parallel()

	importTestFs := getImportFs(t)
	ctx := ctxfs.SourceFsOnto(context.Background(), importTestFs)
	cases := []struct {
		filePath, expected string
	}{
		{"/path1/path2/path3/import_from_root.arrai", "2"},
		{"/path1/path2/path4/import_from_same_dir_root.arrai", "1"},
		{"/path1/path2/path4/path5/import_from_nested_root.arrai", "1"},
	}
	for _, c := range cases {
		var buf bytes.Buffer
		require.NoError(t, evalFile(ctxrootcache.WithRootCache(ctx), syntax.MustAbs(t, c.filePath), &buf, ""))
		require.Equal(t, c.expected+"\n", buf.String())
	}
}

func TestNoImportRoot(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	require.NoError(t, fs.MkdirAll(syntax.MustAbs(t, "path/to/file"), os.ModeDir))
	f, err := fs.Create(syntax.MustAbs(t, "path/to/file/test.arrai"))
	require.NoError(t, err)
	defer f.Close()
	mustWrite(t, f, []byte("//{/file}"))

	f, err = fs.Create(syntax.MustAbs(t, "file.arrai"))
	require.NoError(t, err)
	defer f.Close()
	mustWrite(t, f, []byte("1"))
	err = evalFile(
		ctxrootcache.WithRootCache(ctxfs.SourceFsOnto(context.Background(), fs)),
		syntax.MustAbs(t, "path/to/file/test.arrai"), &bytes.Buffer{}, "",
	)
	require.Error(t, err)
	//FIXME: the error contains stack trace that is specific to local machine due to module imports
	assert.Contains(t, err.Error(), "module root not found")
}

func TestLocalImportErrors(t *testing.T) {
	t.Parallel()
	// syntax error in let expression
	fs := ctxfs.CreateTestMemMapFs(t, map[string]string{
		"a.arrai": "let x = //{./b}; x",
		"b.arrai": "//{./c}",
		"c.arrai": "123string",
	})
	err := evalFile(
		ctxfs.SourceFsOnto(context.Background(), fs),
		"a.arrai", &bytes.Buffer{}, "",
	)
	assert.EqualError(t, err, "unconsumed input\n \n\x1b[1;37mc.arrai:1:4:\x1b[0m\n123\x1b[1;31mstring\x1b[0m")

	// file not found error
	fs = ctxfs.CreateTestMemMapFs(t, map[string]string{
		"a.arrai": "let x = //{./b}; x",
	})
	err = evalFile(
		ctxfs.SourceFsOnto(context.Background(), fs),
		"a.arrai", &bytes.Buffer{}, "",
	)
	require.Error(t, err)
	// error contains local directory
	assert.Contains(t,
		err.Error(),
		"file does not exist\n\n\x1b[1;37ma.arrai:1:11:\x1b[0m\nlet x = //\x1b[1;31m{./b}\x1b[0m; x",
	)
}

// FIXME: remote import cannot be tested as it requires downloading modules into cache and the memory filesystem
// does not have access to the cache.
// func TestRemoteImportErrors(t *testing.T) {
// 	t.Parallel()
// 	fs := ctxfs.CreateTestMemMapFs(t, map[string]string{
// 		"a.arrai": "//{github.com/nofun97/test-arrai/syntax-error}",
// 		"b.arrai": "//{github.com/nofun97/test-arrai/wrong-import}",
// 		"c.arrai": "//{github.com/nofun97/test-arrai/import-syntax-error}",
// 	})
// 	ctx := ctxfs.SourceFsOnto(context.Background(), fs)
// 	assert.EqualError(t,
// 		evalFile(ctx, "a.arrai", &bytes.Buffer{}, ""),
// 		"",
// 	)
// 	assert.EqualError(t,
// 		evalFile(ctx, "b.arrai", &bytes.Buffer{}, ""),
// 		"",
// 	)
// 	assert.EqualError(t,
// 		evalFile(ctx, "c.arrai", &bytes.Buffer{}, ""),
// 		"",
// 	)
// }

func mustWrite(t *testing.T, f afero.File, content []byte) {
	_, err := f.Write(content)
	require.NoError(t, err)
}

func TestEvalFile(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	ctx := arraictx.InitRunCtx(context.Background())
	require.NoError(t, evalFile(ctx, "../../examples/jsfuncs/jsfuncs.arrai", &buf, ""))
	require.NoError(t, evalFile(ctx, "../../examples/grpc/app.arrai", &buf, ""))
}

func TestEvalNotExistingFile(t *testing.T) {
	t.Parallel()
	ctx := arraictx.InitRunCtx(context.Background())
	require.Equal(t, `"version": not a command and not found as a file in the current directory`,
		evalFile(ctx, "version", nil, "").Error())

	require.Equal(t, `"`+string([]rune{'.', os.PathSeparator})+`version": file not found`,
		evalFile(ctx, string([]rune{'.', os.PathSeparator})+"version", nil, "").Error())
}
