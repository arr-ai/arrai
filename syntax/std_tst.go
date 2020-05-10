package syntax

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"
)

func createTestCompareFuncAttr(name string, ok func(a, b rel.Value) bool, message string) rel.Attr {
	return createNestedFuncAttr(name, 2, func(args ...rel.Value) rel.Value {
		expected := args[0]
		actual := args[1]
		if !ok(expected, actual) {
			panic(fmt.Errorf("%s\nexpected: %v\nactual:   %v", message, expected, actual))
		}
		return rel.None
	})
}

func createTestCheckFuncAttr(name string, ok func(v rel.Value) bool) rel.Attr {
	return rel.NewNativeFunctionAttr(name, func(value rel.Value) rel.Value {
		if !ok(value) {
			panic(fmt.Errorf("not %s\nvalue: %v", name, value))
		}
		return rel.None
	})
}

func stdTest() rel.Attr {
	return rel.NewTupleAttr("test",
		rel.NewTupleAttr("assert",
			createTestCompareFuncAttr("equal", func(e, a rel.Value) bool { return e.Equal(a) }, "not equal"),
			createTestCompareFuncAttr("unequal", func(e, a rel.Value) bool { return !e.Equal(a) }, "not unequal"),
			createTestCompareFuncAttr("size", func(e, a rel.Value) bool {
				return int(e.(rel.Number).Float64()) == a.(rel.Set).Count()
			}, "unexpected size"),
			createTestCheckFuncAttr("false", func(v rel.Value) bool { return !v.IsTrue() }),
			createTestCheckFuncAttr("true", func(v rel.Value) bool { return v.IsTrue() }),
		),
		rel.NewNativeFunctionAttr("suite", func(value rel.Value) rel.Value {
			return nil
		}),
	)
}
