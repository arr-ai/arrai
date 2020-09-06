package arraictx

import (
	"context"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/pkg/ctxrootcache"
	"github.com/spf13/afero"
)

func InitRunCtx(ctx context.Context) context.Context {
	ctx = ctxfs.SourceFsOnto(ctx, afero.NewOsFs())
	ctx = ctxfs.RuntimeFsOnto(ctx, afero.NewOsFs())
	ctx = ctxrootcache.WithRootCache(ctx)
	return ctx
}
