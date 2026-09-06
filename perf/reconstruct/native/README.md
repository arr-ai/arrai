# Handcrafted native reconstruct

A single-file C++ implementation of the reconstruct scenario — protobuf in,
Sysl source out — that produces byte-identical output to
`arrai run ../vendor/run.arrai ../model.sysl.pb`. It exists to put a
native-speed floor under the evaluator's number for this workload, so the
performance ledger can say not just "faster than before" but "this far from
the machine".

It is a faithful port of the pipeline in
`../vendor/sysl/pkg/arrai/{sysl,reconstruct}.arrai`: the same normalisation,
ordering, quoting and template rules, implemented generally over the protobuf
features that pipeline renders. It does **not** exploit the synthetic model's
regularity; it decodes the wire format (no libprotobuf) and renders whatever
apps, endpoints, types, fields, enums, annotations and tags it finds, with
the original's quirks preserved (endpoint params and tags, events, aliases
and views are unrendered there, so they are unrendered here).

```sh
clang++ -std=c++20 -O2 -o reconstruct reconstruct.cpp
./reconstruct ../model.sysl.pb | cmp - ../expected.arrai
```

`go test ./perf -run TestNativeReconstruct` compiles and byte-checks it
against the same golden the arr.ai pipeline is held to (skipped when no C++
compiler is on the PATH).

Measured 2026-08-31 (M4 Max, hyperfine, 1 warmup + 5 runs):

| Implementation | Time |
|---|---|
| `arrai run` (evaluator, post-optimisation) | 1.90 s |
| native (this) | **6.3 ms** |

The gap (~300×) is the honest measure of remaining interpretation overhead:
the native version does the same protobuf decode, the same normalisation
decisions, the same string assembly — without an evaluator, immutable
relational values, or a grammar-parsed source program. (The arrai figure
also includes parsing run.arrai and the sysl library itself, which the
native binary has compiled in; that is part of what "no evaluator" buys.)
