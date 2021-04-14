package test

import (
	"fmt"
	"github.com/arr-ai/arrai/rel"
	"strings"
)

func ForeachLeaf(val rel.Value, path string, leafAction func(val rel.Value, path string)) {
	path = strings.TrimPrefix(path, "<root>.")

	if isLiteralTrue(val) || isLiteralFalse(val) {
		leafAction(val, path)
		return
	}

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
		leafAction(val, path)
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
