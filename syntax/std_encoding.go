package syntax

import (
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate/pb"
)

func stdEncoding() rel.Attr {
	return rel.NewTupleAttr("encoding",
		stdEncodingJSON(),
		stdEncodingProtoProto(),
		stdEncodingProto(),
	)
}

// this is a placeholder to represent //encoding.proto.proto
func stdEncodingProtoProto() rel.Attr {
	return rel.NewTupleAttr(
		"proto",
		rel.NewAttr("proto", rel.NewTuple()),
	)
}

func stdEncodingProto() rel.Attr {
	return rel.NewTupleAttr(
		"proto",
		rel.NewAttr("decode", pb.StdProtobufDecode),
	)
}
