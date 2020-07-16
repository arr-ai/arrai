package syntax

import (
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate"
)

func stdEncodingYAML() rel.Attr {
	return rel.NewTupleAttr(
		"yaml",
		rel.NewNativeFunctionAttr("decode", func(v rel.Value) (rel.Value, error) {
			var bytes []byte
			switch v := v.(type) {
			case rel.String:
				bytes = []byte(v.String())
			case rel.Bytes:
				bytes = v.Bytes()
			}
			return bytesYAMLToArrai(bytes)
		}),
	)
}

func bytesYAMLToArrai(bytes []byte) (rel.Value, error) {
	val, err := translate.BytesYamlToArrai(bytes)
	if err != nil {
		return nil, err
	}
	return val, nil
}
