// +build !wasm

package syntax

import (
	"io/ioutil"
	"os"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
)

func stdOsGetArgs() rel.Value {
	args, found := tools.Arguments.Get("Arguments")
	if found {
		return strArrToRelArr(args.([]string))
	}
	return rel.NewArray()
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

var stdOsStdinVar = newStdOsStdin(os.Stdin)
