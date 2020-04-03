package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

func stdEncoding() rel.Attr {
	return rel.NewTupleAttr("encoding", stdEncodingJSON())
}
