package tools

import (
	"fmt"

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

// ValueTypeAsString returns a string that describes the type of the value in a human-readable form.
func ValueTypeAsString(v rel.Value) string {
	switch v.(type) {
	case rel.Number, *rel.Number:
		return "number"
	case rel.Array, *rel.Array:
		return "array"
	case rel.Bytes, *rel.Bytes:
		return "bytes"
	case rel.Closure, *rel.Closure:
		return "closure"
	case rel.Dict, *rel.Dict:
		return "dict"
	case rel.ExprClosure, *rel.ExprClosure:
		return "expr closure"
	case rel.String, *rel.String:
		return "string"
	case rel.GenericSet, *rel.GenericSet:
		return "set"
	case rel.Tuple, *rel.GenericTuple:
		return "tuple"
	case rel.ArrayItemTuple, *rel.ArrayItemTuple:
		return "array item tuple"
	case rel.BytesByteTuple, *rel.BytesByteTuple:
		return "bytes byte tuple"
	case rel.DictEntryTuple, *rel.DictEntryTuple:
		return "dict entry tuple"
	case rel.StringCharTuple, *rel.StringCharTuple:
		return "string char tuple"
	}
	return fmt.Sprintf("%T", v)
}
