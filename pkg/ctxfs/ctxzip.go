package ctxfs

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/afero"
)

type zipFs struct {
	sync.RWMutex
	fs afero.Fs
}

func WithZipFs(ctx context.Context, key interface{}) context.Context {
	return context.WithValue(ctx, key, &zipFs{fs: afero.NewMemMapFs()})
}

func ZipFile(ctx context.Context, key interface{}, filePath string, content []byte) error {
	if !strings.HasPrefix(filePath, "/") {
		return fmt.Errorf("path has to be absolute in UNIX path format: %s", filePath)
	}

	fs := ctx.Value(key).(*zipFs)
	fs.Lock()
	defer fs.Unlock()

	if err := fs.fs.MkdirAll(filepath.Dir(filePath), os.ModeDir); err != nil {
		return err
	}

	if _, err := fs.fs.Stat(filePath); err != nil && os.IsExist(err) {
		return fmt.Errorf("BundleFs: file already exists: %s", filePath)
	}

	f, err := fs.fs.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(content)
	return err
}

func FileExists(ctx context.Context, key interface{}, filePath string) (bool, error) {
	fs := ctx.Value(key).(*zipFs)
	fs.Lock()
	defer fs.Unlock()

	fi, err := fs.fs.Stat(filePath)
	if err != nil {
		return false, err
	}
	return !fi.IsDir(), nil
}

func OutputZip(ctx context.Context, key interface{}, w io.Writer) error {
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()
	fs := ctx.Value(key).(*zipFs)
	fs.Lock()
	defer fs.Unlock()
	return afero.Walk(fs.fs, "/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := fs.fs.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		content, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		path = strings.TrimPrefix(path, "/")
		zipF, err := zipWriter.Create(path)
		if err != nil {
			return err
		}
		_, err = zipF.Write(content)
		return err
	})
}
