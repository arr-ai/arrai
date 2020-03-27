package syntax

import "github.com/arr-ai/arrai/rel"

func stdOs() rel.Attr {
	return rel.NewTupleAttr("os",
		rel.NewAttr("args", getArgs()),
		rel.NewAttr("path_separator", pathSeparator()),
		rel.NewAttr("path_list_separator", pathListSeparator()),
		rel.NewAttr("cwd", cwd()),
		rel.NewNativeFunctionAttr("file", file),
		rel.NewNativeFunctionAttr("get_env", getEnv),
	)
}
