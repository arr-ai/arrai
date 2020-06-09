all: parser test lint wasm

install: parser
	go install ./cmd/arrai
	[ -f $$(dirname $$(which arrai))/ai ] || ln -s arrai $$(dirname $$(which arrai))/ai
	[ -f $$(dirname $$(which arrai))/ax ] || ln -s arrai $$(dirname $$(which arrai))/ax

test:
	go test $(GOTESTFLAGS) -tags timingsensitive ./...
	GOARCH=386 go build ./...
	go mod tidy

lint:
	golangci-lint run

wasm:
	GOOS=js GOARCH=wasm go build -o /tmp/arrai.wasm ./cmd/arrai

parser:
	go generate .

build:
	go build -ldflags=$(LDFLAGS) ./cmd/arrai

########
Version=$(shell git describe --tags)
FullCommit=$(shell git log --pretty=format:"%H" -1)
GoVersion=$(shell go version)
BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BuildOS=$(shell echo $(join $(join $(shell uname -s) , /), $(shell uname -p)) | tr A-Z a-z)

LDFLAGS='-X "main.Version=$(Version)" -X "main.GitFullCommit=$(FullCommit)" -X "main.BuildDate=$(BuildDate)" -X "main.GoVersion=$(GoVersion)" -X "main.BuildOS=$(BuildOS)"'
