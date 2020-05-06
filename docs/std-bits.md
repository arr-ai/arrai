# bits

The `bits` library contains functions related to bit operations.

## `//bits.set(n <: number) <: set`

`set` takes in the number `n` and returns a set representing the
bitmask of `n`. The set itself contains positions of `1`-bits of
the binary representation of `n`

`n` must be non-negative number.

Usage:

| example | equals |
|:-|:-|
|`//bits.set(42)` | `{1, 3, 5}` |
|`//bits.set(72061992084439040)` | `{42, 56}` |
|`//bits.set(0)` | `{}` |

## `//bits.mask(s <: set) <: number`

`mask` does the opposite of set. It takes a bitmask in the form
of a set of numbers `s` and returns a numerical value represented
by the bitmask.

`s` must be a set of numbers.

Usage:

| example | equals |
|:-|:-|
|`//bits.mask({1, 3, 5})` | `42` |
|`//bits.mask({42, 56})` | `72061992084439040` |
|`//bits.mask({})` | `0` |
