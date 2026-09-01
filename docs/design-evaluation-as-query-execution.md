# Evaluation as Query Execution

*Design note, 2026-09. Written at the close of the 2026-08 evaluator
performance overhaul (see `perf-ledger-2026-08.html`), which took the
reconstruct workload from a reported 58.8–67.8s to 1.51s. That work mined
out the representation-level wins; this document is the plan for the next
order of magnitude, aimed at the regime that actually matters: workloads
with data sets orders of magnitude larger than our benchmarks.*

## Thesis

Arr.ai's surface syntax already *is* the relational algebra. There is no
SQL-to-algebra translation gap: the compiled expression tree is a logical
query plan that today gets executed exactly as written — every intermediate
materialized, every operator evaluated naively, every index built reactively.

The thesis of this design: **because the language is pure and its syntax is
the algebra, the interpreter and the query executor can be the same
machine.** Evaluation becomes plan execution. There is no embedded-database
subsystem; the evaluator *is* one.

Four seemingly independent improvement programs — streaming, algebraic
rewrites, query planning, and indexing — are four views of the single
artifact this thesis puts at the centre: **the plan**.

## The four forces, unified

1. **The expression tree is the logical plan.** Algebraic rewrites
   (predicate pushdown, projection pruning, semi-join reduction,
   count/exists passthrough) transform it. Purity makes every rewrite
   unconditionally sound with respect to *values*; the semantic-analysis
   layer below licenses the rewrites that also need structural facts.
   Rewrites are not an end in themselves: pushdown shrinks flows early and
   *lengthens fusable chains* — its output is raw material for streaming.

2. **The query plan is the decision layer**, and its decisions factor
   cleanly:
   - **Per node**: which algorithm and access path — scan vs index probe,
     hash-join build side, fuse into neighbour.
   - **Per edge**: the most consequential bit in the system —
     **stream or materialize**.

3. **Streaming is the default edge state.** Materialization stops being how
   evaluation works and becomes an explicit plan operator — a *pipeline
   breaker* — inserted only where something forces it:
   - a deduplication boundary without a distinctness certificate;
   - an ordering requirement;
   - a join's build side;
   - **fan-out**: a `let`-bound relation consumed more than once must be
     materialized or recomputed — a genuine plan decision, since purity
     makes recomputation legal.

   A structural gift: **sequence pipelines (`>>` chains over arrays and
   dicts) carry no deduplication semantics at all.** They fuse freely, with
   no analysis and no risk. A large fraction of real model-pipeline code is
   exactly this shape. Set pipelines fuse where certificates permit.

4. **Indices dissolve into materialization format.** An index is *keyed*
   materialization. "Materialize this edge" and "build the demanded index"
   are one fused decision: where a breaker must exist anyway, compile-time
   demand dictates the shape to freeze it in, and the index falls out of the
   same pass — a byproduct of the stream, not a second scan.

### Why streaming leads

For the large-data regime, streaming is the one force whose payoff grows
with data size *even when the program is already well-written*. Rewrites
fire only where the program is suboptimal relative to its access paths —
and users demonstrably hand-optimize (the reconstruct library's indexing
patch was manual query planning). Fusing an n-stage pipeline removes n−1
full-size intermediates regardless of authorial skill, converts per-operator
fork/join parallelism (with a materialization barrier per operator) into one
parallel pass over the whole chain, and defuses the super-linear
memory-pressure cliffs that only appear at scale.

## The three-state set model

The runtime shadow of the plan. A set value is in one of three states:

| State | Meaning | Already shipped as |
|---|---|---|
| **Persistent** | materialized, immutable, structurally shared | frozen-backed sets; the arena's committed prefixes |
| **Frontier** | exclusively extended in place between snapshots | `colStore` ownership (relations), `appendBuf` (strings, arrays), the `canonical` flag's fold path |
| **Suspended** | a plan: base + pending operators; streams under enumeration; forces at identity boundaries; answers `Count` from metadata without forcing | selection-vector views (`where` as base + row ids) — the degenerate case |

Transitions: operator application suspends; `Has`/`Hash128`/`Equal` force;
single-reference extension keeps the frontier. The 2026-08 work built the
first two states and one instance of the third without naming them; this
design generalizes the suspended state from "selection vector" to "operator
pipeline".

This is also the synthesis of two prior art threads: **frozen** contributes
the persistent state (cheap retention), **linqgo/linq** contributes the
suspended state (lazy operator edges with metadata passthrough — its
count-without-enumeration was the seed of the property layer below), and the
overhaul contributed the frontier state between them.

## The semantic-analysis layer (certificates)

Rewrites and fusion need structural facts. Codd bought analyzability by
restricting the operators; arr.ai bought expressiveness with `=>`, so
analyzability is **recovered by classification** — bodies sort into a
lattice of recognized classes, with "opaque" as the always-sound floor:

