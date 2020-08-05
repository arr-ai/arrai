# fmt

The `fmt` library contains helper functions for formatting strings.

## `//fmt.pretty(v <: value) <: string`

Returns a pretty (i.e. with newlines and indents) string representation of the given value.

Usage:

```bash
$ arrai eval "//fmt.pretty({'a': (b: [1,2])})"
{
  'a': (
    b: [1, 2]
  )
}
```

Note: the rendering of output strings differs between `arrai eval` and the interactive shell. As such, the result of `//fmt.pretty` in the interactive shell will not appear very pretty:

<!-- TODO: Update once `/print` is implemented. --> 

```bash
@> //fmt.pretty({'a': (b: [1,2])})
"{\n  'a': (\n    b: [1, 2]\n  )\n}"
```
