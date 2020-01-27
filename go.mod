module github.com/arr-ai/arrai

go 1.13

require (
	github.com/arr-ai/frozen v0.13.0
	github.com/arr-ai/hash v0.4.0
	github.com/arr-ai/proto v0.0.0-20180422074755-2ffbedebee50
	github.com/arr-ai/wbnf v0.2.0
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/go-errors/errors v1.0.1
	github.com/gorilla/websocket v1.4.1
	github.com/mediocregopher/seq v0.1.1-0.20170116151952-4c22a2e6eca9
	github.com/pkg/errors v0.9.1
	github.com/rjeczalik/notify v0.9.2
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	github.com/tealeg/xlsx v1.0.5
	github.com/urfave/cli/v2 v2.1.1
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	golang.org/x/sys v0.0.0-20200124204421-9fbb57f87de9 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200122232147-0452cf42e150 // indirect
	google.golang.org/grpc v1.26.0
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/arr-ai/wbnf => ../wbnf

replace github.com/arr-ai/frozen => ../frozen
