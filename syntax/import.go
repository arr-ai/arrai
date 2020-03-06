package syntax

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/arr-ai/arrai/rel"
)

var importLocalFileOnce sync.Once
var importLocalFileVar rel.Value

func importLocalFile(fromRoot bool) rel.Value {
	importLocalFileOnce.Do(func() {
		importLocalFileVar = rel.NewNativeFunction("//./", func(v rel.Value) rel.Value {

			filepath := v.String()
			if fromRoot {
				pwd, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				rootPath, err := findProjectRootDir(pwd)
				if err != nil {
					panic(err)
				}
				filepath = rootPath + "/" + filepath
			}
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

func findProjectRootDir(currentDir string) (string, error) {
	files, err := ioutil.ReadDir(currentDir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if f.Name() == "go.mod" {
			return currentDir, nil
		}
	}
	if currentDir == "" {
		return "", errors.New("cannot find root")
	}

	ss := strings.Split(currentDir, "/")
	ss[len(ss)-1] = ""
	currentDir = strings.Join(ss, "/")
	currentDir = strings.TrimRight(currentDir, "/")

	return findProjectRootDir(currentDir)

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
