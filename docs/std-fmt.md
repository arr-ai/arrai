# fmt

The `fmt` library contains helper functions for formatting strings.

## `//fmt.pretty(v <: value) <: string`

Returns a pretty (i.e. with newlines and indents) string representation of the given value.

Usage:

```bash
@> //fmt.pretty({'a': (b: [1,2])})
{
  'a': (
    b: [
      1,
      2
    ]
  )
```
