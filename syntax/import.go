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

var importLocalFileOnce sync.Once
var importLocalFileVar rel.Value

// go run ./cmd/arrai e "//./examples/jsfuncs/jsfuncs"
func importLocalFile() rel.Value {
	importLocalFileOnce.Do(func() {
		importLocalFileVar = rel.NewNativeFunction("//./", func(v rel.Value) rel.Value {
			filename := v.String()
			if path.Ext(filename) == "" {
				filename += ".arrai"
			}
			return fileValue(filename)
		})
	})
	return importLocalFileVar
}

// go run ./cmd/arrai e "//github.com/'arr-ai'/arrai/examples/grpc/grpc"
var importModule = rel.NewNativeFunction("//", func(v rel.Value) rel.Value {
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
	// return fileValue(filepath.Join(m.Dir, relname))
	return rel.NewString([]rune(string(data)))
})

func fileValue(filename string) rel.Value {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
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

// go run ./cmd/arrai e "//jsonplaceholder.typicode.com/todos/'1'/userId"
// var importURLOnce sync.Once
var importURL = rel.NewNativeFunction("//", func(v rel.Value) rel.Value {
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
	// ...
	return rel.NewString([]rune(string(data)))
})
