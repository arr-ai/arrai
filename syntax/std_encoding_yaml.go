package syntax

import (
	"context"

	"github.com/go-errors/errors"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate"
)

func newYAMLDecodeConfig() yamlDecodeConfig {
	return yamlDecodeConfig{strict: true}
}

type yamlDecodeConfig struct {
	strict bool
}

func stdEncodingYAML() rel.Attr {
	return rel.NewTupleAttr(
		"yaml",
		rel.NewNativeFunctionAttr(decodeAttr, func(_ context.Context, value rel.Value) (rel.Value, error) {
			return yamlDecodeFnBody("yaml.decode", value, newYAMLDecodeConfig())
		}),

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
