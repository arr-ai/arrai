package syntax

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
)

// Joins byte array joiner to subject.
func bytesJoin(joiner rel.Value, subject rel.Bytes) rel.Value {
	var j []byte
	switch v := joiner.(type) {
	case rel.Bytes:
		j = v.Bytes()
	case rel.EmptySet:
		j = []byte{}
	default:
		panic(fmt.Errorf("bytesJoin: unsupported joiner type %T", joiner))
	}
	result := make([]byte, 0, subject.Count())
	for index, e := range subject.Bytes() {
		if index > 0 && index < subject.Count() {
			result = append(result, j...)
		}
		result = append(result, e)
	}

	return rel.NewBytes(result)
}

// Splits byte array subject by delimiter.
func bytesSplit(delimiter rel.Value, subject rel.Bytes) (rel.Value, error) {
	var splitted []string

	switch delimiter := delimiter.(type) {
	case rel.Bytes:
		splitted = strings.Split(subject.String(), delimiter.String())
	case rel.GenericSet, rel.EmptySet:
		delimStr, is := tools.ValueAsString(delimiter)
		if !is {
			return nil, errors.Errorf("//seq.split: delim not a string: %v", delimiter)
		}
		splitted = strings.Split(subject.String(), delimStr)
	default:
		return nil, errors.Errorf("//seq.split: delimiter and subject different types: "+
			"delimiter: %s, subject: %s", rel.ValueTypeAsString(delimiter), rel.ValueTypeAsString(subject))
	}

	result := make([]rel.Value, 0, len(splitted))
	for _, s := range splitted {
		result = append(result, rel.NewBytes([]byte(s)).(rel.Value))
	}
	return rel.NewArray(result...), nil
}
