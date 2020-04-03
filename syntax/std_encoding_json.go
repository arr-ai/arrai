package syntax

import (
	"encoding/json"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate"
)

func stdEncodingJSON() rel.Attr {
	return rel.NewTupleAttr(
		"json",
		rel.NewNativeFunctionAttr("decode", func(v rel.Value) rel.Value {
			s := mustAsString(v)
			var data interface{}
			var err error
			if err = json.Unmarshal([]byte(s), &data); err == nil {
				var d rel.Value
				if d, err = translate.JSONToArrai(data); err == nil {
					return d
				}
			}
			panic(err)
		}),
	)
}
