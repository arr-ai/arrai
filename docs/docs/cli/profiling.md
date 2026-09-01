---
id: profiling
title: Profiling
---

`arrai` can write Go [pprof](https://pkg.go.dev/runtime/pprof) CPU and heap
profiles for any invocation, via two global flags:

| flag | description |
|-|-|
| `-cpuprofile=<path>` | Write a CPU profile covering the whole run to `<path>`. |
| `-memprofile=<path>` | Write a heap profile (taken after a GC, at the end of the run) to `<path>`. |

These flags are consumed by `arrai` itself, before the command line is
otherwise parsed, so they must come **before** the command name and must use
the `=` form (`-cpuprofile <path>`, with a space, is not supported):

```bash
$ arrai -cpuprofile=cpu.prof run script.arrai
$ arrai -memprofile=mem.prof eval '2 + 2'
```

They work with any binary personality (`arrai`, `ai`, `ax`) and with
long-running commands like `serve`: stopping the command normally, letting it
error out, or interrupting it with Ctrl-C (`SIGINT`)/`SIGTERM` all flush the
profile correctly.

View the resulting profile with:

```bash
$ go tool pprof cpu.prof
```
