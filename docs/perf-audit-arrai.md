# Performance Audit: arrai Post-Generics Migration

## Overview

After migrating arrai to the generic frozen library (`Set[T]`, `Map[K,V]`, `Key[T]`),
join benchmarks show a 1.4–1.7x slowdown with 50–60% more heap allocations. GC
consumes 54–58% of total CPU time. This report covers arrai-side findings and
optimizations. See the companion frozen report for library-level root causes.

## Benchmark Baseline (Apple M4 Max, Go 1.25)

| Benchmark | Pre-generics | Post-generics | Regression |
|---|---|---|---|
| RelationSetJoin100 | 86 us / 2,498 allocs | 141 us / 3,695 allocs | +64% / +48% |
| RelationSetJoin1000 | 1.08 ms / 24.7K allocs | 1.84 ms / 37.8K allocs | +70% / +53% |
| RelationSetJoin10000 | 14.5 ms / 247K allocs | 24.6 ms / 375K allocs | +70% / +52% |
| GenericSetJoin100 | 2.78 ms / 93K allocs | 3.87 ms / 143K allocs | +39% / +54% |
| GenericSetJoin1000 | 29.1 ms / 928K allocs | 42.8 ms / 1,455K allocs | +47% / +57% |
| GenericSetJoin10000 | 307 ms / 9.3M allocs | 434 ms / 14.5M allocs | +42% / +55% |

## Micro-Benchmark Findings

### Tuple operations (backed by `frozen.Map[string, Value]`)

| Operation | ns/op | B/op | allocs/op | Notes |
|---|---|---|---|---|
| Tuple.Get | 114 | 144 | 6 | Should be 0 allocs — frozen closure overhead |
| Tuple.Equal (2 attrs) | 539 | 832 | 26 | ~13 allocs per attribute |
| Tuple.Equal (5 attrs) | 1,767 | 2,368 | 86 | Scales linearly |
| Tuple.Hash (2 attrs) | 87 | 160 | 4 | Should be 0 allocs |
| Tuple.Enumerator | 69 | 176 | 3 | Acceptable |
| Tuple.Names() (10 attrs) | 1,300 | — | 51 | Rebuilds Set[string] every call |
| Tuple.Project (per elem) | — | — | 25 | Creates new frozen.Map + GenericTuple |
| getBucket() (10 attrs) | 1,800 | — | 64 | Rebuilds Names + sorts every call |

### Names operations (backed by `frozen.Set[string]`)

| Operation | ns/op | B/op | allocs/op |
|---|---|---|---|
| Names.Has | 76 | 64 | 4 |
| Names.Equal | 85 | 40 | 2 |
| Names.Intersect | 236 | 224 | 12 |
| Names.Minus | 170 | 176 | 9 |
| Names.Create (3 names) | 118 | 152 | 6 |

### Set operations (backed by `frozen.Set[Value]`)

| Operation | ns/op | B/op | allocs/op |
|---|---|---|---|
| Set[Value].Has | 35 | 8 | 1 |
| Set[Value].With | 190 | 507 | 2 |
| Set[Value].Range (1000 elems) | 22,352 | 20,256 | 657 |
| Set[Value].Equal (1000 elems) | 17,009 | 3,496 | 74 |
| SetBuilder.Add (1000 elems) | 721,368 | 804,263 | 24,975 |

### Relation vs GenericSet

| Operation | Relation | GenericSet | Ratio |
|---|---|---|---|
| RelationAttrs (1000 elems) | 59 ns / 3 allocs | 324 us / 11,600 allocs | 5,500x / 3,867x |
| Join grouping (500 elems) | 849 us / 18K allocs | 12.4 ms / 418K allocs | 14.5x / 23x |
| Builder.Add (1000 elems) | 436 us / 14K allocs | 721 us / 25K allocs | 1.7x / 1.8x |
| Enumeration (100 elems) | 34 us / 1,474 allocs | 2.2 us / 71 allocs | GenericSet wins here |

Note: Relation enumeration is expensive (14.7 allocs/element) because it reconstructs
a full `GenericTuple` with a new `frozen.Map[string, Value]` for each element.

## Root Causes (arrai-side)

### 1. `Tuple.Names()` rebuilds `frozen.Set[string]` every call

`GenericTuple.Names()` creates a new `SetBuilder[string]`, iterates all map keys,
and finishes a new `Set[string]` — every single time. This is called:
- Per element in `RelationAttrs()` to verify name consistency
- Per element in `SetBuilder.Add()` via `getBucket()`
- During join setup via `Names.Intersect()`

