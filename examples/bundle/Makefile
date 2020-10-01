.PHONY: all build tidy run

all: build

build: main

run: main
	./main hello world

tidy:
	go mod tidy
	gofmt -s -w .
	goimports -w .

# Builds the main binary, including the generated bindata.go.
main: internal/arrai/bindata.go
	go build -o main main.go && chmod +x main

# Generates bindata.go, containing a []byte constant for internal/arrai/echo.arraiz.
internal/arrai/bindata.go: internal/arrai/echo.arraiz
	go-bindata -pkg arrai -o $@ $^

# Generates .arraiz files by calling arrai bundle on .arrai files.
%.arraiz: %.arrai
	arrai bundle $< > $@
