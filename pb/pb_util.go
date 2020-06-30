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

	attrs := []rel.Attr{}

	message.Range(func(descriptor pr.FieldDescriptor, value pr.Value) bool {
		attr := rel.NewAttr(string(descriptor.Name()), rel.NewString([]rune("Test")))
		attrs = append(attrs, attr)
		return true
	})

	return rel.NewTuple(attrs...), nil
}
