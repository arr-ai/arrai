package pb

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
)

func TestFastWay(t *testing.T) {
	in, err := ioutil.ReadFile("../../examples/protocolbuffer/sysl/petshop")
	assert.Nil(t, err)
	module := &sysl.Module{}
	err = proto.Unmarshal(in, module)
	assert.Nil(t, err)

	apps := module.GetApps()
	assert.NotNil(t, apps["PetShopApi"])
}

func TestGenerateBinary(t *testing.T) {
	fd, err := protoregistry.GlobalFiles.FindFileByPath("sysl.proto")
	assert.Nil(t, err)

	msgDescriptors := fd.Messages()
	md := msgDescriptors.ByName("Module")
	message := dynamicpb.NewMessage(md)

	in, err := ioutil.ReadFile("../../examples/protocolbuffer/sysl/petshop")
	assert.Nil(t, err)

	err = proto.Unmarshal(in, message)
	assert.Nil(t, err)

	fieldApps := message.Descriptor().Fields().ByName("apps")
	valueApps := message.Get(fieldApps)

	fmt.Println(valueApps)
}
