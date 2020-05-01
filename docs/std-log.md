# log

The `log` library contains functions for printing values to the
standard output.

## `print(v <: any) <: any`

It takes the value `v` of any type and print it to the standard
output. It returns `v` itself.

Usage:

| example |
|:-|
| `//log.print(123)` |
| `//log.print("string")` |
| `//log.print([1, "b", "pew"])` |

## `printf(format <: string, values <: array_of_any) <: any`

It prints each member of `values` corresponding to `format`. It
returns the first member of `values`.

The syntax for `format` uses [golang's syntax](https://pkg.go.dev/fmt?tab=doc#hdr-Printing).

| example |
|:-|
| `//log.printf("My laser went %s and %s", ["pew", "pew again"])` |
| `//log.printf("I bought %d watermelons, each one is worth %d %s", [2000, 200, "dollars"])` |
