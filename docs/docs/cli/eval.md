Arr.ai's `eval` (alias `e`) command evaluates its argument as an arr.ai expression
and, by default, prints the result to stdout.

The `eval` command evaluates its argument as an arr.ai expression.

```bash
$ arrai eval '2 + 2'
4
```

The `e` command is a shortcut and behaves identically to `eval`.

```bash
$ arrai e '2 + 2'
4
```

## Output controls

The `eval` command supports the `--out` flag (shorthand `-o`), which changes
arr.ai's output behaviour as follows:

| option | description |
|-|-|
| `--out=file:<path>` | Outputs the evaluated result, which must be a string or byte array, to the file specified by `<path>`. `file` may be abbreviated to `f`. |
| `--out=dir:<path>` | Outputs the evaluated result to the directory specified by `<path>`. The result to be output must be a recursively nested dict of dicts, with strings and/or byte arrays at the leaves. Each dict represents the contents of a directory at the location given by its key while each string or byte array correspondingly represents the contents of a file. `dir` may be abbreviated to `d`.  |

For both `--out=file` and `--out=dir`, strings are UTF-8 encoded when written to the corresponding file.

### More fine-grained control

An arr.ai program can exercise more control over the `--out=dir` option by
returning tuples in place of dicts (for directories) or strings (for files). The
structure of these tuples is as follows:

| atttribute | value(s) | purpose |
|-|-|-|
| `ifExists` | `'fail'` \| `'ignore'` \| `'merge'` \| `'remove'` \| `'replace'` | Control the behaviour when encountering existing directory entries in the target directory. If `ifExists` is omitted, it defaults to `'merge'` if the `dir` attribute is present, or `'replace'` if the `file` attribute is present. |
| `dir` | a dict | Output the dict a directory. |
| `file` | a string or byte array | Output the content as a file. |

The tuple `(dir: {})` can be used to output an empty directory with , since `{}`
by itself is indistinguishable from an empty file, `""`.

Strings, byte arrays and dicts are in fact shorthands for the tuple form, as
follows:

| value | example | equivalent |
|-|-|-|
| string or byte array | "hello" | `(file: "hello")` |
| dict | {"foo.txt": "hello"} | `(dir: {"foo.txt": "hello"})` |

### Examples

```bash
$ arrai e --out=file:example.txt '$"Hello, ${//os.args(1)}!\n"' Bob
$ cat example.txt
Hello, Bob!
$ arrai e -o f:example.txt '"hello\n"'
$ cat example.txt
hello
$ arrai e -o :example.txt '"hello\n"'
$ cat example.txt
hello
$ arrai e -o example.txt '"hello\n"'
$ cat example.txt
hello
```

Note that `--out=file` offers several shortenings. The final example above only
works if the path itself doesn't have a `:` in it.

```bash
$ arrai e --out=dir:out 'let [_, name, ...] = //os.args;
{
    "foo.txt": $"Hello, ${name}.\n",
    "bar": {
        "baz.txt": <<"Goodbye, ", name, "!", 10>>,
        "empty": (
            ifExists: "replace",
            dir: {},
        ),
        "template.c": (
            ifExists: "remove",
        ),
        "template2.c": (
            ifExists: "ignore",
            file: $`
                int main() {
                    // Write some code here...
                    return 1;
                }
            `,
        ),
    }
}' Bob
$ tree out
out
├── bar
│   └── baz.txt
│   └── empty
│   └── template2.c
└── foo.txt

1 directory, 2 files
$ cat out/foo.txt
Hello, Bob.
$ cat out/bar/baz.txt
Goodbye, Bob!
$ cat out/bar/template2.c
$ cat out/bar/template2.c
int main() {
    // Write some code here...
    return 1;
}
```
