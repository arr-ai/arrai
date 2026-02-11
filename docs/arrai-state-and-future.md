# Arr.ai: Current State and Future Directions

## 1. What arr.ai is

Arr.ai is a functional data transformation language built on a single unifying idea: **everything is a set**. Numbers, booleans, strings, arrays, dictionaries, and relations are all sets of tuples under the hood. This isn't just a theoretical nicety — it means set operations (union, intersection, difference, join) work uniformly across all data types, and relational algebra is a first-class citizen rather than a library bolted on top.

The implementation is a Go-hosted AST-walking interpreter with an interactive shell, a gRPC/WebSocket server mode for reactive data distribution, a bundler for packaging scripts, and a standard library covering encoding (JSON, YAML, XML, CSV, protobuf, XLSX), string/sequence manipulation, regular expressions, networking, and filesystem access.

The language originated at ANZ Bank and is open-sourced under Apache 2.0. The core team comprises roughly 5–7 contributors, with Marcelo Cantos as lead maintainer.

## 2. Current state

### 2.1 Project health

| Metric | Value |
|---|---|
| Total commits | ~1,215 |
| Stars / forks | 21 / 16 |
| Open issues | 146 |
| Latest release | v0.321.0 (Oct 2024) |
| Go version | 1.21 (CI), 1.18 (go.mod) |
| Codebase size | ~38K lines Go (24.5K source, 13.4K test) |

Commit velocity peaked in 2021 (182 commits), dropped to 21 in 2022, 6 in 2023, and effectively zero on master in 2024–2025. The project is in **maintenance mode**: dependencies are kept current (dependabot PRs), occasional bug fixes land, and releases still ship, but no new features are being developed.

The 146 open issues include several from 2020–2021 with P0/P1 labels that were never resolved. Many are edge-case language bugs (comment breaks tuple inference, complement on empty set fails, string template whitespace issues).

### 2.2 Feature completeness

The language is **substantially complete** for its intended purpose:

**What works well:**
- Full relational algebra with 8 join variants, nest/unnest, reduce, aggregation
- Pattern matching with structural destructuring, fallback patterns, extra-element capture
- Comprehensive encoding/decoding (JSON, YAML, XML, CSV, protobuf, XLSX)
- String interpolation with format specifiers
- Recursive functions via `let rec` and fixed-point combinators
- Safe/unsafe stdlib split for sandboxed execution
- Bundle system for dependency-frozen distribution
- Server mode with streaming gRPC/WebSocket observers

**What's missing or incomplete:**
- No concurrency primitives (all evaluation is single-threaded)
- No type annotations or compile-time type checking
- Macro system marked NYI in the grammar
- Array slicing syntax (`[1:5:2]`) marked NYI
- No module versioning (piggybacks on Go modules indirectly)
- No user-defined types beyond tuple/set composition
- Error handling is fail-fast only — no try/catch, limited recovery via `?:` operator
- WASM target is a bare cross-compilation with no JavaScript bindings

### 2.3 Architecture assessment

**Strengths of the current design:**

1. **The set-theoretic foundation is genuinely powerful.** Representing strings as `{(@: 0, @char: 'h'), (@: 1, @char: 'e'), ...}` means you can join a string against a lookup table, intersect two strings to find common characters, or use relational operations on what other languages would treat as opaque primitives. This isn't academic — it makes data transformation pipelines composable in ways that other languages struggle with.

2. **Immutability is total and well-implemented.** Every value, scope, and intermediate result uses the `frozen` library's persistent data structures. Structural sharing via hash array mapped tries means modifications are O(log n) and previous versions are retained. This makes the evaluation model inherently safe for caching and (potentially) parallelism.

3. **Specialised representations behind uniform interfaces.** While everything is conceptually a set, the implementation dispatches to optimised representations: `StringCharTuple` for string elements, `ArrayItemTuple` for array items, `Relation` with columnar storage for typed tuple sets. The `SetBuilder` bucketing system routes values to the right representation transparently.

4. **Clean separation of safe/unsafe capabilities.** The stdlib split and context-based capability injection mean untrusted arrai code can be evaluated in a sandbox with no filesystem or network access, without any changes to the language or evaluator.

**Weaknesses of the current design:**

1. **Pure AST walking with no optimisation passes.** Every `Eval()` call traverses the expression tree. There's no bytecode, no compilation to an intermediate form, no dead code elimination, no common subexpression elimination, no operator fusion. A chain `a + b + c + d` creates three separate `BinExpr` nodes each making independent `Eval()` calls. For data-heavy workloads this is dominated by the frozen collection overhead, but for compute-heavy expressions it's unnecessarily slow.

