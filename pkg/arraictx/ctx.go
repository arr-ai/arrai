package arraictx

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/pkg/ctxrootcache"
	"github.com/spf13/afero"
)

type ctxKey int

const (
	argsKey ctxKey = iota
)

// InitCliCtx returns an arr.ai context with the arguments set from the CLI context.
//
// For example, command line `arrai -d r file.arrai arg1 arg2 arg3` will set
// `[]string{"file.arrai", "arg1", "arg2", "arg3"}` as the context value of `argsKey`.
func InitCliCtx(ctx context.Context, c *cli.Context) context.Context {
	return WithArgs(InitRunCtx(ctx), c.Args().Slice()...)
}

// InitRunCtx returns a context for evaluating arr.ai programs.
func InitRunCtx(ctx context.Context) context.Context {
	ctx = ctxfs.SourceFsOnto(ctx, afero.NewOsFs())
	ctx = ctxfs.RuntimeFsOnto(ctx, afero.NewOsFs())
	ctx = ctxrootcache.WithRootCache(ctx)
	return ctx
}

// WithArgs sets the CLI arguments on the Go context.
func WithArgs(ctx context.Context, args ...string) context.Context {
	return context.WithValue(ctx, argsKey, args)
}

// Args returns the args stored in ctx.
func Args(ctx context.Context) []string {
	a, ok := ctx.Value(argsKey).([]string)
	if !ok {
		return []string{}
	}
	return a
}
