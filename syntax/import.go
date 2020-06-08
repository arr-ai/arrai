package syntax

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
	"github.com/arr-ai/arrai/tools/module"
	"github.com/arr-ai/arrai/translate"
)

const arraiRootMarker = "go.mod"

var importLocalFileOnce sync.Once
var importLocalFileVar rel.Value
var cache *importCache = newCache()

func importLocalFile(fromRoot bool) rel.Value {
	importLocalFileOnce.Do(func() {
		importLocalFileVar = rel.NewNativeFunction("//./", func(v rel.Value) (rel.Value, error) {
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

			v, err := fileValue(filename)
			if err != nil {
				return nil, err
			}

			return v, nil
		})
	})
	return importLocalFileVar
}

var importExternalContentOnce sync.Once
var importExternalContentVar rel.Value

func importExternalContent() rel.Value {
	importExternalContentOnce.Do(func() {
		importExternalContentVar = rel.NewNativeFunction("//", func(v rel.Value) (rel.Value, error) {
			s, ok := rel.AsString(v.(rel.Set))
			if !ok {
				return nil, fmt.Errorf("cannot convert %#v to string", v)
			}
			importpath := s.String()

			var moduleErr error

			if !strings.HasPrefix(importpath, "http://") && !strings.HasPrefix(importpath, "https://") {
				v, err := importModuleFile(importpath)
				if err == nil {
					return v, nil
				}
				moduleErr = err

				// Since an explicit schema is allowed, it's OK to assume https as the default.
				importpath = "https://" + importpath
			}

			v, err := importURL(importpath)
			if err != nil {
				if moduleErr != nil {
					return nil, fmt.Errorf("failed to import %s - %s, and %s", importpath, moduleErr.Error(), err.Error())
				}
				return nil, err
			}

			return v, nil
		})
	})
	return importExternalContentVar
}

func importModuleFile(importpath string) (rel.Value, error) {
	var mod module.Module = module.NewGoModule()

	m, err := mod.Get(importpath)
	if err != nil {
		return nil, err
	}

	return fileValue(filepath.Join(m.Dir, strings.TrimPrefix(importpath, m.Name)))
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

func importURL(url string) (rel.Value, error) {
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
		val := cache.getOrAdd(url, func() rel.Value { return bytesValue(NoPath, data) })
		return val, nil
	}
	return nil, fmt.Errorf("request %s failed: %s", url, resp.Status)
}

func fileValue(filename string) (rel.Value, error) {
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
	return bytesValue(filename, bytes), nil
}

func bytesValue(filename string, data []byte) rel.Value {
	eval := func() rel.Value {
		expr, err := Compile(filename, string(data))
		if err != nil {
			panic(err)
		}
		value, err := expr.Eval(rel.EmptyScope)
		if err != nil {
			panic(err)
		}
		return value
	}
	if filename != NoPath {
		return cache.getOrAdd(filename, eval)
	}

	return eval()
}