Since tuples are immutable, the result should be cached on first access.

### 2. `SetBuilder.getBucket()` is unnecessarily expensive

For every `Add()` call, `getBucket()`:
1. Calls `v.Names()` — which rebuilds a frozen.Set (see above)
2. Calls `OrderedNames()` — which sorts the names into a `[]string`
3. Constructs a `hashableNamesSlice` string for map lookup

For homogeneous sets (all tuples have the same names), this repeats identical
work millions of times.

### 3. GenericJoin creates per-element HAMT structures

The generic join path (`ops_rel.go:274`) calls `Tuple.Project(common)` for every
element in both input sets. Each `Project` creates:
- A new `frozen.Map[string, Value]` (HAMT) containing the projected attributes
- A new `GenericTuple` wrapping it
- Then uses this tuple as an HAMT key (triggering expensive Hash + Equal)

The positional Relation join avoids all of this by using `[]int` index remapping.

### 4. Relation enumeration reconstructs full tuples

`Relation.Enumerator().Current()` calls `valuesToTuple()` which:
- Allocates a `[]Attr` slice
- Calls `NewTuple(attrs...)` which builds a new `frozen.Map[string, Value]`
- 14.7 allocations per element

This is needed because the external API expects `Value` (which means `Tuple`),
but the internal storage is positional `Values` (`[]Value`).

## Optimization Strategies

### Strategy A: Cache `Tuple.Names()` [Very Low Risk]

Add a cached `names Names` field and `namesOnce sync.Once` to `GenericTuple`
(the pattern already exists for `orderNamesOnce`). Since tuples are immutable,
this is completely safe.

**Expected impact**: Eliminates millions of `Set[string]` constructions per join.
`RelationAttrs` on GenericSet drops from 11,600 allocs to near zero.

### Strategy B: Optimize `getBucket()` for homogeneous sets [Low Risk]

Cache the bucket key after the first computation. For the common case where all
tuples share the same names (which is always true for relation-like sets), this
avoids redundant `Names()` + `OrderedNames()` + string construction.

**Expected impact**: ~64 fewer allocs per `SetBuilder.Add()` call for 10-attr tuples.

### Strategy C: Lightweight join grouping key [Medium Risk]

Instead of `Tuple.Project(common)` creating a full `GenericTuple` + HAMT per
element, compute a hash directly from the common attribute values and group using
a Go native `map[uintptr][]Value`. This avoids per-element HAMT construction.

**Expected impact**: GenericJoin grouping could approach Relation join performance
(14.5x improvement potential).

### Strategy D: Lazy tuple reconstruction in Relation.Enumerator [Medium Risk]

Return a lightweight "view" tuple backed by the positional data + attribute names
instead of constructing a full `GenericTuple`. The view would implement the `Value`
interface but avoid building a `frozen.Map` until actually needed (e.g., if
`Get()` is called).

**Expected impact**: Eliminate 14K+ allocs per 1000-element Relation enumeration.

## Phase 2 Results: arrai-Side Caching (28 Feb 2026)

Strategies A and B were implemented along with additional optimizations to
`SetBuilder.Add()`, `Relation.Enumerator()`, and `valuesToTuple()`. Strategies C
and D are deferred to a follow-up.

### Changes

**`rel/value_tuple.go`** — `GenericTuple` struct gains three cached fields:

- `cachedNames Names` + `cachedNamesOnce sync.Once` — `Names()` now computes
  the `frozen.Set[string]` once and returns the cached value on subsequent calls.
  (Strategy A)
- `cachedBucket fmt.Stringer` + `cachedBucketOnce sync.Once` — `getBucket()`
  caches the `hashableNamesSlice` result after first computation. (Strategy B)
- `getSetBuilder()` now calls `TupleOrderedNames(t)` (already cached via
  `orderNamesOnce`) instead of `t.Names().OrderedNames()`.

**`rel/value_set_builder.go`** — `SetBuilder.Add()` called `v.getBucket()` twice
when creating a new bucket (once for lookup, once for insertion). Fixed to call
once and reuse the result.

**`rel/value_set_rel.go`** — Three changes:

1. `Relation` struct gains an `attrMap map[string]int` field, eagerly computed in
   `newRelation()`. All six call sites that previously called `mapIndices(r.attrs,
   r.p)` per-invocation now use the cached `r.attrMap`.
