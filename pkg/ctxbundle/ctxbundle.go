package ctxbundle

import (
	"archive/zip"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/afero"
	"github.com/spf13/afero/zipfs"
)

type bundleMode int

type bundleFs struct {
	sync.RWMutex
	zip afero.Fs
}

const (
	bundleModeKey bundleMode = iota
	localDir                 = "local"
	remoteDir                = "remote"
)

func WithBundleFs(ctx context.Context, r *zip.Reader) context.Context {
	return context.WithValue(ctx, bundleModeKey, &bundleFs{zip: zipfs.New(r)})
}

func fromBundleFs(ctx context.Context) *bundleFs {
	return ctx.Value(bundleModeKey).(*bundleFs)
}

func isBundling(ctx context.Context) bool {
	return ctx.Value(bundleModeKey) != nil
}

func ZipLocalImport(ctx context.Context, filePath string, content []byte) error {
	if !isBundling(ctx) {
		return nil
	}
	return zipFile(fromBundleFs(ctx), filepath.Join(localDir, filePath), content)
}

func ZipRemoteImport(ctx context.Context, filePath string, content []byte) error {
	if !isBundling(ctx) {
		return nil
	}
	return zipFile(fromBundleFs(ctx), filepath.Join(remoteDir, filePath), content)
}

func zipFile(fs *bundleFs, filePath string, content []byte) error {
	fs.Lock()
	defer fs.Unlock()

	if err := fs.zip.MkdirAll(filepath.Dir(filePath), os.ModeDir); err != nil {
		return err
	}

	if _, err := fs.zip.Stat(filePath); err != nil && os.IsExist(err) {
		return fmt.Errorf("BundleFs: file already exists: %s", filePath)
	}

	f, err := fs.zip.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(content)
	return err
}
