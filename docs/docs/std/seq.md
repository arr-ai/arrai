The `seq` library contains functions that are used for manipulating sequenced data structures including string and array.

## `//seq.concat(seqs <: array) <: array` <br/> `//seq.concat(seqs <: string) <: string`

`concat` takes an array of sequences `seqs` and returns a sequence that is
the concatenation of the sequences in the array.

Usage:
| example | equals |
|:-|:-|
| `//seq.concat(["ba", "na", "na"])` | `"banana"` |
| `//seq.concat([[1, 2], [3, 4, 5]])` | `[1, 2, 3, 4, 5]` |

## `//seq.contains(sub <: array, subject <: array) <: bool` <br/> `//seq.contains(sub <: string, subject <: string) <: bool`

`contains` checks whether sequence `sub` is contained in sequence `subject` and returns true if it is, or false otherwise.

Usage:
| example | equals |
|:-|:-|
| `//seq.contains("substring", "the full string which has substring")` | `true` |
| `//seq.contains("microwave", "just some random sentence")` | `false` |
| `//seq.contains([1,2,3,4,5], [1,2,3,4,5])` | `true` |
| `//seq.contains([['B','C']],[['A', 'B'], ['B','C'],['D','E']])` | `true` |

## `//seq.has_prefix(prefix <: array, subject <: array) <: bool` <br/> `//seq.has_prefix(prefix <: string, subject <: string) <: bool` <br/> `//seq.has_prefix(prefix <: bytes, subject <: bytes) <: bool`

`has_prefix` checks whether the sequence `subject` is prefixed by sequence `prefix` and returns true if it is, or false otherwise.

Usage:
| example | equals |
|:-|:-|
| `//seq.has_prefix("I'm", "I'm running out of stuff to write")` | `true` |
| `//seq.has_prefix("to write", "I'm running out of stuff to write")` | `false` |
| `//seq.trim_prefix(<<'dive'>>, <<"divesting">>)` | `true` |
| `//seq.has_prefix(['A'],['A','B','C'])` | `true` |
| `//seq.has_prefix([1, 2],[1, 2, 3])` | `true` |
| `//seq.has_prefix([[1, 2]],[[1, 2], [3]])` | `true` |


## `//seq.has_suffix(suffix <: array, subject <: array) <: bool` <br/> `//seq.has_suffix(suffix <: string, subject <: string) <: bool` <br/> `//seq.has_suffix(suffix <: bytes, subject <: bytes) <: bool`

`has_suffix` checks whether the sequence `subject` is suffixed by sequence `suffix` and returns true if it is, or false otherwise.

Usage:
| example | equals |
|:-|:-|
| `//seq.has_suffix("I'm", "I'm running out of stuff to write")` | `false` |
| `//seq.has_suffix("to write", "I'm running out of stuff to write")` | `true` |
| `//seq.has_suffix(<<'ary'>>, <<'binary'>>)` | `true` |
| `//seq.has_suffix(['E'],['A','B','C','D','E'])` | `true` |
| `//seq.has_suffix([[3, 4]],[[1 ,2], [3, 4]])` | `true` |

## `//seq.join(joiner <: array, subject <: array) <: array` <br/> `//seq.join(joiner <: string, subject <: array_of_string) <: string`

`join` returns a concatenated sequence with each member of sequence `subject` delimited by sequence `joiner`

Usage:
| example | equals |
|:-|:-|
| `//seq.join(", ", ["pew", "another pew", "and more pews"])` | `"pew, another pew, and more pews"` |
| `//seq.join(" ", ["this", "is", "a", "sentence"])` | `"this is a sentence"` |
| `//seq.join(["", "this", "is", "a", "sentence"])` | `"thisisasentence"` |
| `//seq.join([0], [[1, 2], [3, 4], [5, 6]]` | `[1, 2, 0, 3, 4, 0, 5, 6]` 
| `//seq.join([0], [[2, [3, 4]], [5, 6]])` | `[2, [3, 4], 0, 5, 6]` |
| `//seq.join([[0],[1]], [[[1, 2], [3, 4]],[[5, 6],[7, 8]]])` | `[[1, 2], [3, 4], [0], [1], [5, 6], [7, 8]]` |

## `//seq.split(delimiter <: array, subject <: array) <: array` <br/> `//seq.split(delimiter <: string, subject <: string) <: array of string`

