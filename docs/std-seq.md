# seq

The `seq` library contains functions that are used for manipulating sequenced data structures.

## String

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


## `//seq.has_prefix(prefix <: string, s <: string) <: bool`

`has_prefix` checks whether the string `s` is prefixed by `prefix`. It returns a boolean.

Usage:

| example | equals |
|:-|:-|
| `//seq.has_prefix("I'm", "I'm running out of stuff to write")` | `true` |
| `//seq.has_prefix("to write", "I'm running out of stuff to write")` | `{}` which is equal to `false` |

## `//seq.has_suffix(suffix <: string, s <: string) <: bool`

`has_suffix` checks whether the string `s` is suffixed by `suffix`. It returns a boolean.

Usage:

| example | equals |
|:-|:-|
| `//seq.has_suffix("I'm", "I'm running out of stuff to write")` | `{}` which is equal to `false` |
| `//seq.has_suffix("to write", "I'm running out of stuff to write")` | `true` |

## `//seq.join(delimiter <: string, s <: array_of_string) <: string`

`join` returns a concatenated string with each member of `s` delimited by `delimiter`

Usage:

| example | equals |
|:-|:-|
| `//seq.join(", ", ["pew", "another pew", "and more pews"])` | `"pew, another pew, and more pews"` |
| `//seq.join(" ", ["this", "is", "a", "sentence"])` | `"this is a sentence"` |
| `//seq.join(["", "this", "is", "a", "sentence"])` | `"thisisasentence"` |

## `//seq.contains(substr <: string, str <: string) <: bool`

`contains` checks whether `substr` is contained in `str`. It returns a
boolean.

Usage:

| example | equals |
|:-|:-|
| `//seq.contains("substring", "the full string which has substring")` | `true` |
| `//seq.contains("microwave", "just some random sentence")` | `{}` which is equal to `false` |

## `//seq.sub(old <: string, new <: string, s <: string) <: string`

`sub` replaces occurrences of `old` in `s`  with `new`. It returns the modified string.

Usage:

| example | equals |
|:-|:-|
| `//seq.sub("old string", "new sentence", "this is the old string")` | `"this is the new sentence"` |
| `//seq.sub("string", "stuff", "just another sentence")` | `"just another sentence"` |

## `//seq.split(delimiter <: string, s <: string) <: array of string`

`split` splits the string `s` based on the provided `delimiter`. It returns an array of strings
which are split from the string `s`.

Usage:

| example | equals |
|:-|:-|
| `//seq.split(" ", "deliberately adding spaces to demonstrate the split function")` | `["deliberately", "adding", "spaces", "to", "demonstrate", "the", "split", "function"]` |
| `//seq.split("random stuff", "this is just a random sentence")` | `["this is just a random sentence"]` |

## Array

## `//seq.has_prefix(prefix <: array, subject <: array) <: bool`

`has_prefix` checks whether the array `subject` is prefixed by array `prefix`. It returns a boolean.

Usage:

| example | equals |
|:-|:-|
| `//seq.has_prefix(['A'],['A','B','C'])` | `true` |
| `//seq.has_prefix([1, 2],[1, 2, 3])` | `true` |
| `//seq.has_prefix([[1, 2]],[[1, 2], [3]])` | `true` |

## `//seq.has_suffix(suffix <: array, subject <: array) <: bool`

`has_suffix` checks whether the array `subject` is suffixed by array `suffix`. It returns a boolean.

Usage:

| example | equals |
|:-|:-|
| `//seq.has_suffix(['E'],['A','B','C','D','E'])` | `true` |
| `//seq.has_suffix([[3, 4]],[[1 ,2], [3, 4]])` | `true` |

## `//seq.join(joiner <: array, subject <: array) <: array`

`join` returns a concatenated array with each member of `subject` delimited by `joiner`

Usage:

| example | equals |
|:-|:-|
| `//seq.join([0], [1,2,3,4,5])` | `[1,0,2,0,3,0,4,0,5]` |
| `//seq.join([0], [[2, [3, 4]], [5, 6]])` | `[2, [3, 4], 0, 5, 6]` |
| `//seq.join([[0],[1]], [[[1, 2], [3, 4]],[[5, 6],[7, 8]]])` | `[[1, 2], [3, 4], [0], [1], [5, 6], [7, 8]]` |

## `//seq.contains(sub <: array, subject <: array) <: bool`

`contains` checks whether array `sub` is contained in array `subject`. It returns a boolean.

Usage:

| example | equals |
|:-|:-|
| `//seq.contains([1,2,3,4,5], [1,2,3,4,5])` | `true` |
| `//seq.contains([['B','C']],[['A', 'B'], ['B','C'],['D','E']])` | `true` |

## `//seq.sub(old <: array, new <: array, subject <: array) <: array`

`new` replaces occurrences of array `old` in array `subject` with array `new`. It returns the modified array.

Usage:

| example | equals |
|:-|:-|
| `//seq.sub([1], [2], [1, 2, 3])` | `[2, 2, 3]` |
| `//seq.sub([[2,2]], [[4,4]], [[1,1], [2,2], [3,3]])`| `[[1,1], [4,4], [3,3]]` |

## `//seq.split(delimiter <: array, subject <: array) <: array`

`split` splits the array `subject` based on the provided array `delimiter`. It returns an array
which are split from the array `subject`.

Usage:

| example | equals |
|:-|:-|
| `//seq.split([1],[1, 2, 3])` | `[[],[2,3]]` |
| `//seq.split([3],[1, 2, 3])` | `[[1,2],[]]` |
| `//seq.split(['A'],['B', 'A', 'C', 'A', 'D', 'E'])` | `[['B'],['C'], ['D', 'E']]` |
| `//seq.split([['C','D'],['E','F']],[['A','B'], ['C','D'], ['E','F'], ['G']])`) | `[[['A','B']], [['G']]]` |
