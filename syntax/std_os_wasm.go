// +build wasm

package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

func stdOsGetArgs() rel.Value {
	panic("not implemented")
}

func stdOsGetEnv(value rel.Value) (rel.Value, error) {
	panic("not implemented")
}

func stdOsIsATty(value rel.Value) (rel.Value, error) {
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

func stdOsFile(rel.Value) (rel.Value, error) {
	panic("not implemented")
}

var stdOsStdinVar = newStdOsStdin(nil)