1. **Projection/rename** — a tuple literal of bare `.attr` references:
   `(:.x)`, `(y: .x)`. Codd's π+ρ, decidable from syntax. Count passes
   through iff the projected attrs contain a candidate key. Needs no numeric
   semantics; with class 3, carries most of the analyzer's weight.
2. **Injective transform** — *demoted to conditional*. Under any
   finite-significand floating representation, `x + 1 == x` across the
   entire upper half of the representable range (beyond 2⁵³ for float64,
   beyond 10³⁴ for decimal128): absorption is the defining property, not an
   edge case. Arithmetic injectivity therefore exists only:
   - by **provenance** — `@` indices of arrays/dicts/strings are exact
     integers bounded by container size, so `(@: .@ + 1, ...)` reindexing is
     injective because of where the value came from;
   - by **runtime certificate** — lazy per-column min/max stats (zone maps,
     built like the lazy indexes) proving *this input's* domain sits inside
     the representation's **exact-integer window** — the significand bound,
     not the representation max. The language never assumes; the evaluator
     checks or declines.
   - unconditionally, only if/when integers become exact (bigint default) —
     which is an additional argument for the numeric-tower migration:
     exactness buys back analyzability that floating semantics destroy.
3. **Key-carrying** — a candidate key passed through untouched among
   arbitrary computation: `(a: .a, b: f(.))` with `{a}` a key. Injective by
   consequence, and the key survives into the result — the door to
   functional-dependency propagation through the whole algebra.
4. **Opaque** — assume nothing; dedup as today.

**Where keys come from**, both directions:
- *Statically*: dicts, arrays and strings have `@` as a key **by
  construction**; group-by results are keyed by their grouping attrs; join
  keys compose; FDs propagate per the lattice.
- *Empirically*: the lazy group index is already a key detector — when
  `groupBy(p)` yields as many groups as rows, `p` is a verified candidate
  key, discovered for free as a byproduct of work the evaluator was doing
  anyway. Cache the bit; subsequent projections get count-passthrough and
  dedup elision with no static inference at all.

Payoffs beyond `count`: an injective-body `=>` needs **no deduplication on
materialization** (the entire builder/hash-index cost of the result
vanishes), and class membership is a **fusion-legality certificate** —
injective stages fuse without dedup barriers.

Existing machinery this extends rather than replaces: `matchEqAttrPredicate`
(body analysis in miniature), `TupleExpr.staticShape` (static schemas), the
memoised `groupBy` (runtime fallback and key detector).

## The hash-identity Rubicon

Decision baked into the 128-bit hash design from the start: **treat hash
equality as value equality.** The risk decomposition:

- *Coincidental* collision is a non-issue: at 2⁴⁰ live values, collision
  odds ≈ 2⁻⁴⁹ — UUID-class safety.
- The real threat is *algorithmic* convergence: structurally different
  values whose hash computations land on the same result by construction.
  We have a live specimen — `["x","x"]` vs `["y","y"]` collided because xor
  cancels structurally identical terms (fixed in #738 by switching arrays to
  `Mix`).

Crossing safely therefore requires an **injectivity audit of the hash
algebra**, written down per kind: Mix binds position; salts bind kind;
lengths bind arity; and the rule the #738 fix taught — *never xor terms that
can be identical for structurally different values*. The empirical
counterpart is 🎯T16's property layer running both directions
(`Equal ⇒ hash equal` and, fuzzed, `hash equal ⇒ Equal`).

The crossing itself is staged like everything else here: a `hashidentity`
build tag replaces deep `Equal` with 16-byte comparison, runs differentially
against the whole corpus (the `slowpath` discipline), gets measured, then
becomes the default. The payoff: every index-probe verification,
`equalValues` walk, and set-membership check collapses to one compare;
hash-consing (global value interning, pointer-fast equality, automatic
dedup of repetitive data) becomes trivial rather than heroic.

Interaction worth knowing before sequencing: cheap equality makes
deduplication cheap, which *lowers* the stakes of dedup-elision analysis —
the analyzer's remaining prizes are then count-passthrough, fusion legality
and index remapping. The two big ideas partially substitute for each other.

## Indices: compile-time demand, runtime choice

Purity makes dataflow SSA-clean, so **index-demand inference** is highly
decidable: joins name their keys syntactically, `where .a = e` names a probe
key, groupings name theirs. Annotating each relation-producing expression
with its future probe keys enables:

- **Indexes as byproducts of construction** — built during the fused
  materialization pass, not by a later rescan;
- **Storage-layout selection** — a relation only ever probed by `{a}` can
  *be* a keyed map (the index is the storage);
