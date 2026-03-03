# Arr.ai Agent Guide

## What is arr.ai?

Arr.ai is a functional data representation, transformation, and query language
built on set/relational algebra. Every arr.ai program is a single expression
that evaluates to a single value. There are no statements, no assignments, and
no side effects.

The primary CLI binary is `arrai`. Two shortcut symlinks are installed alongside
it:

- `ai` — opens the interactive shell directly
- `ax` — runs the transform command

### Key characteristics

- All values are immutable.
- The entire type system consists of only three kinds of values: **numbers**,
  **tuples**, and **sets**. Everything else (strings, arrays, booleans, dicts,
  functions, relations) is syntactic sugar layered over these three.
- `true` is `{()}` (the set containing the empty tuple). `false` is `{}`.
- Numbers are 64-bit binary floats.
- Programs are pure (no side effects). File I/O and network access are
  available only through explicit stdlib calls (`//os`, `//net`).

## Installation

### From source

```bash
git clone https://github.com/arr-ai/arrai.git
cd arrai
make install    # builds binary, installs to GOPATH/bin, creates ai/ax symlinks
```

Requires Go 1.24 or later.

### From releases

Download the relevant archive from the
[Releases page](https://github.com/arr-ai/arrai/releases).

## CLI Reference

### Global flags

| Flag | Alias | Description |
|---|---|---|
| `--debug` | `-d` | On evaluation failure, drop into the interactive shell with the last scope. Also set via `ARRAI_DEBUG=1`. |
| `--version` | `-v` | Print version, OS, arch. |
| `--help` | `-h` | Print help. |
| `--help-agent` | | Print this agent guide. |

### Direct file execution

When the first non-flag argument is not a known subcommand, `arrai` treats it
as a file path and runs it directly:

```bash
arrai path/to/file.arrai
arrai run path/to/file.arrai   # equivalent
arrai ./run                    # use ./ to disambiguate a file named "run"
```

### Commands

| Command | Aliases | Description |
|---|---|---|
| `shell` | `i` | Start the interactive REPL |
| `run` | `r` | Evaluate an arrai file or `.arraiz` bundle |
| `eval` | `e` | Evaluate an inline expression |
| `bundle` | `b` | Bundle a script and its imports into an `.arraiz` file |
| `compile` | `c` | Compile a script to a standalone binary |
| `test` | `t` | Run `*_test.arrai` test files |
| `json` | `jx` | Convert JSON from stdin to arrai representation |
| `serve` | `s` | Start a gRPC/WebSocket server |
| `observe` | `o` | Observe an expression on a running server |
| `update` | `u` | Update a running server's state |
| `sync` | `s` | Sync local files to a server |
| `transform` | `x` | Transform a stream of input data |
| `info` | | Display release information |

### Output modes (`--out` / `-o` flag)

Available on `eval`, `run`, and `bundle`:

| Form | Description |
|---|---|
| `--out filename` | Write result to a file |
| `--out f:filename` | Same (explicit file mode) |
| `--out d:dirname` | Write dict result as a directory tree |

## Language Basics

### The type system

Arr.ai has exactly three value kinds:

1. **Numbers** — 64-bit binary floats.
2. **Tuples** — Named collections. `(x: 1, y: 2)`.
3. **Sets** — Unordered collections with no duplicates. `{1, 2, 3}`.

Everything else is sugar:

| Sugar | Underlying form |
|---|---|
| `true` | `{()}` |
| `false` | `{}` |
| `"hello"` | `{(@:0, @char:104), (@:1, @char:101), ...}` |
| `[3, 9, 27]` | `{(@:0, @item:3), (@:1, @item:9), (@:2, @item:27)}` |
| `{"a": 1}` | `{(@:"a", @value:1)}` |
| `<<'hi'>>` | `{(@:0, @byte:104), (@:1, @byte:105)}` |
| `%a` | `97` (character literal = number) |

Relations are sets of tuples sharing the same attribute names — the arr.ai
equivalent of SQL tables.

### Literals

```arrai
# Numbers
42         1.5e3       -10

# Tuples
()                      # empty tuple
(x: 1, y: 2)
(:x, :y)                # shorthand for (x: x, y: y)

# Sets
{}                      # empty set / false
{1, 2, 3}

# Relation shorthand
{|x, y| (1, 2), (3, 4)}   # = {(x:1,y:2), (x:3,y:4)}

# Strings (all equivalent)
"hello"    'hello'    `hello`

# Expression string (template)
$"Hello ${name}!"
$`${[1,2,3]::, }`     # = "1, 2, 3"

# Arrays
[1, 2, 3]

# Dicts
{"a": 1, "b": 2}

# Bytes
<<'hello'>>
<<104, 101, 108, 108, 111>>
```

### Binding and functions

```arrai
let x = 42; x + 1                              # let-binding
\x x * 2                                       # lambda
\x \y x + y                                    # curried binary
let rec fib = \n 1 if n < 2 else n * fib(n-1); # named recursion
let [head, ...tail] = [1, 2, 3]; tail          # pattern matching
```

### Transform operators

| Operator | Description |
|---|---|
| `->` | Transform a single value |
| `=>` | Map over set members |
| `>>` | Map over sequence/dict values |
| `>>>` | Map with key and value |
| `:>` | Map over tuple attribute values |

```arrai
42 -> . + 1                        # 43
{2,4,6} => . * 2                   # {4,8,12}
[1,2,3] >> . * 10                  # [10,20,30]
(r:0.5, g:0.2, b:0.7) :> 1 - .    # (b:0.3, g:0.8, r:0.5)
```

### Set operators

```arrai
{1,2,3} | {3,4,5}                  # union: {1,2,3,4,5}
{1,2,3} & {3,4,5}                  # intersection: {3}
{1,2,3} &~ {3,4,5}                 # difference: {1,2}
{1,2,3} where . > 1                # filter: {2,3}
2 <: {1,2,3}                       # membership: true
```

### Relational operators

```arrai
A <&> B                             # natural join
A <-> B                             # compose (project away common attrs)
A <&- B                             # left semijoin
A -&> B                             # right semijoin
```

### Imports

```arrai
//str                              # standard library module
//{./other.arrai}                  # relative file import
//{/root/path/file}                # module-root-relative import
//encoding.json.decode(data)       # stdlib function call
```

Module root is the directory of the nearest `go.mod` file, searching upward.

### Conditional

```arrai
expr1 if test else expr2

cond (
    age < 40: "young",
    age < 60: "middle",
    *:        "old",
)
```

## Standard Library

All stdlib is accessed via `//`.

| Module | Key functions |
|---|---|
| `//str` | `lower`, `upper`, `title`, `repr`, `expand` |
| `//seq` | `concat`, `contains`, `has_prefix`, `has_suffix`, `join`, `split`, `sub`, `repeat`, `trim_prefix`, `trim_suffix` |
| `//math` | `pi`, `e`, `sin`, `cos`, `sqrt`, `inf` |
| `//re` | `compile(pattern)` → `{match, sub, subf}` |
| `//encoding.json` | `decode`, `encode`, `encode_indent` |
| `//encoding.yaml` | `decode`, `encode` |
| `//encoding.csv` | `decode`, `encode` |
| `//encoding.xml` | `decode`, `encode` |
| `//encoding.proto` | `descriptor`, `decode` |
| `//encoding.xlsx` | `decodeToRelation` |
| `//os` | `args`, `cwd`, `file`, `exists`, `tree`, `get_env`, `stdin` |
| `//net.http` | `get`, `post` |
| `//log` | `print`, `printf` |
| `//fmt` | `pretty` |
| `//fn` | `fix` (fixed-point combinator), `fixt` (mutual recursion) |
| `//rel` | `union` |
| `//bits` | `set`, `mask` |
| `//eval` | `value`, `eval`, `evaluator` |
| `//grammar` | `parse`, `lang.wbnf`, `lang.arrai` |
| `//dict` | Convert tuple to dict |
| `//tuple` | Convert dict to tuple |

## Testing

Test files must be named `*_test.arrai` and evaluate to a structure where all
leaves are `true` or `false`.

```arrai
# math_test.arrai
(
    addition:    1 + 1 = 2,
    subtraction: 5 - 3 = 2,
    nested: [
        2 * 2 = 4,
        2 * 3 = 6,
    ],
)
```

```bash
arrai test                              # all tests recursively
arrai test path/to/dir                  # specific directory
arrai test path/to/specific_test.arrai  # single file
```

## Common Patterns

### Read and transform JSON

```arrai
let data = //encoding.json.decode(//os.file('data.json'));
data('users') => .name -> //str.upper(.)
```

### HTTP request

```arrai
let resp = //net.http.get((), 'https://api.example.com/data');
//encoding.json.decode(resp.body)
```

### Code generation with expression strings

```arrai
let funcs = [
    (name: "square", params: ["x"], body: "x ^ 2"),
    (name: "sum",    params: ["x", "y"], body: "x + y"),
];
$`${funcs >> $`
    function ${.name}(${.params::, }) {
        return ${.body}
    }
`::\i:\n}`
```

### Write output to files

```bash
arrai run --out d:output/ script.arrai
```

The script must return a dict where keys are filenames and values are strings
or byte arrays.

## Architecture (for contributors)

### Evaluation pipeline

```
Source text → Parser (syntax/arrai.wbnf) → AST → Compiler (syntax/compile.go) → Expr → Eval(ctx, Scope) → Value
```

### Key packages

| Package | Description |
|---|---|
| `rel/` | Core types: `Value`, `Expr`, `Scope`, `Pattern`. Immutable values over `github.com/arr-ai/frozen`. |
| `syntax/` | Parser, compiler, stdlib (`std_*.go`), imports, bundler. |
| `cmd/arrai/` | CLI entry point (urfave/cli). |
| `pkg/` | REPL, test runner, bundle system, context-based filesystem. |
| `engine/` | Stateful server evaluation engine. |
| `translate/` | Format translators (protobuf, XML, YAML). |
