package pb

import (
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

//TransformProtoBufToTuple transforms the protocol buffer message to a tuple
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

	return rel.NewTuple(travel(message)...), nil
}

func travel(message *dynamicpb.Message) []rel.Attr {
	attrs := []rel.Attr{}

	message.Range(func(desc pr.FieldDescriptor, val pr.Value) bool {
		if desc.IsMap() {
			val.Map().Range(func(key pr.MapKey, value pr.Value) bool {
				attrs = append(attrs, rel.NewAttr(key.String(), rel.NewString([]rune(value.String()))))
				return true
			})
		}
		return true
	})

	return attrs
}
