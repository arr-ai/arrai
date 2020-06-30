package syntax

import (
	"fmt"

	"github.com/arr-ai/arrai/pb"

	"github.com/arr-ai/arrai/rel"
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
		createFunc3Attr("decode", stdProtobufDecode),
	)
}

func stdProtobufDecode(definition, data, rootMessageName rel.Value) (rel.Value, error) {
	definitionBytes, is := definition.(rel.Bytes)
	if !is {
		return nil, fmt.Errorf("//encoding.proto.decode: definition not bytes")
	}

	dataBytes, is := data.(rel.Bytes)
	if !is {
		return nil, fmt.Errorf("//encoding.proto.decode: data not bytes")
	}

	rootMessageNameStr, is := rootMessageName.(rel.String)
	if !is {
		return nil, fmt.Errorf("//encoding.proto.decode: rootMessageName not string")
	}

	str, isStr := valueAsString(rootMessageNameStr)
	if !isStr {
		return nil, fmt.Errorf("//encoding.proto.decode: rootMessageName not string")
	}

	return pb.TransformProtoBufToTuple(definitionBytes.Bytes(), dataBytes.Bytes(), str)
}
