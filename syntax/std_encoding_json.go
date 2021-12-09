package syntax

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/go-errors/errors"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate"
)

func newJSONDecodeConfig() jsonDecodeConfig {
	return jsonDecodeConfig{strict: true}
}

func newJSONEncodeConfig() jsonEncodeConfig {
	return jsonEncodeConfig{prefix: "", indent: "", strict: true, escapeHTML: false}
}

type jsonDecodeConfig struct {
	strict bool
}

type jsonEncodeConfig struct {
	prefix     string
	indent     string
	escapeHTML bool
	strict     bool
}

func stdEncodingJSON() rel.Attr {
	return rel.NewTupleAttr(
		"json",
		rel.NewNativeFunctionAttr(decodeAttr, func(_ context.Context, value rel.Value) (rel.Value, error) {
			return jsonDecodeFnBody("json.decode", value, newJSONDecodeConfig())
		}),

		//nolint:dupl // Not a duplicate of JSON encoder
		rel.NewNativeFunctionAttr(decoderAttr, func(_ context.Context, configValue rel.Value) (rel.Value, error) {
			fn := "json.decoder"
			config := newJSONDecodeConfig()

			configTuple, ok := configValue.(rel.Tuple)
			if !ok {
				return nil, errors.Errorf("first arg to %s must be tuple, not %s", fn, rel.ValueTypeAsString(configValue))
			}

			if strict, ok := getConfigBool(configTuple, "strict"); ok {
				config.strict = strict
			}

			return rel.NewNativeFunction("decode", func(_ context.Context, value rel.Value) (rel.Value, error) {
				return jsonDecodeFnBody("json.decoder payload", value, config)
			}), nil
		}),

		rel.NewNativeFunctionAttr("encode", func(_ context.Context, value rel.Value) (rel.Value, error) {
			return jsonEncodeFnBody(value, newJSONEncodeConfig())
		}),

		rel.NewNativeFunctionAttr("encode_indent", func(_ context.Context, value rel.Value) (rel.Value, error) {
			config := newJSONEncodeConfig()
			config.indent = "  "
			return jsonEncodeFnBody(value, config)
		}),

		rel.NewNativeFunctionAttr("encoder", func(_ context.Context, configValue rel.Value) (rel.Value, error) {
			fn := "json.encoder"
			config := newJSONEncodeConfig()

			configTuple, ok := configValue.(rel.Tuple)
			if !ok {
				return nil, errors.Errorf("first arg to %s must be tuple, not %s", fn, rel.ValueTypeAsString(configValue))
			}

			prefix, err := getConfigString(configTuple, fn, "prefix", config.prefix)
			if err != nil {
				return nil, err
			}
			config.prefix = prefix

			indent, err := getConfigString(configTuple, fn, "indent", config.indent)
			if err != nil {
				return nil, err
			}
			config.indent = indent

			config.escapeHTML, _ = getConfigBool(configTuple, "escapeHTML")
			if strict, ok := getConfigBool(configTuple, "strict"); ok {
				config.strict = strict
			}

			return rel.NewNativeFunction("encode", func(_ context.Context, value rel.Value) (rel.Value, error) {
				return jsonEncodeFnBody(value, config)
			}), nil
		}),
	)
}

func jsonDecodeFnBody(fn string, value rel.Value, config jsonDecodeConfig) (rel.Value, error) {
	bs, ok := bytesOrStringAsUTF8(value)
	if !ok {
		return nil, errors.Errorf("first arg to %s must be string or bytes, not %s", fn, rel.ValueTypeAsString(value))
	}
	var data interface{}
	var err error
	if err = json.Unmarshal(bs, &data); err == nil {
		var d rel.Value
		t := translate.NewTranslator(config.strict)
		if d, err = t.ToArrai(data); err == nil {
			return d, nil
		}
	}
	return nil, err
}

func jsonEncodeFnBody(value rel.Value, config jsonEncodeConfig) (rel.Value, error) {
	t := translate.NewTranslator(config.strict)
	data, err := t.FromArrai(value)
	if err != nil {
		return nil, err
	}
	result := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(result)
	enc.SetEscapeHTML(config.escapeHTML)
	enc.SetIndent(config.prefix, config.indent)
	err = enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return rel.NewBytes(result.Bytes()), nil
}
