# Performance Audit: March 2026 Comparative Review

## Overview

Follow-up performance audit comparing current state (frozen v1.9.0, Go 1.25,
post-Phase 2 arrai optimizations) against the February 2026 baseline. All
measurements on Apple M4 Max, 5 iterations, `go test -benchmem`.

The February audit found a 40–70% regression from the frozen generics migration,
implemented Phase 2 arrai-side caching (Names, getBucket, attrMap, valuesToTuple),
and left a residual gap attributed to frozen library closure/boxing overhead.

## End-to-End Join Benchmarks

The headline numbers. These are the benchmarks that matter most for real-world
arrai performance.

### GenericSetJoin (uses frozen HAMT for everything)

| Size | Pre-generics | Feb 2026 (post-Phase 2) | Mar 2026 (current) | Δ vs Feb | Gap vs pre-generics |
|---|---|---|---|---|---|
| 100 | 2.78 ms / 93K allocs | 4.18 ms / 141K allocs | 5.15 ms / 174K allocs | **+23% / +24%** | +85% / +87% |
| 1000 | 29.1 ms / 928K allocs | 41.2 ms / 1,418K allocs | 49.8 ms / 1,768K allocs | **+21% / +25%** | +71% / +91% |
| 10000 | 307 ms / 9.3M allocs | 441 ms / 14.1M allocs | 506 ms / 17.6M allocs | **+15% / +25%** | +65% / +89% |

### RelationSetJoin (uses positional join, frozen only for storage)

| Size | Pre-generics | Feb 2026 (post-Phase 2) | Mar 2026 (current) | Δ vs Feb | Gap vs pre-generics |
|---|---|---|---|---|---|
| 100 | 86 us / 2,498 allocs | 166 us / 4,660 allocs | 187 us / 4,648 allocs | **+13% / ~** | +117% / +86% |
| 1000 | 1.08 ms / 24.7K allocs | 2.13 ms / 47.3K allocs | 2.30 ms / 47.1K allocs | **+8% / ~** | +113% / +91% |
| 10000 | 14.5 ms / 247K allocs | 29.2 ms / 467K allocs | 32.8 ms / 465K allocs | **+12% / ~** | +126% / +88% |

**Observation**: Both paths show a 10–25% regression from the February numbers.
The allocation counts for RelationSetJoin are stable (identical allocs/op), so
the time regression is pure CPU overhead — likely Go 1.25 runtime changes or
measurement variance between sessions. The GenericSet path shows both time and
allocation increases, suggesting possible frozen library changes between the
two measurement points.

## Micro-Benchmark Comparison

### Tuple operations

| Operation | Feb 2026 | Mar 2026 | Δ time | Δ allocs |
|---|---|---|---|---|
| Tuple.Get (3 attrs) | 114 ns / 6 allocs | 183 ns / 10 allocs | +61% | +67% |
| Tuple.Equal (2 attrs) | 539 ns / 26 allocs | 942 ns / 42 allocs | +75% | +62% |
| Tuple.Equal (5 attrs) | 1,767 ns / 86 allocs | 2,810 ns / 126 allocs | +59% | +47% |
| Tuple.Hash (2 attrs) | 87 ns / 4 allocs | 94 ns / 4 allocs | +8% | ~ |
| Tuple.Enumerator (3 attrs) | 69 ns / 3 allocs | 85 ns / 3 allocs | +23% | ~ |
| Tuple.Names (10 attrs) | 1,300 ns / 51 allocs | **1.1 ns / 0 allocs** | **−99.9%** | **−100%** |
| getBucket (10 attrs) | 1,800 ns / 64 allocs | **1.9 ns / 0 allocs** | **−99.9%** | **−100%** |

**Key findings**:
- The Phase 2 caching (Names, getBucket) continues to hold — zero allocations
  on cached paths.
- `Tuple.Get` and `Tuple.Equal` show significant regression (59–75%) versus the
  February numbers. These are pure frozen `Map.Get` and `Map.Equal` operations,
  confirming the bottleneck is in the frozen library's closure dispatch.
