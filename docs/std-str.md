# str

`str` library contains functions that are used for string manipulations.

## `contains(str <: string, substr <: string) <: bool`

Checks whether `substr` is contained in `str`. It returns a
boolean.

Usage:

| example | equals |
|:-|:-|
| `//str.contains("the full string which has substring", "substring")` | `true` |
| `//str.contains("just some random sentence", "microwave")` | `{}` which is equal to `false` |

## `sub(s <: string, old <: string, new <: string) <: string`

Replaces occurrences of `old` in `s`  with `new`. It returns the modified string.

Usage:

| example | equals |
|:-|:-|
| `//str.sub("this is the old string", "old string", "new sentence")` | `"this is the new sentence"` |
| `//str.sub("just another sentence", "string", "stuff")` | `"just another sentence"` |

## `split(s <: string, delimiter <: string) <: array of string`

Split the string `s` based on the provided `delimiter`. It returns an array of string
which are splitted from the string `s`.

Usage:

| example | equals |
|:-|:-|
| `//str.split("deliberately adding spaces to demonstrate the split function", " ")` | `["deliberately", "adding", "spaces", "to", "demonstrate", "the", "split", "function"]` |
| `//str.split("this is just a random sentence", "random stuff")` | `["this is just a random sentence"]` |

## `lower(s <: string) <: string`

Returns the string `s` except all of the character is converted to its lowercase form.

Usage:

| example | equals |
|:-|:-|
| `//str.lower("HeLLo ThErE")` | `"hello there"` |
| `//str.lower("GENERAL KENOBI WHAT A SURPRISE")` | `"general kenobi what a surprise"` |
| `//str.lower("123")` | `"123"` |

## `upper(s <: string) <: string`

Returns the string `s` except all of the character is converted to its uppercase form.

Usage:

| example | equals |
|:-|:-|
| `//str.upper("HeLLo ThErE")` | `"HELLO THERE"` |
| `//str.upper("did you ever hear the tragedy of darth plagueis the wise")` | `"DID YOU EVER HEAR THE TRAGEDY OF DARTH PLAGUEIS THE WISE"` |
| `//str.upper("321")` | `"321"` |

## `title(s: string) <: string`

Returns the string `s` except all the first letters of each words delimited by
a white space are capitalized.

Usage:

| example | equals |
|:-|:-|
| `//str.title("laser noises pew pew pew")` | `"Laser Noises Pew Pew Pew"` |
| `//str.title("pew")` | `"Pew"` |

## `has_prefix(s <: string, prefix <: string) <: bool`

Checks whether the string `s` is prefixed by `prefix`. It returns a boolean.

Usage:

| example | equals |
|:-|:-|
| `//str.has_prefix("I'm running out of stuff to write", "I'm")` | `true` |
| `//str.has_prefix("I'm running out of stuff to write", "to write")` | `{}` which is equal to `false` |

## `has_suffix(s <: string, suffix <: string) <: bool`

Checks whether the string `s` is suffixed by `suffix`. It returns a boolean.

Usage:

| example | equals |
|:-|:-|
| `//str.has_suffix("I'm running out of stuff to write", "I'm")` | `{}` which is equal to `false` |
| `//str.has_suffix("I'm running out of stuff to write", "to write")` | `true` |

## `join(s <: array_of_string, delimiter <: string) <: string`

It returns a string which is a concatenated string of each member of `s` delimited
by the `delimiter`

Usage:

| example | equals |
|:-|:-|
| `//str.join(["pew", "another pew", "and more pews"], ", ")` | `"pew, another pew, and more pews"` |
| `//str.join(["this", "is", "a", "sentence"], " ")` | `"this is a sentence"` |
| `//str.join(["this", "is", "a", "sentence"], "")` | `"thisisasentence"` |
