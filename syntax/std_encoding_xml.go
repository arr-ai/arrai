package syntax

import (
	"context"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate"
	"github.com/go-errors/errors"
)

func stdEncodingXML() rel.Attr {
	return rel.NewTupleAttr(
		"xml",
		rel.NewNativeFunctionAttr("decode", func(_ context.Context, v rel.Value) (rel.Value, error) {
			return stdXMLDecode(v, translate.XMLDecodeConfig{TrimSurroundingWhitespace: false})
		}),
		rel.NewNativeFunctionAttr("decoder", func(_ context.Context, config rel.Value) (rel.Value, error) {
			xmlConfig := translate.XMLDecodeConfig{TrimSurroundingWhitespace: false}

			configTuple, ok := config.(rel.Tuple)
			if !ok {
				return nil, errors.Errorf("first arg to xml.decoder must be tuple, not %s", rel.ValueTypeAsString(config))
			}

			xmlConfig.TrimSurroundingWhitespace = getConfigBool(configTuple, "trimSurroundingWhitespace")

			return rel.NewNativeFunction("decode", func(_ context.Context, v rel.Value) (rel.Value, error) {
				return stdXMLDecode(v, xmlConfig)
			}), nil
		}),
		rel.NewNativeFunctionAttr("encode", func(_ context.Context, v rel.Value) (rel.Value, error) {
			return translate.BytesXMLFromArrai(v)
		}),
	)
}

func stdXMLDecode(v rel.Value, config translate.XMLDecodeConfig) (rel.Value, error) {
	bytes, ok := bytesOrStringAsUTF8(v)
	if !ok {
		return nil, errors.New("unhandled type for xml decoding")
	}
	return translate.BytesXMLToArrai(bytes, config)
}
