package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

func stdUnicodeUTF8() rel.Attr {
	return rel.NewTupleAttr(
		"utf8",
		rel.NewNativeFunctionAttr("encode", func(v rel.Value) rel.Value {
			return rel.NewBytes([]byte(mustAsString(v)))
		}),
	)
}

func stdUnicode() rel.Attr {
	return rel.NewTupleAttr("unicode",
		stdUnicodeUTF8(),
	)
}
