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
	"github.com/arr-ai/arrai/tools/module"
)

const arraiRootMarker = "go.mod"

var importLocalFileOnce sync.Once
var importLocalFileVar rel.Value
var cache importCache = newCache()

func importLocalFile(fromRoot bool) rel.Value {
	importLocalFileOnce.Do(func() {
		importLocalFileVar = rel.NewNativeFunction("//./", func(v rel.Value) rel.Value {
			s, ok := rel.AsString(v.(rel.Set))
			if !ok {
				panic(fmt.Errorf("cannot convert %#v to string", v))
			}

			filename := s.String()
			if fromRoot {
				pwd, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				rootPath, err := findRootFromModule(pwd)
				if err != nil {
					panic(err)
				}
				filename = rootPath + "/" + filename
			}

			v, err := fileValue(filename)
			if err != nil {
				panic(err)
			}

			return v
		})
	})
	return importLocalFileVar
}

var importExternalContentOnce sync.Once
var importExternalContentVar rel.Value

func importExternalContent() rel.Value {
	importExternalContentOnce.Do(func() {
		importExternalContentVar = rel.NewNativeFunction("//", func(v rel.Value) rel.Value {
			s, ok := rel.AsString(v.(rel.Set))
			if !ok {
				panic(fmt.Errorf("cannot convert %#v to string", v))
			}
			importpath := s.String()

			var moduleErr error

			if !strings.HasPrefix(importpath, "http://") && !strings.HasPrefix(importpath, "https://") {
				v, err := importModuleFile(importpath)
				if err == nil {
					return v
				}
				moduleErr = err

				// Since an explicit schema is allowed, it's OK to assume https as the default.
				importpath = "https://" + importpath
			}

			v, err := importURL(importpath)
			if err != nil {
				if moduleErr != nil {
					panic(fmt.Errorf("failed to import %s - %s, and %s", importpath, moduleErr.Error(), err.Error()))
				}
				panic(err)
			}

			return v
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
		exists := fileExists(filepath.Join(currentPath, arraiRootMarker))
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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
		if exists, val := cache.exists(url); exists {
			return val, nil
		}
		// Pass space as file path as it is remote url and can't be found in local env
		val := bytesValue(NoPath, data)
		cache.add(url, val)
		return val, nil
	}
	return nil, fmt.Errorf("request %s failed: %s", url, resp.Status)
}

func fileValue(filename string) (rel.Value, error) {
	if filepath.Ext(filename) == "" {
		filename += ".arrai"
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return bytesValue(filename, data), nil
}

func bytesValue(filename string, data []byte) rel.Value {
	// maybe filename is "", so add check first
	if filename != NoPath {
		if exists, val := cache.exists(filename); exists {
			return val
		}
	}

	expr, err := Compile(filename, string(data))
	if err != nil {
		panic(err)
	}
	value, err := expr.Eval(rel.EmptyScope)
	if err != nil {
		panic(err)
	}

	// maybe filename is "", so add check first
	if filename != NoPath {
		cache.add(filename, value)
	}
	return value
}
