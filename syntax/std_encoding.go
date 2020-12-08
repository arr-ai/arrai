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

func bytesOrStringAsUTF8(v rel.Value) ([]byte, bool) {
	switch v := v.(type) {
	case rel.String:
		return []byte(v.String()), true
	case rel.Bytes:
		return v.Bytes(), true
	default:
		return nil, false
	}
}
