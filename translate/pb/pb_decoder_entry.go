package pb

import (
	"context"
	"fmt"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

const fileDescriptorSet = "FileDescriptorSet"

// StdProtobufDecoder transforms the protocol buffer message to a tuple.
var StdProtobufDecoder = rel.NewNativeFunction("decode", func(_ context.Context, param rel.Value) (rel.Value, error) {
	tuple, isTuple := param.(rel.Tuple)
	if !isTuple {
		return nil, fmt.Errorf("//encoding.proto.decode: param not tuple")
	}
	// TODO will change it to use a tuple to represent a file descriptor, #496
	fdVal, found := tuple.Get(fileDescriptorSet)
	if !found {
		return nil, fmt.Errorf("//encoding.proto.decode: couldn't find %s in tuple", fileDescriptorSet)
	}
	fdBytes, isBytes := fdVal.(rel.Bytes)
	if !isBytes {
		return nil, fmt.Errorf("//encoding.proto.decode: %s is not bytes", fileDescriptorSet)
	}
	fd, err := decodeFileDescriptor(fdBytes.Bytes())
	if err != nil {
		return nil, err
	}

	return rel.NewNativeFunction("decode$2", func(_ context.Context, messageTypeName rel.Value) (rel.Value, error) {
		nameStr, isStr := tools.ValueAsString(messageTypeName)
		if !isStr {
			return nil, fmt.Errorf("//encoding.proto.decode: messageTypeName not string")
		}
		rootMessageDesc := fd.Messages().ByName(protoreflect.Name(nameStr))
		message := dynamicpb.NewMessage(rootMessageDesc)

		return rel.NewNativeFunction("decode$3", func(_ context.Context, data rel.Value) (rel.Value, error) {
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

// StdProtobufDescriptor transforms the protocol buffer `.proto` binary file to a tuple.
var StdProtobufDescriptor = rel.NewNativeFunction(
	"decode",
	func(_ context.Context, param rel.Value) (rel.Value, error) {
		definitionBytes, isBytes := param.(rel.Bytes)
		if !isBytes {
			return nil, fmt.Errorf("//encoding.proto.descriptor: param not bytes")
		}

		return rel.NewTuple(rel.NewAttr(fileDescriptorSet, definitionBytes)), nil
	},
)
