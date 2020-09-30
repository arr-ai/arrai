package tools

import (
	"github.com/arr-ai/arrai/rel"
)

// ValueAsString transform rel.Value to string.
// In arr.ai, all empty sets are the same and `""`, `{}` and `[]` will be parsed to rel.GenericSet.
func ValueAsString(v rel.Value) (string, bool) {
	switch v := v.(type) {
	case rel.String:
		return v.String(), true
	case rel.GenericSet:
		return "", !v.IsTrue()
	}
	return "", false
}

// ValueAsBytes transform rel.Value to byte array.
// In arr.ai, all empty sets are the same and `""`, `{}` and `[]` will be parsed to rel.GenericSet.
func ValueAsBytes(v rel.Value) ([]byte, bool) {
	switch v := v.(type) {
	case rel.Bytes:
		return v.Bytes(), true
	case rel.GenericSet:
		return nil, !v.IsTrue()
	}
	return nil, false
}
