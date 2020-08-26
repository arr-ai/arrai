package syntax

import (
	"context"

	"github.com/arr-ai/arrai/rel"
)

func EvalWithScope(ctx context.Context, path, source string, scope rel.Scope) (rel.Value, error) {
	expr, err := Compile(ctx, path, source)
	if err != nil {
		return nil, err
	}

	value, err := expr.Eval(ctx, scope)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func EvaluateExpr(ctx context.Context, path, source string) (rel.Value, error) {
	return EvalWithScope(ctx, path, source, rel.Scope{})
}
