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
main: internal/arrai/echo.arraiz
	go build -o main main.go && chmod +x main

# Generates .arraiz files by calling arrai bundle on .arrai files.
%.arraiz: %.arrai
	arrai bundle $< > $@
