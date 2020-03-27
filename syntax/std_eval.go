package syntax

import "github.com/arr-ai/arrai/rel"

func stdEval() rel.Attr {
	return rel.NewTupleAttr("eval",
		//TODO: eval needs to be changed to only evaluate simple expression
		// e.g. no functions, no math operations etc only simple values
		rel.NewNativeFunctionAttr("value", evalExpr),

		//TODO: eval.expr
	)
}

func evalExpr(v rel.Value) rel.Value {
	evaluated, err := EvaluateExpr(".", v.(rel.String).String())
	if err != nil {
		panic(err)
	}
	return evaluated
}
