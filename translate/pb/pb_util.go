package pb

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/tools"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

// decodeFileDescriptor parses protobuf definition file and decodes to file descriptor
func decodeFileDescriptor(definition []byte) (protoreflect.FileDescriptor, error) {
	fdSet := new(descriptorpb.FileDescriptorSet)
	if err := proto.Unmarshal(definition, fdSet); err != nil {
		return nil, err
	}

	// We know protoc was invoked with a single .proto file only
	file1 := fdSet.GetFile()[0]

	// Initialize the File descriptor object
	fd1, err := protodesc.NewFile(file1, protoregistry.GlobalFiles)
	if err != nil {
		return nil, err
	}

	return fd1, nil
}

// WalkThroughMessageToBuildValue walks through protobuf message and build a value whose type is rel.Value
func WalkThroughMessageToBuildValue(val protoreflect.Value) (rel.Value, error) {
	switch message := val.Interface().(type) {
	case protoreflect.Message:
		attrs := []rel.Attr{}
		message.Range(func(desc protoreflect.FieldDescriptor, val protoreflect.Value) bool {
			item, err := WalkThroughMessageToBuildValue(val)
			if err != nil {
				attrs = append(attrs, rel.NewAttr("@error_"+string(desc.Name()),
					rel.NewString([]rune(err.Error()))))
			} else {
				if desc.Enum() != nil {
					switch t := item.(type) {
					case rel.Number:
						num, success := t.Int()
						if success {
							name := desc.Enum().Values().ByNumber(protoreflect.EnumNumber(num)).Name()
							attrs = append(attrs,
								rel.NewAttr(string(desc.Name()),
									rel.NewString([]rune(string(name)))))
						} else {
							attrs = append(attrs, rel.NewAttr(string(desc.Name()), item))
						}
					default:
						attrs = append(attrs, rel.NewAttr(string(desc.Name()), item))
					}
				} else {
					attrs = append(attrs, rel.NewAttr(string(desc.Name()), item))
				}
			}
			return true
		})
		return rel.NewTuple(attrs...), nil

	case protoreflect.Map:
		entries := []rel.DictEntryTuple{}
		val.Map().Range(func(key protoreflect.MapKey, value protoreflect.Value) bool {
			val, err := WalkThroughMessageToBuildValue(value)
			if err != nil {
				entries = append(entries, rel.NewDictEntryTuple(
					rel.NewString([]rune("@error_"+key.String())),
					rel.NewString([]rune(err.Error()))))
			} else {
				entries = append(entries, rel.NewDictEntryTuple(
					rel.NewString([]rune(key.String())),
					val))
			}
			return true
		})
		return rel.NewDict(false, entries...), nil

	case protoreflect.List:
		list := []rel.Value{}
		len := message.Len()
		for i := 0; i < len; i++ {
			item, err := WalkThroughMessageToBuildValue(message.Get(i))
			if err != nil {
				list = append(list, rel.NewArrayItemTuple(i,
					rel.NewString([]rune("@error:"+err.Error()))))
			} else {
				list = append(list, rel.NewArrayItemTuple(i, item))
			}
		}
		return rel.NewArray(list...), nil

	case string:
		return rel.NewString([]rune(message)), nil
	case bool:
		return rel.NewBool(message), nil
	case int32, int64, uint32, uint64, float32, float64:
		val, err := rel.NewValue(message)
		if err != nil {
			return nil, err
		}
		return val, nil
	case protoreflect.EnumNumber:
		// type EnumNumber int32
		val, err := rel.NewValue(int32(message))
		if err != nil {
			return nil, err
		}
		return val, nil
	default:
		// []byte and protoreflect.Enum etc.
		return rel.NewTuple(rel.NewStringAttr("@error",
			[]rune(fmt.Errorf("%T is not supported data type", message).Error()))), nil
	}
}

//StdProtobufDecode transforms the protocol buffer message to a tuple.
var StdProtobufDecode = rel.NewNativeFunction("decode", func(definition rel.Value) (rel.Value, error) {
	definitionBytes, is := tools.ValueAsBytes(definition)
	if !is {
		return nil, fmt.Errorf("//encoding.proto.decode: definition not bytes")
	}
	fd, err := decodeFileDescriptor(definitionBytes)
	if err != nil {
		return nil, err
	}

	return rel.NewNativeFunction("decode$1", func(rootMessageName rel.Value) (rel.Value, error) {
		nameStr, isStr := tools.ValueAsString(rootMessageName)
		if !isStr {
			return nil, fmt.Errorf("//encoding.proto.decode: rootMessageName not string")
		}
		rootMessageDesc := fd.Messages().ByName(protoreflect.Name(nameStr))
		message := dynamicpb.NewMessage(rootMessageDesc)

		return rel.NewNativeFunction("decode$2", func(data rel.Value) (rel.Value, error) {
			dataBytes, is := tools.ValueAsBytes(data)
			if !is {
				return nil, fmt.Errorf("//encoding.proto.decode: data not bytes")
			}

			err = proto.Unmarshal(dataBytes, message)
			if err != nil {
				return nil, err
			}

			tuple, err := WalkThroughMessageToBuildValue(protoreflect.ValueOf(message.ProtoReflect()))
			if err != nil {
				return nil, err
			}

			return tuple, nil
		}), nil
	}), nil
})
