// +build !wasm

package syntax

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/afero"

	"github.com/arr-ai/arrai/pkg/arraictx"

	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/rel"
)

// stdOsGetArgs returns a rel.Array of the program arguments in the context.
func stdOsGetArgs(ctx context.Context, _ rel.Value) (rel.Value, error) {
	if arraictx.IsCompiling(ctx) {
		return rel.None, nil
	}
	return strArrToRelArr(arraictx.Args(ctx)), nil
}

func stdOsGetEnv(ctx context.Context, value rel.Value) (rel.Value, error) {
	if arraictx.IsCompiling(ctx) {
		return rel.None, nil
	}
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

func stdOsExists(ctx context.Context, v rel.Value) (rel.Value, error) {
	if arraictx.IsCompiling(ctx) {
		return rel.None, nil
	}
	filePath, err := toString("os.exists", v)
	if err != nil {
		return nil, err
	}
	_, err = ctxfs.RuntimeFsFrom(ctx).Stat(filePath.String())
	switch {
	case os.IsNotExist(err):
		return rel.NewBool(false), nil
	case err != nil:
		return nil, err
	}
	return rel.NewBool(true), nil
}

func stdOsFile(ctx context.Context, v rel.Value) (rel.Value, error) {
	if arraictx.IsCompiling(ctx) {
		return nil, os.ErrNotExist
	}
	filePath, err := toString("os.file", v)
	if err != nil {
		return nil, err
	}
	buf, err := afero.ReadFile(ctxfs.RuntimeFsFrom(ctx), filePath.String())
	if err != nil {
		//TODO: wrap this in an arrai error message
		return nil, err
	}

	return rel.NewBytes(buf), nil
}

// stdOsTree returns an array of paths to files within the given directory path.
func stdOsTree(ctx context.Context, v rel.Value) (rel.Value, error) {
	if arraictx.IsCompiling(ctx) {
		return rel.None, nil
	}
	d, ok := v.(rel.String)
	if !ok {
		return nil, errors.Errorf("tree arg must be a string, not %T", v)
	}

	b := rel.NewSetBuilder()
	err := filepath.Walk(d.String(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Fix for consistent behavior on Windows.
		if path == d.String() && runtime.GOOS == windows {
			path = strings.ReplaceAll(path, `/`, string(filepath.Separator))
		}
		b.Add(rel.NewTuple(
			rel.NewStringAttr("name", []rune(info.Name())),
			rel.NewStringAttr("path", []rune(path)),
			rel.NewBoolAttr("is_dir", info.IsDir()),
			rel.NewFloatAttr("size", float64(info.Size())),
			rel.NewFloatAttr("mod_time", float64(info.ModTime().UnixNano())),
		))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return b.Finish()
}

func stdOsIsATty(ctx context.Context, value rel.Value) (rel.Value, error) {
	if arraictx.IsCompiling(ctx) {
		return rel.None, nil
	}
	n, ok := value.(rel.Number)
	if !ok {
		return nil, fmt.Errorf("isatty arg must be a number, not %s", rel.ValueTypeAsString(value))
	}
	fd, ok := n.Int()
	if !ok {
		return nil, fmt.Errorf("isatty arg must be an integer, not %v", value)
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

func toString(fnName string, v rel.Value) (rel.String, error) {
	str, isString := v.(rel.String)
	if !isString {
		return rel.String{}, fmt.Errorf("//%s: argument does not resolve to String: %v", fnName, v)
	}
	return str, nil
}
