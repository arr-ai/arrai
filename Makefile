include VersionReport.mk

all: lint test wasm

parser:
	go generate .

build: parser
	go build -ldflags=$(LDFLAGS) ./cmd/arrai

wasm: parser
	GOOS=js GOARCH=wasm go build -o /tmp/arrai.wasm ./cmd/arrai

install: parser
	go install -ldflags=$(LDFLAGS) ./cmd/arrai
	[ -f $$(dirname $$(which arrai))/ai ] || ln -s arrai $$(dirname $$(which arrai))/ai
	[ -f $$(dirname $$(which arrai))/ax ] || ln -s arrai $$(dirname $$(which arrai))/ax

lint: parser
	golangci-lint run

test: parser
	go test $(GOTESTFLAGS) -tags timingsensitive ./...
	GOARCH=386 go build ./...
docker:
	docker build . -t arrai
