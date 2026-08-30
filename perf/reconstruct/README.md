# frozen generics slowdown — minimal reproduction

This reproduces a real-world performance regression we observed between
[arr-ai/frozen](https://github.com/arr-ai/frozen)'s pre-generics (v0.20.x)
and post-generics (v1.x) implementations, using `arr-ai/arrai` and the
public `anz-bank/sysl` arr.ai library.

On a real (proprietary, not included here) production workload, with a
performance fix **we wrote** for an unrelated inefficiency we found in
`reconstruct.arrai` along the way (details below — it's orthogonal to the
frozen regression itself, but needed to remove a confounding inefficiency
that otherwise masks it), comparing the two official arrai release tags
that bracket the frozen generics migration:

| Build | Real time |
|---|---|
| arrai v0.321.0 (frozen v0.20.3) | 62.6s |
| arrai v0.338.0 (frozen v1.12.0) | 465.3s |

That's a **~7.4x slowdown** for identical logic, identical output. (Without
the `reconstruct.arrai` patch, i.e. on the code as it ships today, we
measured a smaller but still real ~2.3x gap: ~115s vs ~265-285s — see "What
we measured" for why the patched comparison is the more representative
one.) This repro
demonstrates the same class of regression with **entirely synthetic data**
(a generated fake `.sysl` model — no real application/company data), built
from public repos:

- `vendor/sysl/**` — based on
  [anz-bank/sysl](https://github.com/anz-bank/sysl) commit
  `d34f35389ac019db7c37300ac62c2306bfa766b2`
  (`pkg/arrai/{sysl,reconstruct,util}.arrai`, `pkg/arrai/tools/appname.arrai`,
  `pkg/importer/utils.arrai`, `pkg/sysl/sysl.pb`). Two kinds of changes were
  made on top of that commit:
  - `//{...}` import paths were rewritten to plain local-relative paths, so
    this can run standalone without a working `go.mod`-based module fetch
    (arr.ai's `//{github.com/...}` import resolution always fetches
    `@latest` and does not honour `go.mod` version pins or `replace`
    directives — worth being aware of if you try to reproduce via a normal
    `arrai bundle` instead; see "Notes" below).
  - `pkg/arrai/reconstruct.arrai` has a performance fix **we wrote and
    applied ourselves** (not an upstream anz-bank/sysl change) — see "What
    we measured" for why. It indexes several relations that were being
    linearly re-scanned once per app/statement, which is unrelated to the
    frozen regression itself; it's included because it makes the
    frozen-version gap much cleaner to see (4.7x vs 1.5x, below). The
    unpatched file behaves identically, just slower on both frozen
    versions; if you want the untouched original for comparison, fetch
    `pkg/arrai/reconstruct.arrai` directly from the commit above.
- `vendor/arrai-contrib/util.arrai` — copied as-is from
  [arr-ai/arrai](https://github.com/arr-ai/arrai)'s `contrib/util.arrai`.
- `model.sysl` / `model.sysl.pb` — a generated model of N fake "apps", each
  with one endpoint calling the next app in a chain, a type with a few
  annotated fields, and a tag. Regenerate/rescale with `gen_model.py`
  (requires only the Python standard library).

## Reproducing

You need two `arrai` binaries, built from two different points in this same
repo: one from before the frozen generics migration, one from current
master (or any release since). We used the two release tags in the table
below, but any pre/post-migration pair should show the same pattern:

```sh
git checkout v0.321.0   # frozen v0.20.3, pre-migration
go build -o /tmp/arrai-old ./cmd/arrai

git checkout v0.338.0   # frozen v1.12.0, current
go build -o /tmp/arrai-new ./cmd/arrai

git checkout master     # back to your working branch
```

Then, from this directory:

```sh
cd vendor
time /tmp/arrai-old r --out=file:/tmp/out-old.arrai run.arrai ../model.sysl.pb
time /tmp/arrai-new r --out=file:/tmp/out-new.arrai run.arrai ../model.sysl.pb
cmp /tmp/out-old.arrai /tmp/out-new.arrai && echo "outputs are identical"
```

Rescale the model with `python3 gen_model.py 2500 > ../model.sysl` (then
recompile with `sysl pb --mode=pb ../model.sysl -o ../model.sysl.pb --root ..`,
which needs the [sysl CLI](https://github.com/anz-bank/sysl)) if you want to
see how the gap scales with data size.

## What we measured

Using the two official arrai release tags that bracket the frozen generics
migration — `v0.321.0` (frozen v0.20.3, pre-migration) and `v0.338.0`
(frozen v1.12.0, current) — against `gen_model.py 1000` (the model checked
in here), with the patched `reconstruct.arrai` that's actually committed in
`vendor/`:

| Build | Real time |
|---|---|
| v0.321.0 (frozen v0.20.3) | 10.5s |
| v0.338.0 (frozen v1.12.0) | 49.7s |

**~4.73x slower**, byte-identical output confirmed via `cmp`. This is the
number we consider representative, since the patch is what we're actually
running going forward.

For context, the same comparison against the *original*, unpatched
`reconstruct.arrai` (i.e. fetched straight from the anz-bank/sysl commit
above, no local changes) shows a smaller but still real gap — 46.4s vs
71.5s, ~1.54x. The original file has an unrelated O(apps × total_rows)
inefficiency: it repeatedly does `where`/`<--` scans of large, un-indexed
relations once per app/statement instead of indexing once up front. That
inefficiency dominates total time on the unpatched file and *masks* how
much slower frozen itself has gotten per-operation. Once it's patched out
and the workload is dominated by many cheap, already-indexed lookups
instead, the frozen-version gap becomes much starker (4.73x vs 1.54x) —
which points fairly specifically at increased per-operation cost in
frozen's generics implementation (allocation, hashing, tree traversal)
rather than at any algorithmic issue in the calling code.

## Notes for profiling

- `arrai r -cpuprofile=/tmp/cpu.prof ...` is *supposed* to produce a
  `pprof`-compatible CPU profile you can inspect with
  `go tool pprof -top -cum`, but as of v0.338.0 (and every release we
  checked, including v0.321.0) it doesn't actually work: `main.go`'s
  `prepareProfilers()` registers its `defer pprof.StopCPUProfile()` (and
  the profile file's `defer f.Close()`) inside itself, so they fire the
  instant `prepareProfilers()` returns — immediately after starting the
  profile, not at program exit. In practice this captures almost nothing
  (we saw ~200ms profiles for multi-minute runs). We pushed a one-line-ish
  fix for this on the
  [`fix/cpuprofile-defer-scope`](../../tree/fix/cpuprofile-defer-scope)
  branch — build `arrai` from there if you want profiling to actually work;
  it's a small, mechanical change and doesn't touch anything relevant to
  the frozen comparison itself, so it's independent of everything else in
  this directory.
- With that fix, in our own investigation, the dominant cost showed up as
  `frozen.Set.Where` / `internal/tree.(*branch).Where` (frozen-generics
  side) reached via `arr-ai/arrai/rel.Relation.Where` and
  `positionalRelation.Where` — i.e. filtering operations on `frozen.Set[T]`.
  We also saw substantial GC pressure (`runtime.gcDrain`/`scanObjectsSmall`)
  in early profiles, suggesting increased allocation rate is at least part
  of the story, not purely algorithmic complexity in `Where` itself.
