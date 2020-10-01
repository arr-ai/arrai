package syntax

import (
	"context"

	"github.com/arr-ai/arrai/pkg/arraictx"

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

// EvaluateExpr evaluate the passed in arrai script `source` and returns the evaluated arrai value.
// Parameter `path` is used as source context, could be empty.
func EvaluateExpr(ctx context.Context, path, source string) (rel.Value, error) {
	return EvalWithScope(ctx, path, source, rel.Scope{})
}

// EvaluateBundle evaluates the buffer of a bundled scripts using the arrai bundle cmd.
// If args are provided, they override the values of //os.args.
func EvaluateBundle(bundle []byte, args ...string) (rel.Value, error) {
	ctx := arraictx.InitRunCtx(context.Background())
	return EvaluateBundleCtx(ctx, bundle, args...)
}

// EvaluateBundle evaluates the buffer of a bundled scripts using the arrai bundle cmd.
// If args are provided, they override the values of //os.args.
func EvaluateBundleCtx(ctx context.Context, bundle []byte, args ...string) (rel.Value, error) {
	ctx = arraictx.WithArgs(ctx, args...)
	ctx, err := WithBundleRun(ctx, bundle)
	if err != nil {
		return nil, err
	}
	ctx, mainFileSource, path := GetMainBundleSource(ctx)
	return EvaluateExpr(ctx, path, string(mainFileSource))
}
