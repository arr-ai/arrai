// +build !wasm

package syntax

import (
	"io/ioutil"
	"os"

	"github.com/arr-ai/arrai/rel"
)

func getArgs() rel.Value {
	return strArrToRelArr(os.Args[2:])
}

func getEnv(value rel.Value) rel.Value {
	return rel.NewString([]rune(os.Getenv(value.(rel.String).String())))
}

func pathSeparator() rel.Value {
	return rel.NewString([]rune{os.PathSeparator})
}

func pathListSeparator() rel.Value {
	return rel.NewString([]rune{os.PathListSeparator})
}

func cwd() rel.Value {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return rel.NewString([]rune(wd))
}

func file(v rel.Value) rel.Value {
	f, err := ioutil.ReadFile(v.(rel.String).String())
	if err != nil {
		panic(err)
	}
	return rel.NewBytes(f)
}
