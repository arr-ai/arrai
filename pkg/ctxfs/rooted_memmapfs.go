package ctxfs

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/afero"
)

// rootedMemMapFs wraps an afero.MemMapFs and normalizes every path before
// delegating, so it resolves consistently regardless of how the path was
// spelled. afero.MemMapFs (still true as of v1.15.0, the latest release)
// treats "name" and "/name" as distinct, unrelated entries rather than
// aliases of the same file, so code that builds a bare relative path (as
// path.Join(".", name) does) and expects it to nest under root — which is
// how a real OS filesystem behaves — silently fails to find it. On Windows
// it additionally can't resolve a path carrying a volume name (`C:\foo`),
// a bug open upstream since 2019: https://github.com/spf13/afero/issues/225.
// A real OS filesystem doesn't have either problem, so this only wraps
// MemMapFs.
type rootedMemMapFs struct {
	afero.Fs
}

// NewTestMemMapFs returns a fresh afero.MemMapFs made safe to address with
// either a bare relative path or a native absolute path (as filepath.Abs
// produces, including on Windows). Use this instead of afero.NewMemMapFs()
// for any test fixture that will be read through arrai's normal
// (non-bundling) file-resolution code.
func NewTestMemMapFs() afero.Fs {
	return rootedMemMapFs{afero.NewMemMapFs()}
}

// normalizeForMemMapFs strips any Windows volume name and backslashes, then
// roots the result at "/" if it wasn't already absolute.
func normalizeForMemMapFs(p string) string {
	p = ToUnixPath(p)
	switch {
	case p == "" || p == ".":
		return "/"
	case !strings.HasPrefix(p, "/"):
		return "/" + p
	default:
		return p
	}
}

func (fs rootedMemMapFs) Create(name string) (afero.File, error) {
	return fs.Fs.Create(normalizeForMemMapFs(name))
}

func (fs rootedMemMapFs) Mkdir(name string, perm os.FileMode) error {
	return fs.Fs.Mkdir(normalizeForMemMapFs(name), perm)
}

func (fs rootedMemMapFs) MkdirAll(path string, perm os.FileMode) error {
	return fs.Fs.MkdirAll(normalizeForMemMapFs(path), perm)
}

func (fs rootedMemMapFs) Open(name string) (afero.File, error) {
	return fs.Fs.Open(normalizeForMemMapFs(name))
}

func (fs rootedMemMapFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return fs.Fs.OpenFile(normalizeForMemMapFs(name), flag, perm)
}

func (fs rootedMemMapFs) Remove(name string) error {
	return fs.Fs.Remove(normalizeForMemMapFs(name))
}

func (fs rootedMemMapFs) RemoveAll(path string) error {
	return fs.Fs.RemoveAll(normalizeForMemMapFs(path))
}

func (fs rootedMemMapFs) Rename(oldname, newname string) error {
	return fs.Fs.Rename(normalizeForMemMapFs(oldname), normalizeForMemMapFs(newname))
}

func (fs rootedMemMapFs) Stat(name string) (os.FileInfo, error) {
	return fs.Fs.Stat(normalizeForMemMapFs(name))
}

func (fs rootedMemMapFs) Chmod(name string, mode os.FileMode) error {
	return fs.Fs.Chmod(normalizeForMemMapFs(name), mode)
}

func (fs rootedMemMapFs) Chown(name string, uid, gid int) error {
	return fs.Fs.Chown(normalizeForMemMapFs(name), uid, gid)
}

func (fs rootedMemMapFs) Chtimes(name string, atime, mtime time.Time) error {
	return fs.Fs.Chtimes(normalizeForMemMapFs(name), atime, mtime)
}
