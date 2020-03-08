// +build wasm

package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

func getArgs() rel.Value {
	return rel.NewSet()
}

func getEnv(value rel.Value) rel.Value {
	return rel.NewSet()
}

func pathSeparator() rel.Value {
	return rel.NewSet()
}

func pathListSeparator() rel.Value {
	return rel.NewSet()
}

func cwd() rel.Value {
	return rel.NewSet()
}

func file(rel.Value) rel.Value {
	return rel.NewBytes([]byte{})
}
