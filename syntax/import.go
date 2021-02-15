package syntax

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/anz-bank/pkg/mod"
	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/pkg/ctxrootcache"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
	"github.com/arr-ai/wbnf/parser"
	"github.com/sirupsen/logrus"
)

// ModuleRootSentinel is a file which marks the module root of a project.
const ModuleRootSentinel = "go.mod"

var (
	implicitDecoderSyncOnce sync.Once
	implicitDecode          rel.Value

	cache             = newCache()
	errModuleNotExist = errors.New("module root not found")
)

func implicitDecoder() rel.Value {
	implicitDecoderSyncOnce.Do(func() {
		implicitDecode = mustParseLit(string(MustAsset("syntax/implicit_import.arrai")))
	})
	return implicitDecode
}

func importLocalFile(
	ctx context.Context,
	decoder rel.Tuple,
	fromRoot bool,
	importPath, sourceDir string,
) (rel.Expr, error) {
	if fromRoot {
		rootPath, err := findRootFromModule(ctx, sourceDir)
		if err != nil {
			return nil, err
		}
		if err = addModuleSentinel(ctx, rootPath); err != nil {
			return nil, err
		}
		if !strings.HasPrefix(importPath, "/") {
			importPath = rootPath + "/" + strings.ReplaceAll(importPath, "../", "")
		}
	}

	if err := bundleLocalFile(ctx, importPath); err != nil {
		return nil, err
	}

	v, err := fileValue(ctx, decoder, importPath)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func importExternalContent(ctx context.Context, decoder rel.Tuple, importPath string) (rel.Expr, error) {
	var moduleErr error
	if !strings.HasPrefix(importPath, "http://") && !strings.HasPrefix(importPath, "https://") {
		v, err := importModuleFile(ctx, decoder, importPath)
		if err == nil {
			return v, nil
		}
		moduleErr = err

		// Since an explicit schema is allowed, it's OK to assume https as the default.
		importPath = "https://" + importPath
	}

	v, err := importURL(ctx, decoder, importPath)
	if err != nil {
		if moduleErr != nil {
			return nil, fmt.Errorf("failed to import %s - %s, and %s", importPath, moduleErr.Error(), err.Error())
		}
		return nil, err
	}

	return v, nil
}

func importModuleFile(ctx context.Context, decoder rel.Tuple, importPath string) (rel.Expr, error) {
	if isRunningBundle(ctx) {
		return fileValue(ctx, decoder, path.Join(ModuleDir, importPath))
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if err := mod.Config(mod.GoModulesMode,
		mod.GoModulesOptions{ModName: filepath.Base(wd)}, mod.GitHubOptions{}); err != nil {
		return nil, err
	}

	path, ver := mod.ExtractVersion(importPath)
	if ver != "" {
		return nil, errors.New("per-importing versioning is not allowed")
	}

	m, err := mod.Retrieve(path, ver)
	if err != nil {
		return nil, err
	}

	relImportPath := strings.TrimPrefix(importPath, m.Name)
	if ctx, err = bundleModule(ctx, relImportPath, m); err != nil {
		return nil, err
	}

	return fileValue(ctx, decoder, filepath.Join(m.Dir, relImportPath))
}

func findRootFromModule(ctx context.Context, modulePath string) (string, error) {
	currentPath, err := filepath.Abs(modulePath)
	if err != nil {
		return "", err
	}

	currentPath = bundleToValidPath(ctx, currentPath)

	if r, exists, err := ctxrootcache.LoadRoot(ctx, currentPath); err != nil {
		return "", err
	} else if exists {
		return r, nil
	}

	systemRoot, err := filepath.Abs(string(os.PathSeparator))
	if err != nil {
		return "", err
	}

	systemRoot = bundleToValidPath(ctx, systemRoot)

	// 16 is enough for pretty much all cases.
	paths := append(make([]string, 0, 16), currentPath)

	// Keep walking up the directories to find nearest root marker
	for {
		exists, err := tools.FileExists(ctx, filepath.Join(currentPath, ModuleRootSentinel))
		reachedRoot := currentPath == systemRoot || (err != nil && os.IsPermission(err))
		switch {
		case exists:
			for _, p := range paths {
				if err := ctxrootcache.StoreRoot(ctx, p, currentPath); err != nil {
					return "", err
				}
			}
			return currentPath, nil
		case reachedRoot:
			return "", errModuleNotExist
		case err != nil:
			return "", err
		}
		currentPath = filepath.Dir(currentPath)
		paths = append(paths, currentPath)
	}
}

func importURL(ctx context.Context, decoder rel.Tuple, url string) (rel.Expr, error) {
	if isRunningBundle(ctx) {
		url = strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "http://")
		return fileValue(ctx, decoder, path.Join(ModuleDir, url))
	}
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if err = bundleRemoteFile(ctx, url, data); err != nil {
			return nil, err
		}
		val, err := cache.getOrAdd(url, func() (rel.Expr, error) { return bytesValue(ctx, NoPath, data) })
		return val, err
	}
	return nil, fmt.Errorf("request %s failed: %s", url, resp.Status)
}

func fileValue(ctx context.Context, decoder rel.Tuple, filename string) (rel.Expr, error) {
	if filepath.Ext(filename) == "" {
		filename += ".arrai"
	}

	bytes, err := ctxfs.ReadFile(ctxfs.SourceFsFrom(ctx), filename)
	if err != nil {
		return nil, err
	}

	// TODO: add cache of decoded values
	if decoder != nil {
		return decode(ctx, decoder, bytes)
	}
	if ext := filepath.Ext(filename); ext != ".arrai" {
		decoded, err := rel.NewCallExprCurry(
			parser.Scanner{},
			implicitDecoder(),
			rel.NewString([]rune(ext)),
			rel.NewBytes(bytes),
		).Eval(ctx, StdScope())
		if err != nil {
			return nil, err
		}

		if decoded.Equal(rel.EmptyTuple) {
			logrus.Debugf("implicit decoding failed for %s. Importing with bytes decoder.", filename)
			return rel.NewBytes(bytes), nil
		}
		return decoded, nil
	}

	return bytesValue(ctx, filename, bytes)
}

func bytesValue(ctx context.Context, filename string, data []byte) (rel.Expr, error) {
	compile := func() (rel.Expr, error) {
		return Compile(ctx, filename, string(data))
	}
	if filename != NoPath {
		return cache.getOrAdd(filename, compile)
	}
	return compile()
}
