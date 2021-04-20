package syntax

import (
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate"
	"github.com/go-errors/errors"
)

func decodeXML(v rel.Value, config translate.XMLDecodeConfig) (rel.Value, error) {
	bytes, ok := bytesOrStringAsUTF8(v)
	if !ok {
		return nil, errors.New("unhandled type for xml decoding")
	}
	if config == (translate.XMLDecodeConfig{}) {
		config.TrimSurroundingWhitespace = false
	}
	return translate.BytesXMLToArrai(bytes, config)
}

// parseXMLConfig returns the config arg as a xmlDecodeConfig.
func parseXMLConfig(configArg rel.Value) (*translate.XMLDecodeConfig, error) {
	config, ok := configArg.(*rel.GenericTuple)
	if !ok {
		return nil, errors.Errorf("first arg (config) must be tuple, not %s", rel.ValueTypeAsString(configArg))
	}
	whitespace, ok := config.Get("TrimSurroundingWhitespace")
	if !ok {
		return &translate.XMLDecodeConfig{}, nil
	}

	return &translate.XMLDecodeConfig{TrimSurroundingWhitespace: whitespace.IsTrue()}, nil
}