2. `valuesToTuple()` builds a `frozen.Map[string, Value]` directly via
   `frozen.MapBuilder` instead of allocating a `[]Attr` slice and routing through
   `NewTuple()` (which dispatches to specialized tuple types — unnecessary here
   since Relation rows are always generic tuples).

### Benchmark Results (Apple M4 Max, Go 1.25, benchstat p=0.008)

#### Micro-benchmarks

| Benchmark | Before | After | Δ time | Δ allocs |
|---|---|---|---|---|
| getBucket (2 attrs) | 235 ns / 11 allocs | 1.9 ns / 0 allocs | **−99.2%** | **−100%** |
| getBucket (5 attrs) | 491 ns / 17 allocs | 1.9 ns / 0 allocs | **−99.6%** | **−100%** |
| getBucket (10 attrs) | 2,050 ns / 75 allocs | 1.9 ns / 0 allocs | **−99.9%** | **−100%** |
| RelationAttrs GenericSet/100 | 32.1 us / 1,171 allocs | 12.3 us / 268 allocs | **−61.7%** | **−77.1%** |
| RelationAttrs GenericSet/1000 | 351 us / 11,654 allocs | 129 us / 2,622 allocs | **−63.3%** | **−77.5%** |
| SetBuilder.Add/100 | 79.3 us / 3,128 allocs | 54.8 us / 2,007 allocs | **−30.9%** | **−35.8%** |
| SetBuilder.Add/1000 | 830 us / 30,970 allocs | 586 us / 19,981 allocs | **−29.4%** | **−35.5%** |
| Relation.Enumerate/100 | 40.0 us / 1,772 allocs | 22.4 us / 768 allocs | **−43.9%** | **−56.7%** |
| Relation.Enumerate/1000 | 435 us / 17,641 allocs | 239 us / 7,680 allocs | **−43.1%** | **−56.5%** |
| GenericJoinGroupBy/100 | 2.74 ms / 109K allocs | 2.51 ms / 91K allocs | **−8.1%** | **−16.1%** |
| GenericJoinGroupBy/500 | 14.4 ms / 509K allocs | 12.2 ms / 427K allocs | **−15.4%** | **−16.2%** |

#### End-to-end joins

| Benchmark | Before | After | Δ time | Δ allocs |
|---|---|---|---|---|
| GenericSetJoin100 | 4.59 ms / 172.7K allocs | 4.18 ms / 140.6K allocs | **−8.9%** | **−18.6%** |
| GenericSetJoin1000 | 50.6 ms / 1,744K allocs | 41.2 ms / 1,418K allocs | **−18.7%** | **−18.7%** |
| GenericSetJoin10000 | 516 ms / 17.3M allocs | 441 ms / 14.1M allocs | **−14.5%** | **−18.5%** |
| RelationSetJoin100 | 165.7 us / 4,620 allocs | 165.9 us / 4,660 allocs | ~ | ~ |
| RelationSetJoin1000 | 2.15 ms / 47.2K allocs | 2.13 ms / 47.3K allocs | ~ | ~ |
| RelationSetJoin10000 | 29.3 ms / 465K allocs | 29.2 ms / 467K allocs | ~ | ~ |

`RelationSetJoin` is unaffected as expected — it uses the positional join path
which was already optimized.

### Remaining gap vs pre-generics baseline

| Benchmark | Pre-generics | Current (post Phase 2) | Gap |
|---|---|---|---|
| GenericSetJoin100 | 2.78 ms / 93K allocs | 4.18 ms / 141K allocs | +50% / +51% |
| GenericSetJoin1000 | 29.1 ms / 928K allocs | 41.2 ms / 1,418K allocs | +42% / +53% |
| RelationSetJoin100 | 86 us / 2,498 allocs | 166 us / 4,660 allocs | +93% / +87% |
| RelationSetJoin1000 | 1.08 ms / 24.7K allocs | 2.13 ms / 47.3K allocs | +97% / +91% |

The remaining regression is dominated by frozen library closure/boxing overhead
(Phase 1) and HAMT node layout (Phase 3). The arrai-side strategies C and D
offer further gains on the GenericSet path but won't close the Relation gap,
which is entirely frozen-side.

## Benchmark Files

- `rel/micro_bench_test.go` — micro-benchmarks for individual operations
- `rel/perf_bench_test.go` — targeted benchmarks for join phases and hotspots
- `rel/join_bench_test.go` — existing end-to-end join benchmarks