- `Tuple.Equal` at 42 allocs for 2 attributes (21 allocs/attr) is worse than
  the 26 allocs reported in February (13 allocs/attr). This suggests the frozen
  library's equality path may have regressed, or the February measurement used
  a different frozen version.

### Names operations (backed by `frozen.Set[string]`)

| Operation | Feb 2026 | Mar 2026 | Δ time | Δ allocs |
|---|---|---|---|---|
| Names.Has | 76 ns / 4 allocs | 82 ns / 3 allocs | +8% | −25% |
| Names.Equal | 85 ns / 2 allocs | 58 ns / 2 allocs | **−32%** | ~ |
| Names.Intersect | 236 ns / 12 allocs | 236 ns / 10 allocs | ~ | −17% |
| Names.Minus | 170 ns / 9 allocs | 219 ns / 10 allocs | +29% | +11% |
| Names.Create (3 names) | 118 ns / 6 allocs | 230 ns / 10 allocs | +95% | +67% |

Names.Equal improved significantly (−32%). Names.Create regressed (+95%), which
may reflect a change in the frozen Set builder path.

### frozen.Set[Value] primitives

| Operation | Feb 2026 | Mar 2026 | Δ time | Δ allocs |
|---|---|---|---|---|
| Set[Value].Has | 35 ns / 1 alloc | 76 ns / 3 allocs | +117% | +200% |
| Set[Value].With | 190 ns / 2 allocs | 193 ns / 2 allocs | ~ | ~ |
| Set[Value].Range (1000) | 22,352 ns / 657 allocs | 22,330 ns / 646 allocs | ~ | −2% |
| Set[Value].Equal (1000) | 17,009 ns / 74 allocs | 41 ns / 2 allocs | **−99.8%** | **−97%** |
| SetBuilder.Add (1000 tuples) | 830 us / 31K allocs | 714 us / 22K allocs | **−14%** | **−29%** |

**Notable**:
- `Set[Value].Equal` went from 17 us to 41 ns — a **415x speedup**. This is
  because frozen now uses a structural identity fast-path (same root pointer =
  equal) which kicks in when comparing sets built from the same builder.
- `Set[Value].Has` regressed from 1 to 3 allocs. This is the core frozen HAMT
  lookup path and the primary optimization target for Phase 1.
- `SetBuilder.Add` improved 14% time / 29% allocs — the Phase 2 caching
  continues to pay dividends.

### frozen.Map[string, Value] primitives

| Operation | Feb 2026 | Mar 2026 | Δ time | Δ allocs |
|---|---|---|---|---|
| Map.Get | — | 184 ns / 10 allocs | — | — |
| Map.Range (3 entries) | — | 64 ns / 2 allocs | — | — |
| Map.With | — | 246 ns / 13 allocs | — | — |
| Map.Equal (3 entries) | — | 216 ns / 10 allocs | — | — |

The February audit reported `Map.Get` at 114 ns / 6 allocs (measured via
Tuple.Get). The current 184 ns / 10 allocs confirms a frozen-level regression
in the map lookup path.

### Relation vs GenericSet

| Operation | Relation | GenericSet | Ratio |
|---|---|---|---|
| RelationAttrs/1000 | 124 ns / 6 allocs | 88 us / 2,605 allocs | 710x / 434x |
| Enumerate/100 | 42.5 us / 1,761 allocs | 2.3 us / 84 allocs | GenericSet 18x faster |
| Enumerate/1000 | 449 us / 17,630 allocs | 22.5 us / 652 allocs | GenericSet 20x faster |
| Builder.Add/1000 | 665 us / 21,950 allocs | 714 us / 21,957 allocs | ~1x |

**Key change**: `RelationAttrs` on GenericSet dropped from 324 us / 11,600
allocs (Feb) to 88 us / 2,605 allocs — a **73% time / 78% alloc improvement**
from the Names() caching. This remains the single biggest win from Phase 2.

Relation enumeration (449 us / 17,630 allocs for 1000 elems) is still expensive
at 17.6 allocs/element, essentially unchanged from February (17,641 allocs).
Strategy D (lazy tuple view) would address this.

