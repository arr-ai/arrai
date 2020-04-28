package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

func EvalWithScope(path, source string, scope rel.Scope) (rel.Value, error) {
	expr, err := Compile(path, source)
	if err != nil {
		return nil, err
	}

	value, err := expr.Eval(scope)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func EvaluateExpr(path, source string) (rel.Value, error) {
	return EvalWithScope(path, source, rel.Scope{})
}