`split` splits sequence `subject` based on the provided sequence `delimiter`. It returns an array of sequence which are split from the sequence `subject`.

Usage:
| example | equals |
|:-|:-|
| `//seq.split(" ", "deliberately adding spaces to demonstrate the split function")` | `["deliberately", "adding", "spaces", "to", "demonstrate", "the", "split", "function"]` |
| `//seq.split("random stuff", "this is just a random sentence")` | `["this is just a random sentence"]` |
| `//seq.split([1],[1, 2, 3])` | `[[],[2,3]]` |
| `//seq.split([3],[1, 2, 3])` | `[[1,2],[]]` |
| `//seq.split(['A'],['B', 'A', 'C', 'A', 'D', 'E'])` | `[['B'],['C'], ['D', 'E']]` |
| `//seq.split([['C','D'],['E','F']],[['A','B'], ['C','D'], ['E','F'], ['G']])`) | `[[['A','B']], [['G']]]` |

## `//seq.sub(old <: array, new <: array, subject <: array) <: array` <br/> `//seq.sub(old <: string, new <: string, subject <: string) <: string`

`sub` replaces occurrences of sequence `old` in sequence `subject` with sequence `new`. It returns the modified sequence.

Usage:
| example | equals |
|:-|:-|
| `//seq.sub("old string", "new sentence", "this is the old string")` | `"this is the new sentence"` |
| `//seq.sub("string", "stuff", "just another sentence")` | `"just another sentence"` |
| `//seq.sub([1], [2], [1, 2, 3])` | `[2, 2, 3]` |
| `//seq.sub([[2,2]], [[4,4]], [[1,1], [2,2], [3,3]])`| `[[1,1], [4,4], [3,3]]` |

## `//seq.repeat(n <: number, seq <: array) <: array` <br/> `//seq.repeat(n <: number, seq <: string) <: string`

`repeat` returns a sequence that contains `seq` repeated `n` times.

Usage:
| example | equals |
|:-|:-|
| `//seq.repeat(2, "hots")` | `"hotshots"` |

## `//seq.trim_prefix(prefix <: array, subject <: array) <: array` <br/> `//seq.trim_prefix(prefix <: string, subject <: string) <: string` <br/> `//seq.trim_prefix(prefix <: bytes, subject <: bytes) <: bytes`

`trim_prefix` checks whether the sequence `subject` is prefixed by sequence `prefix` and returns `subject` with `prefix` removed, otherwise it returns `subject` unmodified. It will only remove one copy of `prefix`.

Usage:
| example | equals |
|:-|:-|
| `//seq.trim_prefix("I'm", "I'm running out of stuff to write")` | `" running out of stuff to write"` |
| `//seq.trim_prefix("to write", "I'm running out of stuff to write")` | `"I'm running out of stuff to write"` |
| `//seq.trim_prefix(<<'dive'>>, <<"divesting">>)` | `<<'sting'>>` |
| `//seq.trim_prefix(['A'],['A','B','C'])` | `['B','C']` |
| `//seq.trim_prefix([1, 2],[1, 2, 3])` | `[3]` |
| `//seq.trim_prefix([[1, 2]],[[1, 2], [3]])` | `[[3]]` |


## `//seq.trim_suffix(suffix <: array, subject <: array) <: array` <br/> `//seq.trim_suffix(suffix <: string, subject <: string) <: string` <br/> `//seq.trim_suffix(suffix <: bytes, subject <: bytes) <: bytes`

`trim_suffix` checks whether the sequence `subject` is suffixed by sequence `suffix` and returns `subject` with `suffix` removed, otherwise it returns `subject` unmodified. It will only remove one copy of `suffix`.

Usage:
| example | equals |
|:-|:-|
| `//seq.trim_suffix("I'm", "I'm running out of stuff to write")` | `"I'm running out of stuff to write"` |
| `//seq.trim_suffix("to write", "I'm running out of stuff to write")` | `"I'm running out of stuff "` |
| `//seq.trim_suffix(<<'ary'>>, <<'binary'>>)` | `<<'bin'>>` |
| `//seq.trim_suffix(['E'],['A','B','C','D','E'])` | `['A','B','C','D']` |
| `//seq.trim_suffix([[3, 4]],[[1 ,2], [3, 4]])` | `[[1 ,2]]` |
