---
id: interop
title: Interoperability
---

Arr.ai has compilers implemented in Go and JS (Wasm). The `arrai` command line tool invokes the Go compiler and evaluates the resulting expression tree. That process can also be invoked from other Go programs by importing the `syntax` package and calling:

```go
val, err := syntax.EvaluateExpr(arraictx.InitRunCtx(context.Background()), "", "<arr.ai source>")
```

The `arr.ai source` can be a string constant or loaded from a file. However loading sources dynamically can be complicated, since imports must still be resolved.

An alternative method is to use `arrai bundle` to produce an `arraiz` ZIP archive containing all necessary sources, compile that archive into a Go program (with e.g. [`embed`](https://pkg.go.dev/embed)), and then run it with:

```go
val, err := syntax.EvaluateBundle(arraictx.InitRunCtx(context.Background()), bundle, "", "arg", "...")
```

This will run the bundled script, with the following parameters set as `//os.args`, allowing the parent program to pass inputs.
