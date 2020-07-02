package pb

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	pr "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

func getFileDescriptor(definition []byte) (pr.FileDescriptor, error) {
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

//TransformProtoBufToTuple transforms the protocol buffer message to a tuple.
func TransformProtoBufToTuple(definition []byte, data []byte, rootMessage string) (rel.Value, error) {
	fd, err := getFileDescriptor(definition)
	if err != nil {
		return nil, err
	}

	rootMessageDesc := fd.Messages().ByName(pr.Name(rootMessage))
	message := dynamicpb.NewMessage(rootMessageDesc)
	err = proto.Unmarshal(data, message)
	if err != nil {
		return nil, err
	}

	tuple, err := walkThroughMessage(pr.ValueOf(message.ProtoReflect()))
	if err != nil {
		return nil, err
	}

	return tuple, nil
}

func walkThroughMessage(val pr.Value) (rel.Value, error) {
	switch message := val.Interface().(type) {
	case pr.Message:
		// Only Message type is built to Tuple
		attrs := []rel.Attr{}
		message.Range(func(desc pr.FieldDescriptor, val pr.Value) bool {
			item, err := walkThroughMessage(val)
			if err != nil {
				attrs = append(attrs, rel.NewAttr("@error_"+string(desc.Name()),
					rel.NewString([]rune(err.Error()))))
			} else {
				attrs = append(attrs, rel.NewAttr(string(desc.Name()), item))
			}
			return true
		})
		return rel.NewTuple(attrs...), nil

	case pr.Map:
		entries := []rel.DictEntryTuple{}
		val.Map().Range(func(key pr.MapKey, value pr.Value) bool {
			val, err := walkThroughMessage(value)
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

	case pr.List:
		list := []rel.Value{}
		len := message.Len()
		for i := 0; i < len; i++ {
			item, err := walkThroughMessage(message.Get(i))
			if err != nil {
				list = append(list, rel.NewArrayItemTuple(i,
					rel.NewString([]rune(err.Error()))))
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
	case []byte:
		return rel.NewTuple(), nil
	case pr.Enum, pr.EnumNumber:
		return rel.NewTuple(), nil
	default:
		return rel.NewTuple(rel.NewStringAttr("@error",
			[]rune(fmt.Errorf("%T is not supported data type", message).Error()))), nil
	}
}
