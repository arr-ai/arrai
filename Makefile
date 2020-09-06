include VersionReport.mk

.PHONY: all
all: lint test wasm

.PHONY: parser
parser: syntax/parser.go

syntax/parser.go: tools/parser/generate_parser.go syntax/arrai.wbnf
	go run $^ $@

build: parser
	go build -ldflags=$(LDFLAGS) ./cmd/arrai

wasm: parser
	GOOS=js GOARCH=wasm go build -o /tmp/arrai.wasm ./cmd/arrai

install: parser
	go install -ldflags=$(LDFLAGS) ./cmd/arrai
	[ -f $$(dirname $$(which arrai))/ai ] || ln -s arrai $$(dirname $$(which arrai))/ai
	[ -f $$(dirname $$(which arrai))/ax ] || ln -s arrai $$(dirname $$(which arrai))/ax

tidy:
	go mod tidy
	gofmt -s -w .
	goimports -w .

lint: parser
	golangci-lint run

test: parser
	go test $(GOTESTFLAGS) -tags timingsensitive ./...
	GOARCH=386 go build ./...
	make build && ./arrai test

docker:
	docker build . -t arrai

static:
	go-bindata -o cmd/arrai/static.go internal/build/main.arrai
