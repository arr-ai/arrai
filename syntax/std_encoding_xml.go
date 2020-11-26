package syntax

import (
	"context"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate"
)

func stdEncodingXML() rel.Attr {
	return rel.NewTupleAttr(
		"xml",
		rel.NewNativeFunctionAttr("decode", func(_ context.Context, v rel.Value) (rel.Value, error) {
			return stdXMLDecode(v, translate.DecodeConfig{StripFormatting: false})
		}),
		rel.NewNativeFunctionAttr("decoder", func(_ context.Context, config rel.Value) (rel.Value, error) {
			strip := config.IsTrue()
			return rel.NewNativeFunction("decode", func(_ context.Context, v rel.Value) (rel.Value, error) {
				return stdXMLDecode(v, translate.DecodeConfig{StripFormatting: strip})
			}), nil
		}),
		rel.NewNativeFunctionAttr("encode", func(_ context.Context, v rel.Value) (rel.Value, error) {
			return translate.BytesXMLFromArrai(v)
		}),
	)
}

func stdXMLDecode(v rel.Value, config translate.DecodeConfig) (rel.Value, error) {
	var bytes []byte
	switch v := v.(type) {
	case rel.String:
		bytes = []byte(v.String())
	case rel.Bytes:
		bytes = v.Bytes()
	}
	return translate.BytesXMLToArrai(bytes, config)
}
