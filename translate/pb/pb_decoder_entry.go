package pb

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

//StdProtobufDecoder transforms the protocol buffer message to a tuple.
var StdProtobufDecoder = rel.NewNativeFunction("decode", func(tupleParam rel.Value) (rel.Value, error) {
	tuple, isTuple := tupleParam.(rel.Tuple)
	if !isTuple {
		return nil, fmt.Errorf("//encoding.proto.decode: tupleParam not tuple")
	}

	if !tuple.IsTrue() {
		// the call looks like `/set sysl = //encoding.proto.decode(//encoding.proto.proto, //os.file("sysl.pb"))`
		return rel.NewNativeFunction("decode$1", func(definition rel.Value) (rel.Value, error) {
			_, is := tools.ValueAsBytes(definition)
			if !is {
				return nil, fmt.Errorf("//encoding.proto.decode: definition not bytes")
			}

			return rel.NewTuple(rel.NewAttr("fileDescriptor", definition)), nil
		}), nil
	}

	// the call looks like `/set decodeSyslPb = //encoding.proto.decode(sysl)`
	// after `/set sysl = //encoding.proto.decode(//encoding.proto.proto, //os.file("sysl.pb"))`
	fileDescriptor, has := tuple.Get("fileDescriptor")
	if !has {
		return nil, fmt.Errorf("//encoding.proto.decode: tupleParam doesn't have protobuf file descriptor")
	}

	definitionBytes, is := tools.ValueAsBytes(fileDescriptor)
	if !is {
		return nil, fmt.Errorf("//encoding.proto.decode: fileDescriptor not bytes")
	}

	return rel.NewNativeFunction("decode$2", func(rootMessageName rel.Value) (rel.Value, error) {
		fd, err := decodeFileDescriptor(definitionBytes)
		if err != nil {
			return nil, err
		}

		nameStr, isStr := tools.ValueAsString(rootMessageName)
		if !isStr {
			return nil, fmt.Errorf("//encoding.proto.decode: rootMessageName not string")
		}
		rootMessageDesc := fd.Messages().ByName(protoreflect.Name(nameStr))
		message := dynamicpb.NewMessage(rootMessageDesc)

		return rel.NewNativeFunction("decode$3", func(data rel.Value) (rel.Value, error) {
			dataBytes, is := tools.ValueAsBytes(data)
			if !is {
				return nil, fmt.Errorf("//encoding.proto.decode: data not bytes")
			}

			err = proto.Unmarshal(dataBytes, message)
			if err != nil {
				return nil, err
			}

			tuple, err := walkThroughMessageToBuildValue(protoreflect.ValueOf(message.ProtoReflect()))
			if err != nil {
				return nil, err
			}

			return tuple, nil
		}), nil
	}), nil
})
