package pb

import (
	"fmt"

	"github.com/arr-ai/arrai/rel"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
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

// convertProtoValToSyslVal walks through protobuf message and build a value whose type is rel.Value
func convertProtoValToSyslVal(val protoreflect.Value) (rel.Value, error) {
	switch message := val.Interface().(type) {
	case protoreflect.Message:
		attrs := []rel.Attr{}
		message.Range(func(desc protoreflect.FieldDescriptor, val protoreflect.Value) bool {
			item, err := convertProtoValToSyslVal(val)
			if err != nil {
				attrs = append(attrs, rel.NewAttr("@error_"+string(desc.Name()),
					rel.NewString([]rune(err.Error()))))
				return true
			}

			itemNum, isNum := item.(rel.Number)
			if desc.Enum() != nil && isNum {
				// item is protoreflect.EnumNumber
				num, success := itemNum.Int()
				if success {
					name := desc.Enum().
						Values().
						ByNumber(protoreflect.EnumNumber(num)).Name()
					attrs = append(attrs,
						rel.NewAttr(string(desc.Name()),
							rel.NewString([]rune(string(name)))))
				} else {
					attrs = append(attrs, rel.NewAttr(string(desc.Name()), item))
				}
				return true
			}

			attrs = append(attrs, rel.NewAttr(string(desc.Name()), item))
			return true
		})
		return rel.NewTuple(attrs...), nil

	case protoreflect.Map:
		entries := []rel.DictEntryTuple{}
		val.Map().Range(func(key protoreflect.MapKey, value protoreflect.Value) bool {
			val, err := convertProtoValToSyslVal(value)
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
			item, err := convertProtoValToSyslVal(message.Get(i))
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
