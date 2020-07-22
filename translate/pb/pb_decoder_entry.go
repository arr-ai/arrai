package pb

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

var StdProtobufDecoder = rel.NewNativeFunction("decode", func(param rel.Value) (rel.Value, error) {
	tuple, isTuple := param.(rel.Tuple)
	if !isTuple {
		return nil, fmt.Errorf("//encoding.proto.decode: param not tuple")
	}

	fdVal, _ := tuple.Get(fileDescriptorSet)
	fdBytes, _ := fdVal.(rel.Bytes)
	fd, _ := decodeFileDescriptor(fdBytes.Bytes())

	return rel.NewNativeFunction("decode$2", func(messageTypeName rel.Value) (rel.Value, error) {
		nameStr, isStr := tools.ValueAsString(messageTypeName)
		if !isStr {
			return nil, fmt.Errorf("//encoding.proto.decode: messageTypeName not string")
		}
		rootMessageDesc := fd.Messages().ByName(protoreflect.Name(nameStr))
		message := dynamicpb.NewMessage(rootMessageDesc)

		return rel.NewNativeFunction("decode$3", func(data rel.Value) (rel.Value, error) {
			dataBytes, is := tools.ValueAsBytes(data)
			if !is {
				return nil, fmt.Errorf("//encoding.proto.decode: data not bytes")
			}

			err := proto.Unmarshal(dataBytes, message)
			if err != nil {
				return nil, err
			}

			tuple, err := convertProtoValToSyslVal(protoreflect.ValueOf(message.ProtoReflect()))
			if err != nil {
				return nil, err
			}

			return tuple, nil
		}), nil
	}), nil
})

var StdProtobufDescriptor = rel.NewNativeFunction("decode", func(param rel.Value) (rel.Value, error) {
	definitionBytes, isBytes := param.(rel.Bytes)
	if !isBytes {
		return nil, fmt.Errorf("//encoding.proto.descriptor: param not bytes")
	}

	return rel.NewTuple(rel.NewAttr(fileDescriptorSet, definitionBytes)), nil
})

const fileDescriptorSet = "FileDescriptorSet"
