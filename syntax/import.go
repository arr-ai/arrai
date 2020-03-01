package syntax

import (
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/arr-ai/arrai/rel"
)

var importLocalFileOnce sync.Once
var importLocalFileVar rel.Value

func importLocalFile() rel.Value {
	importLocalFileOnce.Do(func() {
		importLocalFileVar = rel.NewNativeFunction("//./", func(v rel.Value) rel.Value {
			filepath := v.String()
			if path.Ext(filepath) == "" {
				filepath += ".arrai"
			}
			data, err := ioutil.ReadFile(filepath)
			if err != nil {
				panic(err)
			}
			expr, err := Compile(filepath, string(data))
			if err != nil {
				panic(err)
			}
			value, err := expr.Eval(rel.EmptyScope)
			if err != nil {
				panic(err)
			}
			return value
		})
	})
	return importLocalFileVar
}

var importURL = rel.NewNativeFunction("//", func(v rel.Value) rel.Value {
	url := v.String()
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
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
	return rel.NewString([]rune(string(data)))
})
