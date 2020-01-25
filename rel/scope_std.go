package rel

import "math"

var stdScope = EmptyScope.
	With(".", NewTupleExpr(
		MustNewAttrExpr("math", NewTupleExpr(
			MustNewAttrExpr("pi", NewNumber(math.Pi)),
			newFloatFuncAttr("sin", math.Sin),
			newFloatFuncAttr("cos", math.Cos),
		)),
		MustNewAttrExpr("grammar", NewTupleExpr(
			NewNativeFunctionAttrExpr("parse", func(value Value) Value { panic("not implemented") }),
		)),
	))

func newFloatFuncAttr(name string, f func(float64) float64) AttrExpr {
	return NewNativeFunctionAttrExpr(name, func(value Value) Value {
		return NewNumber(f(value.(Number).Float64()))
	})
}
