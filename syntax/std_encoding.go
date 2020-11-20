package syntax

import (
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/translate/pb"
)

func stdEncoding() rel.Attr {
	return rel.NewTupleAttr("encoding",
		stdEncodingCSV(),
		stdEncodingJSON(),
		stdEncodingProto(),
		stdEncodingXlsx(),
		stdEncodingYAML(),
		stdEncodingXML(),
	)
}

func stdEncodingProto() rel.Attr {
	return rel.NewTupleAttr(
		"proto",
		rel.NewAttr("decode", pb.StdProtobufDecoder),
		rel.NewAttr("descriptor", pb.StdProtobufDescriptor),
	)
}
