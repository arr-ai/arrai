---
id: install
title: Installation
---

`arrai` is a command-line tool written in [Go](https://golang.org). There are a variety of ways to install it, depending on your OS and use case.

## Summary

- Docker:
  ```
  docker run --rm -it -v $HOME:$HOME -w $(pwd) anzbank/arrai:latest
  ```
- Go:
  ```
  GO111MODULE=on go get -u github.com/arr-ai/arrai/cmd/arrai
  ```
- Source:
  ```
  git clone https://github.com/arr-ai/arrai.git
  cd arrai
  make install
  ```
- Binary: download from the [GitHub releases page](https://github.com/arr-ai/arrai/releases) to your `PATH`

Check the installation with `arrai help`.

## Requirements

- [Golang](https://golang.org/doc/install) version >= 1.13 (check with `go version`).
- [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports).

## Pre-compiled binary

1. Download the pre-compiled binaries matching your OS from the [releases page](https://github.com/arr-ai/arrai/releases).

1. Uncompress the archive and move the `arrai` binary to your desired location:

   1. On your `PATH` to run it with `arrai`
   1. Elsewhere to run it with `./arrai`, or some other `path/to/arrai`

## Go get it

First make sure you've installed Go:

```bash
go version
```

Fetch the `arrai` command's Go module:

```bash
GO111MODULE=on go get -u github.com/arr-ai/arrai/cmd/arrai
```

:::caution
Do NOT run this from inside a Go source directory that is module enabled, otherwise it gets added to go.mod/go.sum.
:::

## Docker

You can use `arrai` within a [Docker container](https://hub.docker.com/r/anzbank/arrai) (created from [this Dockerfile](https://github.com/arr-ai/arrai/blob/master/Dockerfile)):

```bash
docker run --rm -it -v $HOME:$HOME -w $(pwd) anzbank/arrai:latest
```

For example:

```
docker run --rm \
  -v $PWD:/go/src/github.com/arr-ai/arrai \
  -w /go/src/github.com/arr-ai/arrai \
  anzbank/arrai:latest run -v examples/simple/simple.arrai
```

Mac and Linux users can create an `alias` for the `arrai` command:

```
alias arrai="docker run --rm -it -v $HOME:$HOME -w $(pwd) anzbank/arrai:latest"
```

`arrai` can then be used from the same terminal window. Alternatively, add the `alias` to your `.bashrc` or `.zshrc` file to keep it permanently.

## Compile from source

```
git clone https://github.com/arr-ai/arrai.git
cd arrai
GOPATH=$(go env GOPATH) make install
```

## Try it out

If the installation worked, you should be able to run:

```bash
arrai
usage: arrai [<flags>] <command> [<args> ...]
...
```

You can always check your setup of `arrai` with:

```bash
arrai --version
arrai info
```

## VS Code Extension

Arr.ai has a VS Code extension which provides syntax highlighting for `.arrai` files. [Get it from here](https://marketplace.visualstudio.com/items?itemName=arr-ai.vscode-arrai), or search Extensions for "arrai".

## Shortcuts

Installing via `make install` will set up symlinks:

- `ai` => `arrai i`: interactive shell
- `ax` => `arrai x`: transform

If your installation method did not set up similar shortcuts, you may like to do so yourself (e.g. with `alias ai="arrai i"` in Bash).
