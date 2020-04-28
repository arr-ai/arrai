package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

func EvaluateExpr(path, source string) (rel.Value, error) {
	expr, err := Compile(path, source)
	if err != nil {
		return nil, err
	}

	value, err := expr.Eval(rel.Scope{})
	if err != nil {
		return nil, err
	}

	return value, nil
}
