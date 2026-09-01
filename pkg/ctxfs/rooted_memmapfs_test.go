package ctxfs

import (
	"os"
	"runtime"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestNormalizeForMemMapFs(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"":                "/",
		".":               "/",
		"a.arrai":         "/a.arrai",
		"dir/to/file.txt": "/dir/to/file.txt",
		"/a.arrai":        "/a.arrai",
	}
	for in, want := range cases {
		require.Equal(t, want, normalizeForMemMapFs(in), "input=%q", in)
	}

	// Volume-name stripping (via ToUnixPath) only applies on windows; on
	// other platforms a "C:\..." string is just an unusual relative path.
	if runtime.GOOS == "windows" {
		require.Equal(t, "/a.arrai", normalizeForMemMapFs(`C:\a.arrai`))
	} else {
		require.Equal(t, `/C:\a.arrai`, normalizeForMemMapFs(`C:\a.arrai`))
	}
}

func TestNewTestMemMapFsRootsRelativePaths(t *testing.T) {
	t.Parallel()

	fs := NewTestMemMapFs()

	f, err := fs.Create("a.arrai")
	require.NoError(t, err)
	_, err = f.Write([]byte("1"))
	require.NoError(t, err)
	require.NoError(t, f.Close())

	// The same file must be reachable both root-relative and bare-relative,
	// since production code resolves imports as bare-relative paths while
	// test fixtures are often built from absolute ones.
	for _, name := range []string{"a.arrai", "/a.arrai"} {
		data, err := afero.ReadFile(fs, name)
		require.NoError(t, err, "name=%q", name)
		require.Equal(t, "1", string(data), "name=%q", name)
	}

	var found []string
	require.NoError(t, afero.Walk(fs, "/", func(path string, info os.FileInfo, err error) error {
		require.NoError(t, err)
		if !info.IsDir() {
			found = append(found, path)
		}
		return nil
	}))
	require.Equal(t, []string{"/a.arrai"}, found)
}

func TestNewTestMemMapFsMkdirAllThenCreate(t *testing.T) {
	t.Parallel()

	fs := NewTestMemMapFs()
	require.NoError(t, fs.MkdirAll("dir/to/replace", os.ModeDir))

	f, err := fs.Create("dir/to/replace/file.txt")
	require.NoError(t, err)
	require.NoError(t, f.Close())

	fi, err := fs.Stat("/dir/to/replace")
	require.NoError(t, err)
	require.True(t, fi.IsDir())
}
