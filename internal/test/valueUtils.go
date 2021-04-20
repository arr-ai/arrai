package test

import (
	"fmt"
	"github.com/arr-ai/arrai/rel"
	"strings"
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

var emptyTuple = rel.NewTuple()

func isLiteralTrue(val rel.Value) bool {
	v, ok := val.(rel.GenericSet)
	return ok && v.Count() == 1 && v.Has(emptyTuple)
}

func isLiteralFalse(val rel.Value) bool {
	v, ok := val.(rel.GenericSet)
	return ok && v.Count() == 0
}
