package syntax

import (
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate/pb"
)

func stdEncoding() rel.Attr {
	return rel.NewTupleAttr("encoding",
		stdEncodingJSON(),
		stdEncodingProto(),
	)
}

func stdEncodingProto() rel.Attr {
	return rel.NewTupleAttr(
		"proto",
		rel.NewAttr("decode", pb.StdProtobufDecoder),
		// this is a placeholder to represent //encoding.proto.proto
		rel.NewAttr("proto", rel.NewTuple()),
	)
}