2. **High allocation pressure.** Every `SetCall()` allocates a `SetBuilder` to collect results. Every `Tuple.With()` allocates a new `GenericTuple` wrapping a new `frozen.Map`. Every `Scope.With()` during pattern binding allocates a new `Scope`. In a typical evaluation, these compound: a map operation over 10,000 elements creates 10,000 SetBuilders, 10,000+ intermediate tuples, and 10,000 scope extensions. There's no object pooling, no arena allocation, and no reuse.

3. **The compiler is a monolith.** `syntax/compile.go` is 1,599 lines with 44 compilation methods in a single switch dispatch. Operator definitions are split across three files (grammar, lexer tokens, compiler dispatch maps) with no single source of truth. This makes the language hard to extend and hard to optimise.

4. **No tail call optimisation.** Arrow expressions and recursive functions grow the Go call stack linearly. Deep recursion simply overflows. The `//fn.fix` combinator provides Y-combinator semantics but no stack safety. This is a real limitation for a functional language where recursion is the primary iteration mechanism.

5. **Generic sets are 5–10x slower than typed relations.** The join benchmarks demonstrate this clearly. Generic `Set` operations go through interface dispatch, type assertions, and polymorphic bucketing on every element. `Relation` operations use columnar storage and schema awareness to skip all of that. But the compiler doesn't automatically promote sets to relations — the user has to use relation syntax explicitly.

## 3. Future directions

The following are potential directions, loosely ordered from most pragmatic to most ambitious.

### 3.1 Performance: low-hanging fruit

**Object pooling for SetBuilder and intermediate allocations.** The `SetBuilder` type is created and discarded on every function application. A `sync.Pool` for SetBuilder instances (and their internal bucket maps) would reduce GC pressure substantially. This is a targeted change confined to `rel/value_set_builder.go` with no API changes.

**Fast-path for single-result function application.** `SetCall()` currently creates a SetBuilder, adds elements, and calls `Finish()` even when the result is a single value (which is the common case for `->` arrows). Detecting the single-value case and returning directly would eliminate the builder allocation entirely for the most frequent operation in the evaluator.

**Automatic relation promotion.** When the compiler can determine that a set contains only tuples with a common schema (which is true for most `=>` and `>>` expressions), it could emit a `Relation` instead of a generic `Set`. This would route those operations through the 5–10x faster columnar path automatically.

### 3.2 Performance: medium-term

**Bytecode compilation.** The AST walker could be replaced with a stack-based bytecode interpreter for the expression core (arithmetic, comparisons, field access, pattern binding). This eliminates the per-node method dispatch overhead and enables instruction-level optimisations. The frozen collection operations would still dominate for data-heavy workloads, but compute-heavy expressions (scoring functions, complex predicates) would see 2–3x improvements.

**Expression-level optimisations.** Even without bytecode, the compiler could perform:
- Constant folding beyond top-level literals (propagating through let bindings)
- Dead code elimination (unused let bindings)
- Algebraic simplification (`a & a → a`, `a | {} → a`, `a +> () → a`)
- Join reordering (pushing selections closer to data sources)

These are standard relational query optimisation techniques and would benefit complex transformation pipelines directly.

**Lazy enumerators.** Currently, set operations like `Where()`, `Map()`, and chained arrows materialise intermediate results eagerly. Pull-based lazy enumerators would defer materialisation until elements are consumed, reducing memory pressure for pipelined operations. The `frozen` library's iterator interface would need to support this, or a wrapper layer could be introduced.

### 3.3 Language evolution

**Gradual typing.** Arr.ai's set-theoretic type system has a natural correspondence with structural typing. A tuple `(name: "alice", age: 30)` has type `(name: String, age: Number)`. A relation `{|name, age| ("alice", 30), ("bob", 25)}` has type `{|name: String, age: Number|}`. These types could be inferred at compile time without any syntax additions — the compiler already knows the schema of relation literals and the attribute types of tuple literals.

Adding optional type annotations would enable:
- Compile-time error detection for attribute access on wrong tuple types
- Schema validation for encoding/decoding operations
- Better error messages ("expected tuple with attribute 'name', got Number" instead of runtime panics)
- Documentation of function signatures in the stdlib

This could be introduced incrementally: infer types where possible, accept annotations where provided, and fall back to dynamic typing where inference fails. No existing code would break.

**Concurrency via deterministic parallelism.** The language's purity and immutability make it a natural fit for automatic parallelisation. Set operations like `=>` (map over set elements), `where` (filter), and joins are embarrassingly parallel — each element can be processed independently. A fork-join executor could parallelise these operations transparently, with the immutable value model guaranteeing determinism.

