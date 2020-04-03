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
			var bytes []byte
			switch v := v.(type) {
			case rel.String:
				bytes = []byte(v.String())
			case rel.Bytes:
				bytes = v.Bytes()
			}
			var data interface{}
			var err error
			if err = json.Unmarshal(bytes, &data); err == nil {
				var d rel.Value
				if d, err = translate.JSONToArrai(data); err == nil {
					return d
				}
			}
			panic(err)
		}),
	)
}
