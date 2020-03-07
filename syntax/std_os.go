package syntax

import (
	"io/ioutil"

	"github.com/arr-ai/arrai/rel"
)

func stdOs() rel.Attr {
	return rel.NewAttr("os", rel.NewTuple(
		rel.NewAttr("args", getArgs()),
		rel.NewAttr("path_separator", pathSeparator()),
		rel.NewAttr("path_list_separator", pathListSeparator()),
		rel.NewAttr("cwd", cwd()),
		rel.NewNativeFunctionAttr("file", func(value rel.Value) rel.Value {
			f, err := ioutil.ReadFile(value.(rel.String).String())
			if err != nil {
				panic(err)
			}
			return rel.NewBytes(f)
		}),
		rel.NewNativeFunctionAttr("get_env", getEnv),
	))
}
