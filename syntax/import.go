package syntax

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/arr-ai/arrai/rel"
	"github.com/spf13/afero"
)

const arraiRootMarker = "go.mod"

var importLocalFileOnce sync.Once
var importLocalFileVar rel.Value
var fs = afero.NewOsFs()

func importLocalFile(fromRoot bool) rel.Value {
	importLocalFileOnce.Do(func() {
		importLocalFileVar = rel.NewNativeFunction("//./", func(v rel.Value) rel.Value {
			filepath := v.String()
			if fromRoot {
				pwd, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				rootPath, err := findRootFromModule(pwd, fs)
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

func findRootFromModule(modulePath string, fs afero.Fs) (string, error) {
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
		exists, err := afero.Exists(fs, filepath.Join(currentPath, arraiRootMarker))
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
