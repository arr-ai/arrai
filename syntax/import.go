package syntax

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
	"github.com/arr-ai/arrai/tools/module"
	"github.com/arr-ai/arrai/translate"
)

const arraiRootMarker = "go.mod"

var cache = newCache()

func importLocalFile(ctx context.Context, fromRoot bool, importPath string) (rel.Expr, error) {
	if fromRoot {
		rootPath, err := findRootFromModule(filepath.Dir(importPath))
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(importPath, "/") {
			importPath = rootPath + "/" + strings.ReplaceAll(importPath, "../", "")
		}
	}

	v, err := fileValue(ctx, importPath)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func importExternalContent(ctx context.Context, importPath string) (rel.Expr, error) {
	var moduleErr error
	if !strings.HasPrefix(importPath, "http://") && !strings.HasPrefix(importPath, "https://") {
		v, err := importModuleFile(ctx, importPath)
		if err == nil {
			return v, nil
		}
		moduleErr = err

		// Since an explicit schema is allowed, it's OK to assume https as the default.
		importPath = "https://" + importPath
	}

	v, err := importURL(ctx, importPath)
	if err != nil {
		if moduleErr != nil {
			return nil, fmt.Errorf("failed to import %s - %s, and %s", importPath, moduleErr.Error(), err.Error())
		}
		return nil, err
	}

	return v, nil
}

func importModuleFile(ctx context.Context, importPath string) (rel.Expr, error) {
	var mod module.Module = module.NewGoModule()

	m, err := mod.Get(importPath)
	if err != nil {
		return nil, err
	}

	return fileValue(ctx, filepath.Join(m.Dir, strings.TrimPrefix(importPath, m.Name)))
}

func findRootFromModule(modulePath string) (string, error) {
	currentPath, err := filepath.Abs(modulePath)
	if err != nil {
		return "", err
	}

	systemRoot, err := filepath.Abs(string(os.PathSeparator))
	if err != nil {
		return "", err
	}

	// Keep walking up the directories to find nearest root marker
	for {
		exists := tools.FileExists(filepath.Join(currentPath, arraiRootMarker))
		reachedRoot := currentPath == systemRoot || (err != nil && os.IsPermission(err))
		switch {
		case exists:
			return currentPath, nil
		case reachedRoot:
			//TODO: test this after context filesystem is implemented
			return "", errors.New("module root not found")
		case err != nil:
			return "", err
		}
		currentPath = filepath.Dir(currentPath)
	}
}

func importURL(ctx context.Context, url string) (rel.Expr, error) {
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
		val, err := cache.getOrAdd(url, func() (rel.Expr, error) { return bytesValue(ctx, NoPath, data) })
		return val, err
	}
	return nil, fmt.Errorf("request %s failed: %s", url, resp.Status)
}

func fileValue(ctx context.Context, filename string) (rel.Expr, error) {
	if filepath.Ext(filename) == "" {
		filename += ".arrai"
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	switch filepath.Ext(filename) {
	case ".json":
		return bytesJSONToArrai(bytes)
	case ".yml", ".yaml":
		return translate.BytesYamlToArrai(bytes)
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
