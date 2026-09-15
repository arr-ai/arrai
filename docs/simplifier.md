# The simplifier (🎯T28)

`rel.Simplify` is a rewrite layer over the compiled expression tree. It runs
at the end of `syntax.Compile`, so every path that compiles source (run,
eval, test, serve, bundle) gets the same tree, and it runs before
`WriteCompiledPlan` encodes `plan.bin`, so a bundle carries the simplified
tree. It knows nothing about value representation: it rewrites `Expr`
nodes and rebuilds them through their constructors, which is what keeps it
independent of the arena, row-cursor and index work in `rel`.

```
Source → Parser → AST → Compiler → Expr ─Simplify→ Expr → Eval / plan.bin
```

## Walking the tree

`rel/expr_children.go` is the one place that knows every node's children.
For each child it reports how often the child evaluates relative to its
parent and which identifiers the parent binds for it:

| Kind         | Meaning                                       | Examples                                         |
|--------------|-----------------------------------------------|--------------------------------------------------|
| `ChildOnce`  | exactly once, unconditionally                 | operands of `+`, tuple attributes, `->` bodies   |
| `ChildMaybe` | at most once                                  | `if` branches, `cond` arms, `&&`/`\|\|` rhs      |
| `ChildMany`  | any number of times                           | lambda bodies, `=>` / `where` / `>>` bodies      |

Kinds join downward: a once-child of a maybe-child is maybe. The walker
returns a rebuild function alongside the children, and that function goes
through the node's constructor wherever the constructor analyses its
operands (`where`, `orderby`, `order`, `rank`, `call`, `+>`, `=>`, tuple
and set literals), so a rewritten node carries the same derived state as a
freshly compiled one. Expression types outside `rel` (`//pkg` references,
imports, `$"…"` strings) take part through `rel.Rewritable`; a node type
the walker does not know is opaque, and every let that contains one is
left alone.

## Occurrence analysis

For `let x = rhs; body` (compiled as an `ArrowExpr` applying `\x body`
once), the analysis walks `body` counting syntactic occurrences of `x`,
skipping subtrees whose binder shadows it. It records whether the sole
occurrence sits in a once-evaluated position, and whether any binder on
the path to it binds a name that `rhs` mentions (which substitution would
capture). Expressions embedded in patterns (fallbacks, literal
sub-patterns, dictionary keys) count as many-evaluated uses.

Dynamic uses are exactly the ones under `ChildMany`: inside `=>`, `where`,
`>>`, `:>`, `orderby`, reductions, lambda bodies and `let rec` bodies.

## Rewrites

**One-shot let folding.** When `x` occurs exactly once, in a once position,
uncaptured, and `rhs` has no side effects, the let becomes
`body[x := rhs]`. The payoff is downstream: the use site now sees the
shape of `rhs`, so the construction-time rewrites (predicate pushdown,
projection pruning, column extracts, index-answered `where`) apply across
what used to be a binding boundary.

**Unused let dropping.** When `x` does not occur and `rhs` is total, the
let becomes `body`. Total means it cannot fail or act: a literal, a
lambda, an identifier bound by an enclosing binder, or a tuple, array or
set of total expressions.

Rewrites apply bottom-up, then repeat to a fixpoint (at most four passes).

## Evaluation timing

Folding moves `rhs` from before `body` to its use site. The observable
consequences, and how each is handled:

- **Side effects never move.** `rhs` is folded only if it mentions no
  stdlib package with effects (`//log`, `//os`, `//net`, `//eval`,
  `//test`, `//runtime`) and contains no call at all, since a callee may
  have effects. `//str`, `//seq`, `//rel` and the other pure packages are
  fine.
- **Which error is reported can change.** If `rhs` and an earlier part of
  `body` would both fail, the program used to report the `rhs` failure and
  now reports the other. Both runs still fail. This is the one timing
  change the simplifier makes, and it is deliberate.
- **Dropping is invisible.** A dropped `rhs` is one that could not have
  failed or acted.

## Standing oracles

- `rel/simplify_test.go` builds trees directly and checks each rule and
  each refusal (dynamic use, conditional use, capture, effects, recursion,
  non-total drop).
- `syntax/simplify_test.go` compiles source and counts calls to a native
  function bound in scope: a let used inside a mapped body is evaluated
  once before and after simplification, so folding never duplicates work.
- The slowpath build (`-tags slowpath`, run in CI as `Test (slowpath)`)
  disables `Simplify` along with every other fast path, so the whole
  suite runs both with and without the rewrites: a differential oracle
  over the simplifier as well as over the representation shortcuts.
- `perf/TestReconstruct` and `./arrai test` hold: the reconstruct scenario
  matches its v0.321.0 reference and the corpus passes.
