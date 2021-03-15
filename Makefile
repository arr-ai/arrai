include VersionReport.mk

.PHONY: all
all: lint test wasm generate

.PHONY: parser
parser: syntax/parser.go

syntax/parser.go: tools/parser/generate_parser.go syntax/arrai.wbnf
	go run $^ $@

.PHONY: bindata
bindata: syntax/bindata.go

.PHONY: generate check-clean
generate: parser bindata
check-clean: generate
	git --no-pager diff HEAD && test -z "$$(git status --porcelain)"

syntax/bindata.go: syntax/implicit_import.arrai syntax/stdlib.arraiz
	go-bindata -mode 0644 -modtime 1 -pkg syntax -o syntax/bindata.go $^
	# remove all arraiz files after embedded file is generated.
	rm syntax/stdlib.arraiz
	gofmt -s -w $@

syntax/stdlib.arraiz: $(shell find syntax/stdlib -type f)
	go run ./cmd/arrai bundle -o syntax/stdlib.arraiz syntax/stdlib/stdlib.arrai

build: generate
	go build -ldflags=$(LDFLAGS) ./cmd/arrai

wasm: generate
	GOOS=js GOARCH=wasm go build -o /tmp/arrai.wasm ./cmd/arrai

install: generate
	go install -ldflags=$(LDFLAGS) ./cmd/arrai
	arrai info
	[ -f $$(dirname $$(which arrai))/ai ] || ln -s arrai $$(dirname $$(which arrai))/ai
	[ -f $$(dirname $$(which arrai))/ax ] || ln -s arrai $$(dirname $$(which arrai))/ax

tidy: generate
	go mod tidy
	gofmt -s -w .
	goimports -w .

lint: generate
	golangci-lint run

test: generate
	go test $(GOTESTFLAGS) -tags timingsensitive ./...
	[ "$$(go env GOOS)" == "darwin" ] || GOARCH=386 go build ./...
	make build && ./arrai test

docker: generate
	docker build . -t arrai
