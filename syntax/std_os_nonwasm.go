// +build !wasm

package syntax

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"

	"github.com/mattn/go-isatty"

	"github.com/arr-ai/arrai/pkg/ctxfs"
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

func stdOsExists(ctx context.Context, v rel.Value) (rel.Value, error) {
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
	filePath, err := toString("os.file", v)
	if err != nil {
		return nil, err
	}
	buf, err := ctxfs.ReadFile(ctxfs.RuntimeFsFrom(ctx), filePath.String())
	if err != nil {
		//TODO: wrap this in an arrai error message
		return nil, err
	}

	return rel.NewBytes(buf), nil
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

func stdOsExec(_ context.Context, value rel.Value) (rel.Value, error) {
	var cmd *exec.Cmd
	switch t := value.(type) {
	case rel.Array:
		if len(t.Values()) == 0 {
			return nil, errors.Errorf("//os.exec arg must not be empty")
		}

		name := t.Values()[0].String()
		args := make([]string, len(t.Values()))
		for i, v := range t.Values() {
			if i == 0 {
				continue
			}
			args = append(args, v.String())
		}
		cmd = exec.Command(name, args...)
	case rel.String:
		cmd = exec.Command(t.String())
	default:
		return nil, errors.Errorf("//os.exec arg must be a string or array, not %T", value)
	}
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return rel.NewBytes(out), nil
}

var stdOsStdinVar = newStdOsStdin(os.Stdin)

func toString(fnName string, v rel.Value) (rel.String, error) {
	str, isString := v.(rel.String)
	if !isString {
		return rel.String{}, fmt.Errorf("//%s: argument does not resolve to String: %v", fnName, v)
	}
	return str, nil
}
