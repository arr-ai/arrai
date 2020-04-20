# seq library

seq library contains functions that are used for string manipulations.

## concat

`//.seq.concat(seqs)` takes an array of sequences and returns a sequence that is
the concatenation of the sequences in the array.

| example | returns |
|:-|:-|
| `//.seq.concat(["ba", "na", "na"])` | `"banana"` |
| `//.seq.concat([[1, 2], [3, 4, 5]])` | `[1, 2, 3, 4, 5]` |

## repeat

`//.seq.repeat(n, seq)` returns a sequence that contains `seq` repeated `n`
times.

| example | returns |
|:-|:-|
| `//.seq.repeat(2, "hots")` | `"hotshots"` |
