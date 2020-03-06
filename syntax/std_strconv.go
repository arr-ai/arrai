package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

func stdStrconv() rel.Attr {
	return rel.NewAttr("strconv", rel.NewTuple(
		rel.NewNativeFunctionAttr("eval", func(v rel.Value) rel.Value {
			evaluated, err := EvaluateExpr(".", v.(rel.String).String())
			if err != nil {
				panic(err)
			}
			return evaluated
		}),
	))
}
