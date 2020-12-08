package syntax

import (
	"context"
	"errors"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate"
)

func stdEncodingYAML() rel.Attr {
	return rel.NewTupleAttr(
		"yaml",
		rel.NewNativeFunctionAttr("decode", func(_ context.Context, v rel.Value) (rel.Value, error) {
			bytes, ok := bytesOrStringAsUTF8(v)
			if !ok {
				return nil, errors.New("unhandled type for yaml decoding")
			}
			return bytesYAMLToArrai(bytes)
		}),
	)
}

func bytesYAMLToArrai(bytes []byte) (rel.Value, error) {
	return translate.BytesYamlToArrai(bytes)
}
