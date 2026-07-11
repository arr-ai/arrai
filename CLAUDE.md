# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is arr.ai?

Arr.ai is a functional data transformation and query language built on set/relational algebra. It includes a CLI (`arrai`), interactive shell, server mode, and a bundler. Written in Go.

## Build & Development Commands

```bash
make generate        # Regenerate parser + bundle stdlib (run after grammar or stdlib changes)
make build           # Build the arrai binary
make install         # Build, install, and create ai/ax symlinks
make test            # Lint + test with -race, build tags, and arrai test suite
make lint            # Run golangci-lint (v1.54.0, strict config)
make tidy            # gofmt + goimports + go mod tidy
make check-clean     # Verify no uncommitted generated file changes (CI gate)
```

### Running tests

```bash
go test ./...                          # All Go tests
go test -race -tags timingsensitive ./...  # Full test suite as CI runs it
go test -run TestFoo ./rel             # Single test in a package
go test -v -run TestFoo ./syntax       # Verbose single test
./arrai test                           # Run all *_test.arrai files
./arrai test path/to/dir               # Run arrai tests in a specific directory
```

### Code generation

The parser is auto-generated from `syntax/arrai.wbnf` via `tools/parser/generate_parser.go`. The standard library is bundled from `syntax/stdlib/` into `syntax/embed/*.arraiz`. Always run `make generate` after modifying the grammar or stdlib, and ensure `make check-clean` passes before pushing.

## Architecture

### Evaluation pipeline

```
Source text → Parser (WBNF grammar) → AST → Compiler (syntax/compile.go) → Expr → Eval(ctx, Scope) → Value
```

### Key packages

- **`rel/`** — Core types: `Value` interface (Number, String, Set, Tuple, Dict, etc.), `Expr` interface, `Scope`, `Pattern`, set/tuple/relational operations. All values are immutable.
- **`syntax/`** — Parser, compiler (`compile.go` is the main file ~47KB), standard library (`std_*.go`), import system, bundler. Grammar defined in `arrai.wbnf`.
- **`cmd/arrai/`** — CLI entry point using urfave/cli. Commands: run, eval, shell, serve, test, bundle, transform, etc. Binary name detection (`ai` → shell, `ax` → transform).
- **`pkg/`** — Utilities: shell/REPL, test runner, bundle system, context-based filesystem and caching.
- **`engine/`** — Expression evaluation engine core.
- **`translate/`** — Format translators (protobuf, XML, YAML).

### Core interfaces (`rel/value.go`)

- **`Value`** — Immutable value type. Embeds `Expr` and `frozen.Key`. Methods: `Kind()`, `IsTrue()`, `Less()`, `Export()`.
- **`Expr`** — Evaluatable expression. Methods: `Eval(ctx, Scope) (Value, error)`, `Source()`.
- **`Pattern`** (`rel/pattern.go`) — Structural pattern matching with `Bind(ctx, Scope, Value)`.

### Immutable data structures

The codebase uses `github.com/arr-ai/frozen` extensively for immutable maps and sets. `Scope` is built on frozen maps. Set values wrap `frozen.Set`.

### Import system (`syntax/import.go`)

Supports local files, `//` standard library modules (`//str`, `//math`, `//seq`, etc.), and remote HTTP imports with caching. Module root resolved via `go.mod` sentinel file.

## Coding conventions

- Go 1.21, module path `github.com/arr-ai/arrai`
- Max line length: 120 chars
- Linting: golangci-lint with ~20 linters enabled (see `.golangci.yml`)
- Testing: `testify/assert` and `testify/require` for assertions
- CI runs on Ubuntu, macOS, and Windows

## Release mechanics

- Releases are cut by **goreleaser** triggered on `push: tags: 'v*.*.*'` — not by `gh release create`. Just push the tag from master and goreleaser creates the GitHub release and uploads binaries automatically. Config lives at `.github/goreleaser_configs/goreleaser.yml`.
- arrai is **not** in `marcelocantos/homebrew-tap` — don't try to install via `brew install marcelocantos/tap/arrai`. Users get binaries from the GitHub release page.
- Version is injected at build time via `-ldflags -X main.Version=v{{.Version}}`; no in-source version constant to bump.

## docs/ build uses npm, not yarn

- `docs/` uses **npm** with a committed `package-lock.json`. There is **no `yarn.lock`** — do not add one.
- Transitive dep pins live in the `overrides` field of `docs/package.json` (e.g., `minimatch`, `serialize-javascript`, `lodash`). Yarn v1 **ignores** the `overrides` field, so any CI step that runs `yarn install` will re-resolve from scratch and silently drop the pins, pulling newer transitives that currently break the build (webpackbar/webpack ProgressPlugin schema mismatch). Always use `npm ci && npm run build` in CI and locally.
- When fixing a transitive CVE, add a line to `overrides` and run `npm install` in `docs/` — do **not** run `npm update`, which pushes webpack forward past webpackbar compatibility.

## CI / deploy previews

- **GitHub Actions** (`.github/workflows/`): `go.yml`, `docs.yml`, `release.yml`, `docker.yml`, `wasm.yml`, `generate-tag.yml`. `docs.yml` runs only on `pull_request`.
- **Netlify** is connected as a GitHub App and posts deploy-preview checks (`deploy/netlify`, `Header rules`, `Pages changed`, `Redirect rules`) on every docs PR. Netlify runs its **own** build pipeline separate from GitHub Actions — fixing `docs.yml` does not fix Netlify. There is no `netlify.toml` in the repo; config is in the Netlify dashboard, which I usually don't have direct access to.
- `arr-ai/arrai`'s master branch protection has **empty** `required_status_checks.contexts` and `checks` — so Netlify failures are informational only and don't block merge. Don't get stuck waiting for them, and don't let `gh pr checks --watch --fail-fast` bail out early because of them — either drop `--fail-fast` or filter to the GitHub Actions checks only.