This doesn't require language-level concurrency primitives (which would conflict with the pure functional model). It's purely an evaluator optimisation: the semantics remain sequential, but the implementation distributes work across goroutines. The `context.Context`-based design already supports this — each parallel evaluation gets its own context with shared import caches.

**Tail call optimisation via trampolining.** Implementing a trampoline in the evaluator would make recursive functions stack-safe. The approach: detect tail-position calls in the compiler, emit a special `TailCallExpr` that returns a thunk instead of evaluating, and have the top-level evaluator loop on thunks until a value is produced. This is a well-understood technique (Scheme, Clojure) and would remove a real limitation of the language.

**Completing the macro system.** The grammar already has syntax for macros (`{:macro[rule]?:content:}`) but the implementation is marked NYI. Macros with custom grammars would enable domain-specific notation — SQL-like query syntax, template languages, configuration formats — parsed and type-checked at compile time. This is one of arr.ai's most distinctive planned features and would significantly expand its applicability.

### 3.4 Ecosystem

**IDE support.** The VS Code extension (`arr-ai/vscode-arrai`) is dormant. A language server (LSP) providing completion, hover documentation, go-to-definition, and inline diagnostics would make the language dramatically more accessible. The compiler already tracks source positions via `parser.Scanner` — the infrastructure for source mapping exists.

**WASM as a first-class target.** The current WASM build is a bare Go cross-compilation (~5MB binary, full Go runtime). A purpose-built WASM target could:
- Provide JavaScript bindings for browser-side data transformation
- Enable arrai expressions in web applications (e.g., client-side filtering/aggregation)
- Support IndexedDB-backed import caching for offline use
- Use Web Workers for parallel evaluation

This would position arrai as an embeddable query/transformation engine for web applications, which is a use case with very few competitors.

**Package registry.** The current import system resolves modules via Go module paths, which works but couples arrai's module ecosystem to Go's. A lightweight package registry (even just a convention around GitHub repositories with `arrai.mod` files) would enable:
- Versioned dependencies for arrai packages
- A community library ecosystem independent of Go
- Bundle-based distribution of libraries (precompiled, sandboxed)

### 3.5 Positioning

Arr.ai occupies an unusual niche: it's more powerful than jq/yq for data transformation, more principled than SQL for relational operations, and more focused than general-purpose functional languages like Haskell or Elixir. Its closest analogues are:

- **jq** — JSON transformation. Arr.ai handles more formats, has richer relational operations, and its set-theoretic foundation is more composable. But jq is ubiquitous and arr.ai is unknown.
- **Datalog/Datomic** — Logic programming over relations. Arr.ai's relational algebra is more procedural but also more general (it handles non-relational data natively). Datalog can express recursive queries more naturally.
- **Q/kdb+** — Array/vector processing for financial data. Q is faster for numeric workloads but has famously hostile syntax. Arr.ai is more readable and handles heterogeneous data better.
- **Malloy** — Relational analytics language. Newer, backed by Google, focused on BI/analytics. Arr.ai is more general-purpose but lacks Malloy's SQL integration.

The most promising positioning for arr.ai is as an **embeddable data transformation engine**: a language that can be sandboxed, bundled, and invoked from Go applications to evaluate user-defined transformation logic safely. The safe/unsafe stdlib split, the bundle system, and the immutable evaluation model are all designed for this use case. Strengthening this story — with better embedding APIs, a stable evaluation ABI, and documentation aimed at application developers rather than language enthusiasts — would give the project its clearest path to adoption.

## 4. Summary

Arr.ai is a well-designed language with a powerful theoretical foundation and solid engineering. Its core idea — relational algebra as the universal data model — is sound and delivers real expressive power. The implementation is correct, well-tested, and architecturally clean.

The project's challenge is not technical quality but adoption. With 21 GitHub stars and declining commit velocity, it risks becoming a well-kept secret that quietly archives. The performance characteristics are adequate for its current use cases but leave significant headroom. The language feature set is substantially complete but missing concurrency, type checking, and tail call safety.

The most impactful investments would be:
1. **Performance quick wins** (object pooling, fast-path function application, automatic relation promotion) to close the gap with specialised tools
2. **Embeddable engine story** (stable Go API, sandboxing documentation, WASM target) to reach developers who need transformation capabilities inside their applications
3. **IDE support** (LSP, VS Code extension revival) to reduce the barrier to learning the language

The theoretical foundation is strong enough to support all of these directions. The question is whether the investment will be made.
