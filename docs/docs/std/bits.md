The `bits` library contains functions related to bit operations.

## `//bits.set(n <: number) <: set`

`set` takes in the bitmask `n` and returns a set of integers representing `n`.
The set itself contains the positions of 1-bits in the binary representation
of `n`

`n` must be a non-negative number.

Usage:

| example | equals |
|:-|:-|
|`//bits.set(42)` | `{1, 3, 5}` |
|`//bits.set(72061992084439040)` | `{42, 56}` |
|`//bits.set(0)` | `{}` |

## `//bits.mask(s <: set) <: number`

`mask` is the inverse of `set`. It takes a set `s` which represents the position
of 1-bits in a binary representation of a bitmask and returns the bitmask.

`s` must be a set of numbers.

Usage:

| example | equals |
|:-|:-|
|`//bits.mask({1, 3, 5})` | `42` |
|`//bits.mask({42, 56})` | `72061992084439040` |
|`//bits.mask({})` | `0` |
