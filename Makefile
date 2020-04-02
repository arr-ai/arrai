all: test lint wasm

# TODO: If this Makefile is ever used for CI, suppress timingsensitive there.
test:
	go test $(GOTESTFLAGS) -tags timingsensitive ./...

lint:
	golangci-lint run

wasm:
	GOOS=js GOARCH=wasm go build -o /tmp/arrai.wasm ./cmd/arrai
