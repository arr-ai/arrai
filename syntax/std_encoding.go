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

func stdEncodingBytesOrStringAsUTF8(v rel.Value) ([]byte, bool) {
	var bytes []byte
	switch v := v.(type) {
	case rel.String:
		bytes = []byte(v.String())
	case rel.Bytes:
		bytes = v.Bytes()
	default:
		return nil, false
	}

	return bytes, true
}
