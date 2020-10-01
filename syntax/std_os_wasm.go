// +build wasm

package syntax

import (
	"context"

	"github.com/arr-ai/arrai/rel"
)

func stdOsGetArgs(ctx context.Context, _ rel.Value) (rel.Value, error) {
	panic("not implemented")
}

func stdOsGetEnv(context.Context, rel.Value) (rel.Value, error) {
	panic("not implemented")
}

func stdOsIsATty(context.Context, rel.Value) (rel.Value, error) {
	panic("not implemented")
}

func stdOsPathSeparator() rel.Value {
	panic("not implemented")
}

func stdOsPathListSeparator() rel.Value {
	panic("not implemented")
}

func stdOsCwd() rel.Value {
	panic("not implemented")
}

func stdOsExists(context.Context, rel.Value) (rel.Value, error) {
	panic("not implemented")
}

func stdOsFile(context.Context, rel.Value) (rel.Value, error) {
	panic("not implemented")
}

func stdOsExec(context.Context, rel.Value) (rel.Value, error) {
	panic("not implemented")
}

func stdOsTree(context.Context, rel.Value) (rel.Value, error) {
	panic("not implemented")
}

var stdOsStdinVar = newStdOsStdin(nil)
