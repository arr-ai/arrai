package syntax

import (
	"context"
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/arr-ai/arrai/pkg/ctxfs"
)

type bundleMode int

const (
	bundleModeKey bundleMode = iota
	moduleDir                = "/module"
	// noModuleDir = "unnamed"
)

var (
	rootModuleRE = regexp.MustCompile("^module ([^\n]+)\n")
)

func isBundling(ctx context.Context) bool {
	return ctx.Value(bundleModeKey) != nil
}

//TODO: write config for arrai run
func setupBundleMode(ctx context.Context, filePath string, source []byte) (context.Context, error) {
	ctx = ctxfs.WithZipFs(ctx, bundleModeKey)
	root, err := findRootFromModule(ctx, filePath)
	if err != nil {
		if err != errModuleNotExist {
			return ctx, err
		}
		//FIXME: might want to handle arrai with no module
		if err = ctxfs.ZipFile(ctx, bundleModeKey, path.Join(moduleDir, filepath.Base(filePath)), source); err != nil {
			return ctx, err
		}
	}
	f, err := ctxfs.SourceFsFrom(ctx).Open(filepath.Join(root, ModuleRootSentinel))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	moduleRoot := rootModuleRE.FindStringSubmatch(string(buf))
	if moduleRoot == nil {
		return ctx, errors.New("sentinel does not show module path")
	}

	err = ctxfs.ZipFile(ctx, bundleModeKey, path.Join(moduleDir, moduleRoot[1], ModuleRootSentinel), buf)
	if err != nil {
		return ctx, err
	}

	filePath, err = filepath.Abs(filePath)
	if err != nil {
		return ctx, err
	}

	filePathRelativeToSentinel := toUnixPath(strings.TrimPrefix(filePath, root))

	return ctx, ctxfs.ZipFile(ctx, path.Join(moduleDir, moduleRoot[1], filePathRelativeToSentinel), source)
}

//FIXME: handle nested root
func addLocalRoot(ctx context.Context, rootPath string, source []byte) error {
	if !isBundling(ctx) {
		return nil
	}
	if exists, err := ctxfs.FileExists(ctx, bundleModeKey, rootPath); err != nil {
		return err
	} else if !exists {
		return nil
	}

	return ctxfs.ZipFile(ctx, bundleModeKey, toUnixPath(rootPath), source)
}

func bundleLocalFile(ctx context.Context, filePath, rootPath string, source []byte) error {
	if !isBundling(ctx) {
		return nil
	}
	return nil
}

func bundleRemoteFile(ctx context.Context, url string, source []byte) error {
	if !isBundling(ctx) {
		return nil
	}
	return nil
}

func toUnixPath(p string) string {
	if runtime.GOOS == "windows" {
		if vol := filepath.VolumeName(p) + ":"; strings.HasPrefix(p, vol) {
			p = strings.TrimPrefix(p, vol)
		}
		return strings.ReplaceAll(p, "\\", "/")
	}
	return p
}
