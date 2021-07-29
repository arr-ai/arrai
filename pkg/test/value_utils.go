package test

import (
	"fmt"
	"strings"

	"github.com/arr-ai/arrai/rel"
)

// ForeachLeaf visits all leaves in an test tree, invoking the leafAction callback for each leaf encountered.
// Tuples, arrays and dictionaries are considered test containers and not leaves, they are recursed into.
func ForeachLeaf(val rel.Value, path string, leafAction func(val rel.Value, path string)) {
	path = strings.TrimPrefix(path, ".")

	switch v := val.(type) {
	case rel.Array:
		for i, item := range v.Values() {
			ForeachLeaf(item, fmt.Sprintf("%s(%d)", path, i), leafAction)
		}
	case rel.Dict:
		for _, entry := range v.OrderedEntries() {
			key := entry.MustGet("@")
			keyStr := key.String()
			if _, ok := key.(rel.String); ok {
				keyStr = "'" + keyStr + "'"
			}
			ForeachLeaf(entry.MustGet(rel.DictValueAttr), fmt.Sprintf("%s(%s)", path, keyStr), leafAction)
		}
	case rel.Tuple:
		for e := v.Enumerator(); e.MoveNext(); {
			name, attr := e.Current()
			ForeachLeaf(attr, fmt.Sprintf("%s.%s", path, name), leafAction)
		}
	default:
		if path == "" {
			leafAction(val, "<root>")
		} else {
			leafAction(val, path)
		}
	}
}

// isLiteralTrue returns true if and only if the provided value is a literal true, i.e. "true" or "{()}"
func isLiteralTrue(val rel.Value) bool {
	switch v := val.(type) {
	case rel.TrueSet:
		return true
	case rel.GenericSet:
		if v.Count() == 1 && v.Has(rel.EmptyTuple) {
			panic(fmt.Errorf("true set is not of type TrueSet: %v", val))
		}
	}
	return false
}

// isLiteralFalse returns true if and only if the provided value is a literal false, i.e. "false" or "{}"
func isLiteralFalse(val rel.Value) bool {
	if _, ok := val.(rel.EmptySet); ok {
		return true
	}
	v, ok := val.(rel.GenericSet)
	return ok && v.Count() == 0
}
