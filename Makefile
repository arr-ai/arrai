all: parser test lint wasm

install: parser
	go install -ldflags=$(LDFLAGS) ./cmd/arrai
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


######## Build information report
originalVersion=$(shell git describe --tags)

ifeq (-, $(findstring -, $(originalVersion))) #it is in branch
tagName= $(firstword $(subst -,  ,$(originalVersion)))
diffLogs = $(foreach item, $(shell git log --pretty=format:"%h" $(tagName)..HEAD),$(item))

ifeq (true, $(words $(diffLogs)) = 0 || $(shell git status -s) = ) # no changes
Version=$(tagName)
else

ifeq (, $(shell git status -s)) # it doesn't have uncommitted or untracked changes
Version=$(FullCommit) (~$(words $(diffLogs)) = $(tagName))
else
Version=DIRTY-$(FullCommit) (~$(words $(diffLogs)) = $(tagName))
endif

endif

else

# it is in tag
ifeq (,$(shell git status -s)) # no changes
Version=$(originalVersion)
else
Version=DIRTY-$(originalVersion)
endif

endif

FullCommit=$(shell git log --pretty=format:"%H" -1)
GoVersion=$(strip $(subst  go version, ,$(shell go version)))
BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BuildOS=$(shell echo $(join $(join $(shell uname -s) , /), $(shell uname -p)) | tr A-Z a-z)

LDFLAGS='-X "main.Version=$(Version)" -X "main.GitFullCommit=$(FullCommit)" -X "main.BuildDate=$(BuildDate)" -X "main.GoVersion=$(GoVersion)" -X "main.BuildOS=$(BuildOS)"'
