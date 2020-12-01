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
			return stdXMLDecode(v, translate.XMLDecodeConfig{StripFormatting: false})
		}),
		rel.NewNativeFunctionAttr("decoder", func(_ context.Context, config rel.Value) (rel.Value, error) {
			xmlConfig := translate.XMLDecodeConfig{StripFormatting: false}

			configTuple, ok := config.(rel.Tuple)
			if !ok {
				return nil, errors.Errorf("first arg to xml.decoder must be tuple, not %s", rel.ValueTypeAsString(config))
			}

			xmlConfig.StripFormatting = getConfigBool(configTuple, "strip_formatting")

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
	var bytes []byte
	switch v := v.(type) {
	case rel.String:
		bytes = []byte(v.String())
	case rel.Bytes:
		bytes = v.Bytes()
	}
	return translate.BytesXMLToArrai(bytes, config)
}
