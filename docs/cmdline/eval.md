# arrai eval

Arr.ai's `eval` command evaluates its argument as an arr.ai expression and, by
default, prints its value to stdout:

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

### Examples

```bash
arrai e --out=file:example.txt '"hello\n"'
arrai e -o f:example.txt '"hello\n"'
