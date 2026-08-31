package syntax

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/arrai/pkg/importcache"

	"github.com/arr-ai/arrai/rel"
)

// timingEnabled reports whether ARRAI_TIMING is set, in which case every
// top-level evaluation reports its compile and eval phases to stderr.
// Compile covers parsing and compiling the program and everything it
// imports; eval is the program actually running. Separating them keeps
// language-startup cost out of algorithm measurements.
var timingEnabled = os.Getenv("ARRAI_TIMING") != ""

func timePhase(phase, path string) func() {
	if !timingEnabled {
		return func() {}
	}
	start := time.Now()
	return func() {
		fmt.Fprintf(os.Stderr, "arrai timing: %s %s (%s)\n",
			phase, time.Since(start).Round(time.Millisecond), path)
	}
}

func EvalWithScope(ctx context.Context, path, source string, scope rel.Scope) (rel.Value, error) {
	if !importcache.HasImportCacheFrom(ctx) {
		ctx = importcache.WithNewImportCache(ctx)
	}

	done := timePhase("compile", path)
	expr, err := Compile(ctx, path, source)
	done()
	if err != nil {
		return nil, err
	}

	done = timePhase("eval", path)
	value, err := expr.Eval(arraictx.ContextWithIsCompiling(ctx, false), scope)
	done()
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

// EvaluateBundleCtx evaluates the buffer of a bundled scripts using the arrai bundle cmd.
// If args are provided, they override the values of //os.args.
func EvaluateBundleCtx(ctx context.Context, bundle []byte, args ...string) (rel.Value, error) {
	if len(args) > 0 {
		ctx = arraictx.WithArgs(ctx, args...)
	}
	ctx, err := WithBundleRun(ctx, bundle)
	if err != nil {
		return nil, err
	}
	ctx, mainFileSource, path := GetMainBundleSource(ctx)
	return EvaluateExpr(ctx, path, string(mainFileSource))
}
