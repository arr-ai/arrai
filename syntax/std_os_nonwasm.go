// +build !wasm

package syntax

import (
	"context"
	"fmt"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/arr-ai/arrai/tools"

	"github.com/arr-ai/arrai/rel"
)

func stdOsGetArgs() rel.Value {
	return strArrToRelArr(tools.Arguments)
}

func stdOsGetEnv(_ context.Context, value rel.Value) (rel.Value, error) {
	return rel.NewString([]rune(os.Getenv(value.(rel.String).String()))), nil
}

func stdOsPathSeparator() rel.Value {
	return rel.NewString([]rune{os.PathSeparator})
}

func stdOsPathListSeparator() rel.Value {
	return rel.NewString([]rune{os.PathListSeparator})
}

func stdOsCwd() rel.Value {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return rel.NewString([]rune(wd))
}

func stdOsExists(_ context.Context, v rel.Value) (rel.Value, error) {
	_, err := os.Stat(v.(rel.String).String())
	if os.IsNotExist(err) {
		return rel.NewBool(false), nil
	}
	if err != nil {
		return nil, err
	}
	return rel.NewBool(true), nil
}

func stdOsFile(_ context.Context, v rel.Value) (rel.Value, error) {
	f, err := ioutil.ReadFile(v.(rel.String).String())
	if err != nil {
		return nil, err
	}
	return rel.NewBytes(f), nil
}

// stdOsTree returns an array of paths to files within the given directory path.
func stdOsTree(v rel.Value) (rel.Value, error) {
	d, ok := v.(rel.String)
	if !ok {
		return nil, errors.Errorf("tree arg must be a string, not %T", v)
	}

	type entry map[string]interface{}
	root := entry{d.String(): make(entry)}
	err := filepath.Walk(d.String(), func(path string, info os.FileInfo, err error) error {
		ps := strings.Split(path, string(filepath.Separator))
		parent := root
		for i, p := range ps {
			if !info.IsDir() && i == len(ps)-1 {
				parent[p] = true
				return nil
			}

			v, ok := parent[p]
			if !ok {
				e := entry{}
				parent[p] = e
				parent = e
			} else {
				parent, ok = v.(entry)
				if !ok {
					return errors.Errorf("entry is not an entry")
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	//rel.NewDict(false, root)
	return rel.NewValue(root)
}

func stdOsIsATty(_ context.Context, value rel.Value) (rel.Value, error) {
	n, ok := value.(rel.Number)
	if !ok {
		return nil, fmt.Errorf("isatty arg must be a number, not %T", value)
	}
	fd, ok := n.Int()
	if !ok {
		return nil, fmt.Errorf("isatty arg must be an integer, not %s", value)
	}

	switch fd {
	case 0:
		return rel.NewBool(isatty.IsTerminal(os.Stdin.Fd())), nil
	case 1:
		return rel.NewBool(isatty.IsTerminal(os.Stdout.Fd())), nil
	}
	return nil, fmt.Errorf("isatty not implemented for %v", fd)
}

var stdOsStdinVar = newStdOsStdin(os.Stdin)
