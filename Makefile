all: test lint wasm

test:
	go test ./...

lint:
	golangci-lint run

wasm:
	GOOS=js GOARCH=wasm go build -o /tmp/arrai.wasm ./cmd/arrai
