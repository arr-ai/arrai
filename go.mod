module github.com/arr-ai/arrai

go 1.13

replace github.com/spf13/afero v1.3.5 => github.com/anz-bank/afero v1.2.4
replace github.com/joshcarp/gop => ../gop

require (
	cloud.google.com/go/storage v1.12.0 // indirect
	github.com/anz-bank/pkg v0.0.22
	github.com/arr-ai/frozen v0.15.0
	github.com/arr-ai/hash v0.5.0
	github.com/arr-ai/proto v0.0.0-20180422074755-2ffbedebee50
	github.com/arr-ai/wbnf v0.28.0
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/go-errors/errors v1.1.1
	github.com/gorilla/websocket v1.4.2
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/joshcarp/gop v0.0.0-20200924011502-d6f2efad81c9
	github.com/mattn/go-isatty v0.0.12
	github.com/pkg/errors v0.9.1
	github.com/rjeczalik/notify v0.9.2
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/afero v1.4.0
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli/v2 v2.2.0
	golang.org/x/net v0.0.0-20200923182212-328152dc79b1 // indirect
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43 // indirect
	golang.org/x/sys v0.0.0-20200923182605-d9f96fdee20d // indirect
	golang.org/x/tools v0.0.0-20200923182640-463111b69878 // indirect
	google.golang.org/genproto v0.0.0-20200923140941-5646d36feee1 // indirect
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.3.0
)
