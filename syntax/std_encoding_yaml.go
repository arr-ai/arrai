package syntax

import (
	"context"
	"fmt"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate"
)

func stdEncodingYAML() rel.Attr {
	return rel.NewTupleAttr(
		"yaml",
		rel.NewNativeFunctionAttr(decodeAttr, func(_ context.Context, v rel.Value) (rel.Value, error) {
			bytes, ok := bytesOrStringAsUTF8(v)
			if !ok {
				return nil, fmt.Errorf("unexpected arrai object type: %s", rel.ValueTypeAsString(v))
			}
			return bytesYAMLToArrai(bytes)
		}),
		//TODO: configurable decoder
	)
}

func bytesYAMLToArrai(bytes []byte) (rel.Value, error) {
	val, err := translate.BytesYamlToArrai(bytes)
	if err != nil {
		return nil, err
	}
	return val, nil
}
