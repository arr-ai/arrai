package syntax

import (
	"context"
	"fmt"

	"github.com/arr-ai/arrai/rel"
)

func createTestCompareFuncAttr(name string, ok func(_ context.Context, a, b rel.Value) bool, message string) rel.Attr {
	return createNestedFuncAttr(name, 2, func(ctx context.Context, args ...rel.Value) (rel.Value, error) {
		expected := args[0]
		actual := args[1]
		if !ok(ctx, expected, actual) {
			return nil, fmt.Errorf("%s\nexpected: %v\nactual:   %v", message, expected, actual)
		}
		return rel.True, nil
	})
}

func createTestCheckFuncAttr(name string, ok func(_ context.Context, v rel.Value) bool) rel.Attr {
	return rel.NewNativeFunctionAttr(name, func(ctx context.Context, value rel.Value) (rel.Value, error) {
		if ok(ctx, value) {
			return rel.True, nil
		}
		return nil, fmt.Errorf("not %s\nvalue: %v", name, value)
	})
}

func stdTest() rel.Attr {
	return rel.NewTupleAttr("test",
		rel.NewTupleAttr("assert",
			createTestCompareFuncAttr(
				"equal",
				func(_ context.Context, e, a rel.Value) bool { return e.Equal(a) },
				"not equal",
			),
			createTestCompareFuncAttr(
				"unequal",
				func(_ context.Context, e, a rel.Value) bool { return !e.Equal(a) },
				"not unequal",
			),
			createTestCompareFuncAttr("size", func(_ context.Context, e, a rel.Value) bool {
				return int(e.(rel.Number).Float64()) == a.(rel.Set).Count()
			}, "unexpected size"),
			createTestCheckFuncAttr("false", func(_ context.Context, v rel.Value) bool { return !v.IsTrue() }),
			createTestCheckFuncAttr("true", func(_ context.Context, v rel.Value) bool { return v.IsTrue() }),
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
		// 						err = WrapContext(r, elt)
		// 					default:
		// 						panic(WrapContext(fmt.Errorf("unexpected panic: %v", r), expr))
		// 					}
		// 				}()
		// 				_, err = elt.Eval(ctx, local)
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

// func WrapContext(err error, expr rel.Expr) error {
// 	return fmt.Errorf("%s\n%s", err.Error(), expr.Source().Context(parser.DefaultLimit))
// }
