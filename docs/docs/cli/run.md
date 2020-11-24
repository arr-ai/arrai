Arr.ai's `run` (alias `r`) command is effectively the same as [eval](./eval), except it takes as input the path to a runnable arr.ai file instead of an expression.

```bash
$ cat hi.arrai
$"Hello, ${//os.args(1)}!"
$ arrai run hi.arrai Alice
Hello, Alice!
$ arrai r hi.arrai Alice
Hello, Alice!
```

For more details, see [eval](./eval).
