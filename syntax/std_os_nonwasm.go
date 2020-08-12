// +build !wasm

package syntax

import (
	"io/ioutil"
	"os"

	"github.com/arr-ai/arrai/tools"

	"github.com/arr-ai/arrai/rel"
)

func stdOsGetArgs() rel.Value {
	return strArrToRelArr(tools.Arguments)
}

func stdOsGetEnv(value rel.Value) (rel.Value, error) {
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

func stdOsFile(v rel.Value) (rel.Value, error) {
	f, err := ioutil.ReadFile(v.(rel.String).String())
	if err != nil {
		return nil, err
	}
	return rel.NewBytes(f), nil
}

// stdinHasInput returns true if there is data to read on stdin.
func stdinHasInput() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

var stdOsStdinVar = newStdOsStdin(os.Stdin, stdinHasInput())
