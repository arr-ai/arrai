# seq library

seq library contains functions that are used for string manipulations.

## concat

`//.seq.concat(strings)` takes a list of strings and returns a string that is
the concatenation of the strings in the list.

| example | returns |
|:-|:-|
| `//.seq.concat(["ba", "na", "na"])` | `"banana"` |

## repeat

`//.seq.repeat(n, seq)` returns a sequence that contains `seq` repeated `n`
times.

| example | returns |
|:-|:-|
| `//.seq.repeat(2, "hots")` | `"hotshots"` |
