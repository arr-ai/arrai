# str

The `str` library contains functions that are used for string manipulations.

## `//str.contains(str <: string, substr <: string) <: bool`

`contains` checks whether `substr` is contained in `str`. It returns a
boolean.

Usage:

| example | equals |
|:-|:-|
| `//str.contains("the full string which has substring", "substring")` | `true` |
| `//str.contains("just some random sentence", "microwave")` | `{}` which is equal to `false` |

## `//str.sub(s <: string, old <: string, new <: string) <: string`

`sub` replaces occurrences of `old` in `s`  with `new`. It returns the modified string.

Usage:

| example | equals |
|:-|:-|
| `//str.sub("this is the old string", "old string", "new sentence")` | `"this is the new sentence"` |
| `//str.sub("just another sentence", "string", "stuff")` | `"just another sentence"` |

## `//str.split(s <: string, delimiter <: string) <: array of string`

`split` splits the string `s` based on the provided `delimiter`. It returns an array of strings
which are split from the string `s`.

Usage:

| example | equals |
|:-|:-|
| `//str.split("deliberately adding spaces to demonstrate the split function", " ")` | `["deliberately", "adding", "spaces", "to", "demonstrate", "the", "split", "function"]` |
| `//str.split("this is just a random sentence", "random stuff")` | `["this is just a random sentence"]` |

## `//str.lower(s <: string) <: string`

`lower` returns the string `s` with all of the character converted to lowercase.

Usage:

| example | equals |
|:-|:-|
| `//str.lower("HeLLo ThErE")` | `"hello there"` |
| `//str.lower("GENERAL KENOBI WHAT A SURPRISE")` | `"general kenobi what a surprise"` |
| `//str.lower("123")` | `"123"` |

## `//str.upper(s <: string) <: string`

`upper` returns the string `s` with all of the character converted to uppercase.

Usage:

| example | equals |
|:-|:-|
| `//str.upper("HeLLo ThErE")` | `"HELLO THERE"` |
| `//str.upper("did you ever hear the tragedy of darth plagueis the wise")` | `"DID YOU EVER HEAR THE TRAGEDY OF DARTH PLAGUEIS THE WISE"` |
| `//str.upper("321")` | `"321"` |

## `//str.title(s: string) <: string`

`title` returns the string `s` with all the first letter of each word delimited by
a white space capitalised.

Usage:

| example | equals |
|:-|:-|
| `//str.title("laser noises pew pew pew")` | `"Laser Noises Pew Pew Pew"` |
| `//str.title("pew")` | `"Pew"` |

## `//str.has_prefix(s <: string, prefix <: string) <: bool`

`has_prefix` checks whether the string `s` is prefixed by `prefix`. It returns a boolean.

Usage:

| example | equals |
|:-|:-|
| `//str.has_prefix("I'm running out of stuff to write", "I'm")` | `true` |
| `//str.has_prefix("I'm running out of stuff to write", "to write")` | `{}` which is equal to `false` |

## `//str.has_suffix(s <: string, suffix <: string) <: bool`

`has_suffix` checks whether the string `s` is suffixed by `suffix`. It returns a boolean.

Usage:

| example | equals |
|:-|:-|
| `//str.has_suffix("I'm running out of stuff to write", "I'm")` | `{}` which is equal to `false` |
| `//str.has_suffix("I'm running out of stuff to write", "to write")` | `true` |

## `//str.join(s <: array_of_string, delimiter <: string) <: string`

`join` returns a concatenated string with each member of `s` delimited by `delimiter`

Usage:

| example | equals |
|:-|:-|
| `//str.join(["pew", "another pew", "and more pews"], ", ")` | `"pew, another pew, and more pews"` |
| `//str.join(["this", "is", "a", "sentence"], " ")` | `"this is a sentence"` |
| `//str.join(["this", "is", "a", "sentence"], "")` | `"thisisasentence"` |