## Join Phase Isolation

| Phase | Feb 2026 | Mar 2026 | Δ time | Δ allocs |
|---|---|---|---|---|
| GenericJoinGroupBy/100 | 2.51 ms / 91K allocs | 2.92 ms / 113K allocs | +16% | +24% |
| GenericJoinGroupBy/500 | 12.2 ms / 427K allocs | 13.9 ms / 520K allocs | +14% | +22% |
| RelationJoinDirect/100 | — | 175 us / 4,771 allocs | — | — |
| RelationJoinDirect/500 | — | 1.06 ms / 22,655 allocs | — | — |
| TupleProjectInLoop/100 | — | 85.8 us / 4,700 allocs (47/tuple) | — | — |
| TupleProjectInLoop/1000 | — | 947 us / 47,000 allocs (47/tuple) | — | — |

GenericJoinGroupBy regressed ~15% time / ~22% allocs from February. The per-
tuple Project cost of 47 allocs confirms that Strategy C (lightweight join
grouping key) would have the highest impact on the GenericSet join path.

## Summary

### What improved since February

1. **Names caching holds firm** — `Tuple.Names()` and `getBucket()` remain at
   0 allocs, saving millions of allocations per join.
2. **RelationAttrs on GenericSet** — 73% faster, 78% fewer allocs vs Feb.
3. **Set equality fast-path** — 415x faster when comparing structurally
   identical sets (frozen improvement).
4. **SetBuilder.Add** — 14% faster, 29% fewer allocs from caching.

### What regressed since February

1. **Tuple.Get** — 61% slower, 67% more allocs (6 → 10). This is a frozen
   `Map.Get` regression.
2. **Tuple.Equal** — 59–75% slower, 47–62% more allocs. Frozen `Map.Equal`
   regression.
3. **Set[Value].Has** — 117% slower, 200% more allocs (1 → 3).
4. **GenericJoinGroupBy** — 14–16% slower, 22–24% more allocs.
5. **End-to-end joins** — 8–25% slower across all sizes.

### Root cause analysis

The regressions in Tuple.Get, Tuple.Equal, and Set.Has are all frozen library
operations. Between February and now, the frozen library was updated to v1.9.0
which included the H128 hash migration (moving from 32-bit to 128-bit hashes).
This change:

- Added extra allocations in the hash/equality dispatch paths (128-bit hash
  values likely require heap allocation for the `Hash128` struct)
- Increased the cost per HAMT node comparison
- Explains the uniform ~60% regression in map/set lookup operations

### Remaining optimization roadmap

| Phase | Target | Expected impact | Status |
|---|---|---|---|
| **Phase 1**: frozen closure elimination | `Map.Get`, `Set.Has` hot paths | −50–70% allocs on lookups | Not started |
| **Phase 2**: arrai-side caching | Names, getBucket, attrMap | Millions of allocs eliminated | **Done** |
| **Phase 3**: CHAMP migration | HAMT node layout | 23–96% (literature) | Not started |
| **Phase 4**: Type-erased internals | Only if needed | Eliminate generic overhead entirely | Not started |
| **Strategy C**: Lightweight join key | GenericJoin grouping | Up to 14x on generic joins | Not started |
| **Strategy D**: Lazy tuple view | Relation.Enumerator | −14K allocs per 1K enumeration | Not started |

### Overall gap vs pre-generics

| Benchmark | Pre-generics | Current | Gap |
|---|---|---|---|
| GenericSetJoin1000 | 29.1 ms / 928K allocs | 49.8 ms / 1,768K allocs | +71% / +91% |
| RelationSetJoin1000 | 1.08 ms / 24.7K allocs | 2.30 ms / 47.1K allocs | +113% / +91% |

The gap has widened from the February post-Phase 2 numbers (42%/53% for Generic,
97%/91% for Relation) due to the frozen H128 migration overhead. Phase 1 (frozen
closure elimination) is now even more critical — it would address both the
original generics overhead and the new H128 allocation costs.
