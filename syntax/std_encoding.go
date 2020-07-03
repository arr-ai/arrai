package syntax

import (
	"fmt"

	"github.com/arr-ai/arrai/translate/pb"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"

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
		rel.NewAttr("decode", stdProtobufDecode),
	)
}

//StdProtobufDecode transforms the protocol buffer message to a tuple.
var stdProtobufDecode = rel.NewNativeFunction("decode", func(definition rel.Value) (rel.Value, error) {
	definitionBytes, is := valueAsBytes(definition)
	if !is {
		return nil, fmt.Errorf("//encoding.proto.decode: definition not bytes")
	}
	fd, err := pb.DecodeFileDescriptor(definitionBytes)
	if err != nil {
		return nil, err
	}

	return rel.NewNativeFunction("decode$1", func(rootMessageName rel.Value) (rel.Value, error) {
		nameStr, isStr := valueAsString(rootMessageName)
		if !isStr {
			return nil, fmt.Errorf("//encoding.proto.decode: rootMessageName not string")
		}
		rootMessageDesc := fd.Messages().ByName(protoreflect.Name(nameStr))
		message := dynamicpb.NewMessage(rootMessageDesc)

		return rel.NewNativeFunction("decode$2", func(data rel.Value) (rel.Value, error) {
			dataBytes, is := valueAsBytes(data)
			if !is {
				return nil, fmt.Errorf("//encoding.proto.decode: data not bytes")
			}

			err = proto.Unmarshal(dataBytes, message)
			if err != nil {
				return nil, err
			}

			tuple, err := pb.WalkThroughMessageToBuildValue(protoreflect.ValueOf(message.ProtoReflect()))
			if err != nil {
				return nil, err
			}

			return tuple, nil
		}), nil
	}), nil
})
