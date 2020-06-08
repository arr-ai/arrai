package syntax

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"
)

func createTestCompareFuncAttr(name string, ok func(a, b rel.Value) bool, message string) rel.Attr {
	return createNestedFuncAttr(name, 2, func(args ...rel.Value) (rel.Value, error) {
		expected := args[0]
		actual := args[1]
		if !ok(expected, actual) {
			return nil, fmt.Errorf("%s\nexpected: %v\nactual:   %v", message, expected, actual)
		}
		return rel.None, nil
	})
}

func createTestCheckFuncAttr(name string, ok func(v rel.Value) bool) rel.Attr {
	return rel.NewNativeFunctionAttr(name, func(value rel.Value) (rel.Value, error) {
		if ok(value) {
			return rel.None, nil
		}
		return nil, fmt.Errorf("not %s\nvalue: %v", name, value)
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
		//TODO: reimplement testing suite
		// rel.NewNativeExprFunctionAttr("suite", func(expr rel.Expr, local rel.Scope) (rel.Value, error) {
		// 	switch expr := expr.(type) {
		// 	case rel.Value:
		// 		return rel.None, nil
		// 	case *rel.SetExpr:
		// 		errors := []error{}
		// 		var filename string
		// 		for _, elt := range expr.Elements() {
		// 			var err error
		// 			func() {
		// 				defer func() {
		// 					switch r := recover().(type) {
		// 					case nil:
		// 					case error:
		// 						err = wrapContext(r, elt)
		// 					default:
		// 						panic(wrapContext(fmt.Errorf("unexpected panic: %v", r), expr))
		// 					}
		// 				}()
		// 				_, err = elt.Eval(local)
		// 			}()
		// 			if err != nil {
		// 				filename = elt.Source().Filename()
		// 				line, _ := elt.Source().Position()
		// 				errors = append(errors, fmt.Errorf("%d: %v", line, err))
		// 			}
		// 		}
		// 		if len(errors) > 0 {
		// 			var sb strings.Builder
		// 			for _, err := range errors {
		// 				fmt.Fprintln(&sb, err.Error())
		// 			}
		// 			return nil, fmt.Errorf("test failure(s) in %s:\n%s", filename, sb.String())
		// 		}
		// 		return rel.None, nil
		// 	default:
		// 		return nil, fmt.Errorf("//test.suite arg must be a set of tests")
		// 	}
		// }),
	)
}

// func wrapContext(err error, expr rel.Expr) error {
// 	return fmt.Errorf("%s\n%s", err.Error(), expr.Source().Context(parser.DefaultLimit))
// }
