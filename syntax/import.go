package syntax

import (
	"context"
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

func importLocalFile(ctx context.Context, fromRoot bool) rel.Value {
	return rel.NewNativeFunction("//./", func(v rel.Value) (rel.Value, error) {
		s, ok := rel.AsString(v.(rel.Set))
		if !ok {
			return nil, fmt.Errorf("cannot convert %#v to string", v)
		}

		filename := s.String()
		if fromRoot {
			pwd, err := os.Getwd()
			if err != nil {
				return nil, err
			}
			rootPath, err := findRootFromModule(pwd)
			if err != nil {
				return nil, err
			}
			if !strings.HasPrefix(filename, "/") {
				filename = rootPath + "/" + strings.ReplaceAll(filename, "../", "")
			}
		}

		v, err := fileValue(ctx, filename)
		if err != nil {
			return nil, err
		}

		return v, nil
	})
}

func importExternalContent(ctx context.Context) rel.Value {
	return rel.NewNativeFunction("//", func(v rel.Value) (rel.Value, error) {
		s, ok := rel.AsString(v.(rel.Set))
		if !ok {
			return nil, fmt.Errorf("cannot convert %#v to string", v)
		}
		importpath := s.String()

		var moduleErr error

		if !strings.HasPrefix(importpath, "http://") && !strings.HasPrefix(importpath, "https://") {
			v, err := importModuleFile(ctx, importpath)
			if err == nil {
				return v, nil
			}
			moduleErr = err

			// Since an explicit schema is allowed, it's OK to assume https as the default.
			importpath = "https://" + importpath
		}

		v, err := importURL(ctx, importpath)
		if err != nil {
			if moduleErr != nil {
				return nil, fmt.Errorf("failed to import %s - %s, and %s", importpath, moduleErr.Error(), err.Error())
			}
			return nil, err
		}

		return v, nil
	})
}

func importModuleFile(ctx context.Context, importpath string) (rel.Value, error) {
	var mod module.Module = module.NewGoModule()

	m, err := mod.Get(importpath)
	if err != nil {
		return nil, err
	}

	return fileValue(ctx, filepath.Join(m.Dir, strings.TrimPrefix(importpath, m.Name)))
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
			return "", nil
		case err != nil:
			return "", err
		}
		currentPath = filepath.Dir(currentPath)
	}
}

func importURL(ctx context.Context, url string) (rel.Value, error) {
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
		val, err := cache.getOrAdd(url, func() (rel.Value, error) { return bytesValue(ctx, NoPath, data) })
		return val, err
	}
	return nil, fmt.Errorf("request %s failed: %s", url, resp.Status)
}

func fileValue(ctx context.Context, filename string) (rel.Value, error) {
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

func bytesValue(ctx context.Context, filename string, data []byte) (rel.Value, error) {
	eval := func() (rel.Value, error) {
		expr, err := Compile(ctx, filename, string(data))
		if err != nil {
			return nil, err
		}
		value, err := expr.Eval(rel.EmptyScope)
		if err != nil {
			return nil, rel.WrapContext(err, expr, rel.EmptyScope)
		}
		return value, nil
	}
	if filename != NoPath {
		return cache.getOrAdd(filename, eval)
	}

	return eval()
}
