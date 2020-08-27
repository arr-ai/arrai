// +build !wasm

package syntax

import (
	"context"
	"fmt"
	"github.com/arr-ai/arrai/tools"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"

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
func stdOsTree(_ context.Context, v rel.Value) (rel.Value, error) {
	d, ok := v.(rel.String)
	if !ok {
		return nil, errors.Errorf("tree arg must be a string, not %T", v)
	}

	fs := []rel.Value{}
	err := filepath.Walk(d.String(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fs = append(fs, rel.NewTuple(
			rel.NewStringAttr("name", []rune(info.Name())),
			rel.NewStringAttr("path", []rune(path)),
			rel.NewBoolAttr("isDir", info.IsDir()),
			rel.NewFloatAttr("size", float64(info.Size())),
			rel.NewFloatAttr("modTime", float64(info.ModTime().UnixNano())),
		))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rel.NewSet(fs...), nil
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
