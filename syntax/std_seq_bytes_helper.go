package syntax

import (
	"strings"

	"github.com/arr-ai/arrai/rel"
)

// BytesJoin join b to a
func BytesJoin(a, b rel.Bytes) rel.Value {
	result := make([]byte, 0, a.Count())
	for index, e := range a.Bytes() {
		if index > 0 && index < a.Count() {
			result = append(result, b.Bytes()...)
		}
		result = append(result, e)
	}

	return rel.NewBytes(result)
}

// BytesSplit split a by b
func BytesSplit(a rel.Bytes, b rel.Value) rel.Value {
	var splitted []string

	switch b := b.(type) {
	case rel.Bytes:
		splitted = strings.Split(a.String(), b.String())
	case rel.GenericSet:
		splitted = strings.Split(a.String(), mustAsString(b))
	}

	vals := make([]rel.Value, 0, len(splitted))
	for _, s := range splitted {
		vals = append(vals, rel.NewBytes([]byte(s)).(rel.Value))
	}
	return rel.NewArray(vals...)
}
