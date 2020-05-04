# seq

The `seq` library contains functions that are used for string manipulations.

## `//seq.concat(seqs <: array) <: array` <br/> `concat(seqs <: string) <: string`

`concat` takes an array of sequences `seqs` and returns a sequence that is
the concatenation of the sequences in the array.

| example | equals |
|:-|:-|
| `//seq.concat(["ba", "na", "na"])` | `"banana"` |
| `//seq.concat([[1, 2], [3, 4, 5]])` | `[1, 2, 3, 4, 5]` |

## `//seq.repeat(n <: number, seq <: array) <: array` <br/> `repeat(n <: number, seq <: string) <: string`

`repeat` returns a sequence that contains `seq` repeated `n` times.

| example | equals |
|:-|:-|
| `//seq.repeat(2, "hots")` | `"hotshots"` |
