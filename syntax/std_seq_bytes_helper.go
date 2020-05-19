package syntax

import (
	"strings"

	"github.com/arr-ai/arrai/rel"
)

// Joins byte array joiner to subject.
func bytesJoin(joiner, subject rel.Bytes) rel.Value {
	result := make([]byte, 0, subject.Count())
	for index, e := range subject.Bytes() {
		if index > 0 && index < subject.Count() {
			result = append(result, joiner.Bytes()...)
		}
		result = append(result, e)
	}

	return rel.NewBytes(result)
}

// Splits byte array subject by delimiter.
func bytesSplit(delimiter rel.Value, subject rel.Bytes) rel.Value {
	var splitted []string

	switch delimiter := delimiter.(type) {
	case rel.Bytes:
		splitted = strings.Split(subject.String(), delimiter.String())
	case rel.GenericSet:
		splitted = strings.Split(subject.String(), mustAsString(delimiter))
	}

	result := make([]rel.Value, 0, len(splitted))
	for _, s := range splitted {
		result = append(result, rel.NewBytes([]byte(s)).(rel.Value))
	}
	return rel.NewArray(result...)
}
