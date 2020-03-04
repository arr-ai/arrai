package syntax

import (
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools/module"
)

// go run ./cmd/arrai e "//./examples/jsfuncs/jsfuncs"
var importLocalFileOnce sync.Once
var importLocalFileVar rel.Value

func importLocalFile() rel.Value {
	importLocalFileOnce.Do(func() {
		importLocalFileVar = rel.NewNativeFunction("//./", func(v rel.Value) rel.Value {
			filename := v.String()
			if path.Ext(filename) == "" {
				filename += ".arrai"
			}
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				panic(err)
			}
			return bytesValue(filename, data)
		})
	})
	return importLocalFileVar
}

// go run ./cmd/arrai e "//github.com/ChloePlanet/'arrai-examples'/add"
var importModuleFileOnce sync.Once
var importModuleFileVar rel.Value

func importModuleFile() rel.Value {
	importModuleFileOnce.Do(func() {
		importModuleFileVar = rel.NewNativeFunction("//", func(v rel.Value) rel.Value {
			importpath := v.String()

			var mod module.Module
			mod = module.NewGoModule()

			m, err := mod.Get(importpath)
			if err != nil {
				panic(err)
			}

			relname, err := filepath.Rel(m.Name, importpath)
			if err != nil {
				panic(err)
			}

			filename := filepath.Join(m.Dir, relname)
			if path.Ext(filename) == "" {
				filename += ".arrai"
			}

			data, err := ioutil.ReadFile(filename)
			if err != nil {
				panic(err)
			}
			return bytesValue(filename, data)
		})
	})
	return importModuleFileVar
}

// go run ./cmd/arrai e "//jsonplaceholder.typicode.com/todos/'1'/userId"
var importURLOnce sync.Once
var importURLVar rel.Value

func importURL() rel.Value {
	importURLOnce.Do(func() {
		importURLVar = rel.NewNativeFunction("//", func(v rel.Value) rel.Value {
			url := v.String()
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				// TBD: always https?
				url = "https://" + url
			}
			resp, err := http.Get(url) //nolint:gosec
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			return bytesValue("", data)
		})
	})
	return importModuleFileVar
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
