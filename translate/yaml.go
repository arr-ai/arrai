package translate

import (
	"github.com/arr-ai/arrai/rel"
	"gopkg.in/yaml.v2"
)

func BytesYamlToArrai(bytes []byte) (rel.Value, error) {
	var m interface{}
	var err error
	if err = yaml.Unmarshal(bytes, &m); err == nil {
		var d rel.Value
		if d, err = ToArrai(m); err == nil {
			return d, nil
		}
	}
	return nil, err
}
