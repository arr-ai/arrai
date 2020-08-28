package syntax

import (
	"context"
	"fmt"

	"github.com/arr-ai/arrai/rel"
)

func stdEval() rel.Attr {
	return rel.NewTupleAttr("eval",
		//TODO: eval needs to be changed to only evaluate simple expression
		// e.g. no functions, no math operations etc only simple values
		rel.NewNativeFunctionAttr("value", evalExpr),

		//TODO: eval.expr
	)
}

func evalExpr(ctx context.Context, v rel.Value) (rel.Value, error) {
	switch val := v.(type) {
	case rel.String, rel.Bytes:
		evaluated, err := EvaluateExpr(ctx, ".", val.String())
		if err != nil {
			panic(err)
		}
		return evaluated, nil
	}
	return nil, fmt.Errorf("//eval.value: not a byte array or string: %v", v)
}
