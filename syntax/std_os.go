package syntax

import "github.com/arr-ai/arrai/rel"

func stdOs() rel.Attr {
	return rel.NewTupleAttr("os",
		rel.NewAttr("args", stdOsGetArgs()),
		rel.NewAttr("path_separator", stdOsPathSeparator()),
		rel.NewAttr("path_list_separator", stdOsPathListSeparator()),
		rel.NewAttr("cwd", stdOsCwd()),
		rel.NewNativeFunctionAttr("file", stdOsFile),
		rel.NewNativeFunctionAttr("get_env", stdOsGetEnv),
		rel.NewNativeFunctionAttr("&stdin", stdOsStdin),
	)
}
