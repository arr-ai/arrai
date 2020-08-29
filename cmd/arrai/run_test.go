package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/syntax"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestModuleImportRoot(t *testing.T) {
	t.Parallel()

	//
	// path1
	// ├── arraiRootMarker
	// └── path2
	//     ├── path3
	//     │   └── import_from_root.arrai
	//     └── path4
	//         ├── arraiRootMarker
	//         ├── 1.arrai
	//         ├── import_from_same_dir_root.arrai
	//         └── path5
	//             ├── 2.arrai
	//             └── import_from_nested_root.arrai
	//
	importTestFs := afero.NewMemMapFs()
	require.NoError(t, importTestFs.MkdirAll(mustAbs(t, "path1/path2/path3"), os.ModeDir))
	require.NoError(t, importTestFs.MkdirAll(mustAbs(t, "path1/path2/path4"), os.ModeDir))
	files := []struct {
		fileName, expr string
	}{
		{filepath.Join("path1", syntax.ArraiRootMarker), ""},
		{filepath.Join("path1/path2/path4", syntax.ArraiRootMarker), ""},
		{"path1/path2/path3/import_from_root.arrai", "//{/path2/path4/path5/2}"},
		{"path1/path2/path4/1.arrai", "1"},
		{"path1/path2/path4/import_from_same_dir_root.arrai", "//{/1}"},
		{"path1/path2/path4/path5/2.arrai", "2"},
		{"path1/path2/path4/path5/import_from_nested_root.arrai", "//{/1}"},
	}
	for _, af := range files {
		f, err := importTestFs.Create(mustAbs(t, af.fileName))
		require.NoError(t, err)
		defer f.Close()
		mustWrite(t, f, []byte(af.expr))
	}

	ctx := ctxfs.SourceFsOnto(context.Background(), importTestFs)
	cases := []struct {
		filePath, expected string
	}{
		{"path1/path2/path3/import_from_root.arrai", "2"},
		{"path1/path2/path4/import_from_same_dir_root.arrai", "1"},
		{"path1/path2/path4/path5/import_from_nested_root.arrai", "1"},
	}
	for _, c := range cases {
		var buf bytes.Buffer
		require.NoError(t, evalFile(ctx, mustAbs(t, c.filePath), &buf, ""))
		require.Equal(t, c.expected+"\n", buf.String())
	}
}

func TestNoImportRoot(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	require.NoError(t, fs.MkdirAll(mustAbs(t, "path/to/file"), os.ModeDir))
	f, err := fs.Create(mustAbs(t, "path/to/file/test.arrai"))
	require.NoError(t, err)
	defer f.Close()
	mustWrite(t, f, []byte("//{/file}"))

	f, err = fs.Create(mustAbs(t, "file.arrai"))
	require.NoError(t, err)
	defer f.Close()
	mustWrite(t, f, []byte("1"))

	require.EqualError(t,
		evalFile(
			ctxfs.SourceFsOnto(context.Background(), fs),
			mustAbs(t, "path/to/file/test.arrai"), &bytes.Buffer{}, "",
		),
		"module root not found")
}

func mustAbs(t *testing.T, filePath string) string {
	abs, err := filepath.Abs(filePath)
	require.NoError(t, err)
	return abs
}

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
