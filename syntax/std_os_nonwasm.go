// +build !wasm

package syntax

import (
	"io/ioutil"
	"os"

	"github.com/arr-ai/arrai/rel"
)

func stdOsGetArgs() rel.Value {
	var offset int
	switch os.Args[0] {
	case "ai", "ax":
		offset = 1
	default:
		offset = 2
	}
	return strArrToRelArr(os.Args[offset:])
}

func stdOsGetEnv(value rel.Value) rel.Value {
	return rel.NewString([]rune(os.Getenv(value.(rel.String).String())))
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

func stdOsFile(v rel.Value) rel.Value {
	f, err := ioutil.ReadFile(v.(rel.String).String())
	if err != nil {
		panic(err)
	}
	return rel.NewBytes(f)
}

var stdOsStdinVar = newStdOsStdin(os.Stdin)
