package ctxfs

import (
	"context"
	"io/ioutil"

	"github.com/spf13/afero"
)

type ctxkey int

const (
	sourceFsKey ctxkey = iota
	runtimeFsKey
)

var defaultFs afero.Fs

// Sets the default Fs to be returned from SourceFsFrom and RuntimeFsFrom.
func SetDefaultFs(fs afero.Fs) {
	defaultFs = fs
}

// SourceFsOnto returns a new Context with sourceFs added to ctx.
func SourceFsOnto(ctx context.Context, sourceFs afero.Fs) context.Context {
	return context.WithValue(ctx, sourceFsKey, sourceFs)
}

// SourceFsFrom extracts the filesystem used for import resolution from ctx.
func SourceFsFrom(ctx context.Context) afero.Fs {
	if v := ctx.Value(sourceFsKey); v != nil {
		return v.(afero.Fs)
	}
	return defaultFs
}

// RuntimeFsOnto returns a new Context with runtimeFs added to ctx.
func RuntimeFsOnto(ctx context.Context, runtimeFs afero.Fs) context.Context {
	return context.WithValue(ctx, runtimeFsKey, runtimeFs)
}

// RuntimeFsFrom extracts the filesystem used to access the filesystem from ctx.
// It will be used by //os.file and similar functions.
// TODO: Actually use it.
func RuntimeFsFrom(ctx context.Context) afero.Fs {
	if v := ctx.Value(runtimeFsKey); v != nil {
		return v.(afero.Fs)
	}
	return defaultFs
}

func ReadFile(fs afero.Fs, filePath string) ([]byte, error) {
	f, err := fs.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}
