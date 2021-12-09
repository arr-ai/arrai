package syntax

import (
	"bytes"
	"context"

	"gopkg.in/yaml.v3"

	"github.com/go-errors/errors"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate"
)

func newYAMLDecodeConfig() yamlDecodeConfig {
	return yamlDecodeConfig{strict: true}
}

func newYAMLEncodeConfig() yamlEncodeConfig {
	return yamlEncodeConfig{strict: true, indent: 4}
}

type yamlDecodeConfig struct {
	strict bool
}

type yamlEncodeConfig struct {
	strict bool
	indent int
}

func stdEncodingYAML() rel.Attr {
	return rel.NewTupleAttr(
		"yaml",
		rel.NewNativeFunctionAttr(decodeAttr, func(_ context.Context, value rel.Value) (rel.Value, error) {
			return yamlDecodeFnBody("yaml.decode", value, newYAMLDecodeConfig())
		}),

		//nolint:dupl // Not a duplicate of YAML encoder
		rel.NewNativeFunctionAttr(decoderAttr, func(_ context.Context, configValue rel.Value) (rel.Value, error) {
			fn := "yaml.decoder"
			config := newYAMLDecodeConfig()

			configTuple, ok := configValue.(rel.Tuple)
			if !ok {
				return nil, errors.Errorf("first arg to %s must be tuple, not %s", fn, rel.ValueTypeAsString(configValue))
			}

			if strict, ok := getConfigBool(configTuple, "strict"); ok {
				config.strict = strict
			}

			return rel.NewNativeFunction("decode", func(_ context.Context, value rel.Value) (rel.Value, error) {
				return yamlDecodeFnBody("yaml.decoder payload", value, config)
			}), nil
		}),

		rel.NewNativeFunctionAttr("encode", func(_ context.Context, value rel.Value) (rel.Value, error) {
			return yamlEncodeFnBody(value, newYAMLEncodeConfig())
		}),

		rel.NewNativeFunctionAttr("encoder", func(_ context.Context, configValue rel.Value) (rel.Value, error) {
			fn := "yaml.encoder"
			config := newYAMLEncodeConfig()

			configTuple, ok := configValue.(rel.Tuple)
			if !ok {
				return nil, errors.Errorf("first arg to %s must be tuple, not %s", fn, rel.ValueTypeAsString(configValue))
			}

			indent, err := getConfigInt(configTuple, fn, "indent", config.indent)
			if err != nil {
				return nil, err
			}
			config.indent = indent

			if strict, ok := getConfigBool(configTuple, "strict"); ok {
				config.strict = strict
			}

			return rel.NewNativeFunction("encode", func(_ context.Context, value rel.Value) (rel.Value, error) {
				return yamlEncodeFnBody(value, config)
			}), nil
		}),
	)
}

func yamlDecodeFnBody(fn string, value rel.Value, config yamlDecodeConfig) (rel.Value, error) {
	bs, ok := bytesOrStringAsUTF8(value)
	if !ok {
		return nil, errors.Errorf("first arg to %s must be string or bytes, not %s", fn, rel.ValueTypeAsString(value))
	}
	val, err := translate.NewTranslator(config.strict).BytesYamlToArrai(bs)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func yamlEncodeFnBody(value rel.Value, config yamlEncodeConfig) (rel.Value, error) {
	t := translate.NewTranslator(config.strict)
	data, err := t.FromArrai(value)
	if err != nil {
		return nil, err
	}
	result := bytes.NewBuffer([]byte{})
	enc := yaml.NewEncoder(result)
	enc.SetIndent(config.indent)
	err = enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return rel.NewBytes(result.Bytes()), nil
}