- **Index inheritance across lineage** — a live inefficiency today: group
  indexes are memoised per view, so `r where p` rebuilds indexes its parent
  already has, when a selection view's index is the parent's filtered. This
  is a near-term, measurable fix requiring no new analysis.

Limits, correctly assigned: higher-order flow needs interprocedural
summaries (or the conservative fallback); *physical choice* — which
candidate index wins on this data — belongs to runtime statistics, per the
zone-map philosophy. Laziness already performs loop-invariant index hoisting
(the memoised `groupBy` builds once regardless of loop structure) and
remains the sound fallback for everything the analysis can't see. Statically
known-unique indexes drop their bucket lists: half the memory, one-row
probes.

## The plan-invariance contract, and how this stays honest

Streaming and rewriting change *when* things evaluate. Errors are arr.ai's
only observable effect, so the semantic invariant for the entire programme
is:

> **Same value, same error, regardless of plan.**

The precedent is already shipped: the parallel operators' deterministic
first-element-error contract (no worker interrupted; lowest range's error
wins) is plan-invariance for one operator. It generalizes: any plan for an
expression must surface the error the naive left-to-right evaluation would.

The verification strategy comes free, because it is the house method scaled
up: **the current tree-walking evaluator remains in-tree as the reference
implementation**, and the plan engine is differentially tested against it
across the whole corpus — `slowpath` promoted from "fast path vs fallback"
to "planned execution vs naive execution". That oracle style caught a design
hazard on its first run in August; it is the reason this redesign can be
aggressive.

## The plan VM and serialization

The execution-tree serialization goal and the "traditional VMs feel wrong"
instinct resolve together. Arr.ai's hot loops live *inside* operators, not
between them; a stack/bytecode machine optimizes the layer that doesn't
matter. The right shape is a **dataflow/operator VM** — SQLite's VDBE is the
proof of concept that a relational VM *is* a serialization format: nodes are
physical operators with tight native inner loops; user closures compile to
kernels invoked per element; the serialized compiled form (`.arraiz` as
compiled plans) is the same artifact the engine executes.

Timing: the compile phase (parse + compile, ~60% of small-workload wall
clock, negligible at Arieh-scale data) is **explicitly out of scope** until
xbnf replaces wbnf; piecemeal work on the old foundation is waste. The plan
VM's serialization payoff lands naturally in that same window. 🎯T8 is set
aside on the same grounds.

## Out of scope for now, recorded

- **Compile times / wbnf** — parked pending xbnf (above).
- **Unboxed typed columns** — the data layer streaming eventually wants
  (shape inference → monomorphic column storage, real `int64`s instead of
  interface boxes); sequenced after the suspended-set engine exists.
- **Numeric tower** (bigint default, decimal floating point) — tracked as a
  semantics decision with the analyzability argument added; not a
  performance work item here.
- **Match VM** — measured unnecessary after the 🎯T15 bind-protocol work
  (dispatch is 6% of the pattern-stress benchmark); the flat-builder
  protocol is the substrate one would drive if pattern-heavy workloads ever
  justify it.

## Staging

Every stage lands with benchmarks, differential oracles, and byte-identical
goldens, per the overhaul's method. Dependency-ordered:

- **S0 — Scale benchmark suite** *(prerequisite for everything)*.
  Reconstruct at 10k–50k apps plus a millions-of-rows join/pipeline
  workload. All prior optimization was validated at small scale — constants,
  not asymptotics. Expect the profile ranking to change and accidental
  quadratics invisible at 3k rows to surface. Nothing below proceeds on
  small-data evidence alone.
- **S1 — Derived-view index inheritance.** Near-term, no new machinery,
  measurable now; compounds with everything later.
- **S2 — Suspended sets: sequence-pipeline fusion.** Generalize selection
  views to suspended operator pipelines for the semantics-free half (`>>`
  chains over arrays/dicts); metadata passthrough for `count`. First real
  streaming win, no certificate layer required.
- **S3 — Certificates and always-win rewrites.** Body classification
  (classes 1 and 3 first), FD/key propagation, empirical key caching;
  predicate pushdown and projection pruning as Expr-tree passes. Extends
  fusion to set pipelines.
- **S4 — Hash identity.** Algebra audit, `hashidentity` differential tag,
  measure, flip. Cheap, largely independent; can interleave earlier.
- **S5 — The plan proper.** Per-edge stream/materialize decisions, pipeline
  breakers, fan-out handling, index-as-materialization-format, zone-map
  statistics for physical choice.
- **S6 — Plan VM and serialization** *(post-xbnf)*. The operator VM as both
  execution engine and compiled `.arraiz` format.

Targets are filed per stage as work begins, with acceptance criteria naming
their oracles; 🎯T16 (property tests for the value laws) runs alongside as
the standing counterweight.
