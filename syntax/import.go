package syntax

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools/module"
)

var importLocalFileOnce sync.Once
var importLocalFileVar rel.Value

func importLocalFile() rel.Value {
	importLocalFileOnce.Do(func() {
		importLocalFileVar = rel.NewNativeFunction("//./", func(v rel.Value) rel.Value {
			s, ok := rel.AsString(v.(rel.Set))
			if !ok {
				panic(fmt.Errorf("cannot convert %#v to string", v))
			}
			v, err := fileValue(s.String())
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

		return bytesValue("", data), nil
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
