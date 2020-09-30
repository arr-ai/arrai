package rel

import (
	"fmt"
)

// ValueTypeAsString returns a string that describes the type of the value in a human-readable form.
func ValueTypeAsString(v Value) string {
	switch v.(type) {
	case Number, *Number:
		return "number"
	case Array, *Array:
		return "array"
	case Bytes, *Bytes:
		return "bytes"
	case Closure, *Closure:
		return "closure"
	case ExprClosure, *ExprClosure:
		return "expr-closure"
	case Dict, *Dict:
		return "dict"
	case String, *String:
		return "string"
	case GenericSet, *GenericSet:
		return "set"
	case ArrayItemTuple, *ArrayItemTuple:
		return "array-item-tuple"
	case BytesByteTuple, *BytesByteTuple:
		return "bytes-byte-tuple"
	case DictEntryTuple, *DictEntryTuple:
		return "dict-entry-tuple"
	case StringCharTuple, *StringCharTuple:
		return "string-char-tuple"
	case Tuple, *GenericTuple:
		return "tuple"
	}
	return fmt.Sprintf("%T", v)
}
