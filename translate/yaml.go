package translate

import (
	"gopkg.in/yaml.v2"

	"github.com/arr-ai/arrai/rel"
)

func (t Translator) BytesYamlToArrai(bytes []byte) (rel.Value, error) {
	var m interface{}
	var err error
	if err = yaml.Unmarshal(bytes, &m); err == nil {
		var d rel.Value
		if d, err = t.ToArrai(m); err == nil {
			return d, nil
		}
	}
	return nil, err
}
