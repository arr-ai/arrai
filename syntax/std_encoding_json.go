package syntax

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate"
)

func stdEncodingJSON() rel.Attr {
	return rel.NewTupleAttr(
		"json",
		rel.NewNativeFunctionAttr("decode", func(_ context.Context, v rel.Value) (rel.Value, error) {
			bytes, ok := stdEncodingBytesOrStringAsUTF8(v)
			if !ok {
				return nil, errors.New("unhandled type for json decoding")
			}
			return bytesJSONToArrai(bytes)
		}),
		rel.NewNativeFunctionAttr("encode", func(_ context.Context, v rel.Value) (rel.Value, error) {
			data, err := translate.FromArrai(v)
			if err != nil {
				return nil, err
			}
			result := bytes.NewBuffer([]byte{})
			enc := json.NewEncoder(result)
			enc.SetEscapeHTML(false)
			err = enc.Encode(data)
			if err != nil {
				return nil, err
			}
			return rel.NewBytes(result.Bytes()), nil
		}),
		rel.NewNativeFunctionAttr("encode_indent", func(_ context.Context, v rel.Value) (rel.Value, error) {
			data, err := translate.FromArrai(v)
			if err != nil {
				return nil, err
			}
			result := bytes.NewBuffer([]byte{})
			enc := json.NewEncoder(result)
			enc.SetEscapeHTML(false)
			enc.SetIndent("", "  ")
			err = enc.Encode(data)
			if err != nil {
				return nil, err
			}
			return rel.NewBytes(result.Bytes()), nil
		}),
	)
}

func bytesJSONToArrai(bytes []byte) (rel.Value, error) {
	var data interface{}
	var err error
	if err = json.Unmarshal(bytes, &data); err == nil {
		var d rel.Value
		if d, err = translate.ToArrai(data); err == nil {
			return d, nil
		}
	}
	return nil, err
}
