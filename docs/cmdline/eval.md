# arrai eval and run

Arr.ai's `eval` and `run` commands evaluate argument as an arr.ai expression
and, by default, print the result to stdout.

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

The `run` command evalutes the contents of the source file specified by its
argument.

```bash
$ cat hi.arrai
$"Hello, ${//os.args(1)}!"
$ arrai run hi.arrai Alice
Hello, Alice!
$ arrai r hi.arrai Alice
Hello, Alice!
```

## Output controls

The `eval` command supports the `--out` flag (shorthand `-o`), which changes
arr.ai's output behaviour as follows:

| option | description |
|-|-|
| `--out=file:<path>` | Outputs the evaluated result, which must be a string or byte array, to the file specified by `<path>`. `file` may be abbreviated to `f`. |
| `--out=dir:<path>` | Outputs the evaluated result to the directory specified by `<path>`. The result to be output must be a recursively nested dict of dicts, with strings and/or byte arrays at the leaves. Each dict represents the contents of a directory at the location given by its key while each string or byte array correspondingly represents the contents of a file. `dir` may be abbreviated to `d`.  |

For both `--out=file` and `--out=dir`, strings are UTF-8 encoded when written to the corresponding file.

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
    "bar": {"baz.txt": <<"Goodbye, ", name, "!", 10>>}
}' Bob
$ tree out
out
├── bar
│   └── baz.txt
└── foo.txt

1 directory, 2 files
$ cat out/foo.txt
Hello, Bob.
$ cat out/bar/baz.txt
Goodbye, Bob!
```
