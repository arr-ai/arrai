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

// FromProtoValue walks through protobuf message and build an equivalent rel.Value.
func FromProtoValue(val protoreflect.Value) (rel.Value, error) {
	switch message := val.Interface().(type) {
	case protoreflect.Message:
		var attrs []rel.Attr
		var err error
		message.Range(func(desc protoreflect.FieldDescriptor, val protoreflect.Value) bool {
			var item rel.Value
			item, err = FromProtoValue(val)
			if err != nil {
				return false
			}

			fieldName := string(desc.Name())
			itemNum, isNum := item.(rel.Number)
			if desc.Enum() != nil && isNum {
				// item is protoreflect.EnumNumber
				if num, success := itemNum.Int(); success {
					en := protoreflect.EnumNumber(num)
					name := desc.Enum().Values().ByNumber(en).Name()
					attrs = append(attrs, rel.NewStringAttr(fieldName, []rune(name)))
				} else {
					attrs = append(attrs, rel.NewAttr(fieldName, item))
				}
				return true
			}

			attrs = append(attrs, rel.NewAttr(fieldName, item))
			return true
		})
		if err != nil {
			return nil, err
		}
		return rel.NewTuple(attrs...), nil

	case protoreflect.Map:
		entries := make([]rel.DictEntryTuple, 0, val.Map().Len())
		var err error
		val.Map().Range(func(key protoreflect.MapKey, value protoreflect.Value) bool {
			var val rel.Value
			val, err = FromProtoValue(value)
			if err != nil {
				return false
			}

			at := rel.NewString([]rune(key.String()))
			entries = append(entries, rel.NewDictEntryTuple(at, val))
			return true
		})
		if err != nil {
			return nil, err
		}
		return rel.NewDict(false, entries...)

	case protoreflect.List:
		list := make([]rel.Value, 0, message.Len())
		for i := 0; i < message.Len(); i++ {
			item, err := FromProtoValue(message.Get(i))
			if err != nil {
				return nil, err
			}
			list = append(list, item)
		}
		return rel.NewArray(list...), nil

	case string, bool, int, int32, int64, uint, uint32, uint64, float32, float64:
		return rel.NewValue(message)

	case protoreflect.EnumNumber:
		// type EnumNumber int32
		return rel.NewValue(int32(message))
	default:
		return nil, fmt.Errorf("%v (%[1]T) not supported", message)
	}
}
