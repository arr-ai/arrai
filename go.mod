module github.com/arr-ai/arrai

go 1.13

require (
	github.com/arr-ai/frozen v0.6.0
	github.com/arr-ai/hash v0.4.0
	github.com/arr-ai/proto v0.0.0-20180422074755-2ffbedebee50
	github.com/arr-ai/wbnf v0.0.0-20200110014938-ba95372c7523
	github.com/go-errors/errors v1.0.1
	github.com/gorilla/websocket v1.4.1
	github.com/mediocregopher/seq v0.1.1-0.20170116151952-4c22a2e6eca9
	github.com/pkg/errors v0.8.1
	github.com/rjeczalik/notify v0.9.2
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	github.com/tealeg/xlsx v1.0.5
	github.com/urfave/cli v1.22.2
	google.golang.org/grpc v1.26.0
)

replace github.com/arr-ai/wbnf => ../wbnf
