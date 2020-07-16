package syntax

import (
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate/pb"
)

func stdEncoding() rel.Attr {
	return rel.NewTupleAttr("encoding",
		stdEncodingJSON(),
		stdEncodingProto(),
		stdEncodingYAML(),
	)
}

func stdEncodingProto() rel.Attr {
	return rel.NewTupleAttr(
		"proto",
		rel.NewAttr("decode", pb.StdProtobufDecoder),
		// this is a placeholder to represent `//encoding.proto.proto`
		// Now profobuf decoder method is `//encoding.proto.decode`, and the first argument is `rel.Bytes`.
		// In order to implement call `//encoding.proto.decode( , //os.file('sysl.pb'))`, so add this placeholder.
		// The full call can be `//encoding.proto.decode(//encoding.proto.proto , //os.file('sysl.pb'))`.
		rel.NewAttr("proto", rel.Bytes{}),
	)
}
