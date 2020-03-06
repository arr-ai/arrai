package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

func stdStrconv() rel.Attr {
	return rel.NewAttr("strconv", rel.NewTuple(
		//TODO: eval needs to be changed to only evaluate simple expression
		// e.g. no functions, no math operations etc only simple values
		rel.NewNativeFunctionAttr("eval", unsafeEval),
		rel.NewNativeFunctionAttr("unsafe_eval", unsafeEval),
	))
}

func unsafeEval(v rel.Value) rel.Value {
	evaluated, err := EvaluateExpr(".", v.(rel.String).String())
	if err != nil {
		panic(err)
	}
	return evaluated
}
