package syntax

import (
	"context"

	"github.com/arr-ai/arrai/pkg/buildinfo"
	"github.com/arr-ai/arrai/rel"
)

func stdRuntime() rel.Attr {
	return rel.NewTupleAttr("arrai",
		rel.NewNativeFunctionAttr("info", func(context.Context, rel.Value) (rel.Value, error) {
			//TODO
			return GetBuildInfo(buildinfo.BuildInfo), nil
		}),
	)
}

// GetBuildInfo returns arr.ai build information.
func GetBuildInfo(b buildinfo.BuildData) rel.Value {
	return rel.NewTuple(
		rel.NewAttr("version", rel.NewGoString(b.Version)),
		rel.NewAttr("date", rel.NewGoString(b.Date)),
		rel.NewAttr("git", rel.NewTuple(
			rel.NewAttr("commit", rel.NewGoString(b.FullCommit)),
			// param tags has only one tag now.
			rel.NewAttr("tags", rel.MustNewSet(rel.NewGoString(b.Tags)))),
		),
		rel.NewAttr("go", rel.NewTuple(
			rel.NewAttr("os", rel.NewGoString(b.Os)),
			rel.NewAttr("arch", rel.NewGoString(b.Arch)),
			rel.NewAttr("compiler", rel.NewTuple(
				rel.NewAttr("version", rel.NewGoString(b.GoVersion)),
				rel.NewAttr("os", rel.NewGoString(b.Os)),
				rel.NewAttr("arch", rel.NewGoString(b.Arch))),
			)),
		),
	)
}
