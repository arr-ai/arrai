package pb

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// StdProtobufDecoder transforms the protocol buffer message to a tuple.
// Sample code to call this method:
// let sysl = //encoding.proto.decode(//encoding.proto.proto, //os.file('sysl.pb'));
// let decodeSyslPb = //encoding.proto.decode(sysl);
// let shop = decodeSyslPb('Module', //os.file("petshop.pb"));
// petshop.apps("PetShopApi")
var StdProtobufDecoder = rel.NewNativeFunction("decode", func(param rel.Value) (rel.Value, error) {
	bytes, isBytes := param.(rel.Bytes)
	if !isBytes {
		return nil, fmt.Errorf("//encoding.proto.decode: param not bytes")
	}

	if bytes.Bytes() == nil {
		// the call looks like `let sysl = //encoding.proto.decode(//encoding.proto.proto, //os.file('sysl.pb'))`
		return rel.NewNativeFunction("decode$1", func(definition rel.Value) (rel.Value, error) {
			_, is := tools.ValueAsBytes(definition)
			if !is {
				return nil, fmt.Errorf("//encoding.proto.decode: definition not bytes")
			}

			return definition, nil
		}), nil
	}

	// the call looks like `let decodeSyslPb = //encoding.proto.decode(sysl);`
	// after `let sysl = //encoding.proto.decode(//encoding.proto.proto, //os.file('sysl.pb'));`
	definitionBytes, is := tools.ValueAsBytes(bytes)
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

			tuple, err := convertProtoValToSyslVal(protoreflect.ValueOf(message.ProtoReflect()))
			if err != nil {
				return nil, err
			}

			return tuple, nil
		}), nil
	}), nil
})
