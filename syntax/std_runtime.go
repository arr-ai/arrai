package syntax

import (
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
)

func stdRuntime() rel.Attr {
	return rel.NewTupleAttr("arrai",
		rel.NewAttr("info", rel.NewString([]rune("Hello"))),
	)
}

func buildInfo() {
	rel.NewTuple(rel.NewAttr("version", rel.NewString([]rune(tools.Version))))
}
