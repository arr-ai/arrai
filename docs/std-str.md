# str

The `str` library contains functions that are used for string manipulations. 

Note: many of the functions typically found in a string standard library are found in the [seq](./std-seq.md) library instead as they work for any sequenced data structure including strings. These include standard operations such as `//seq.contains`, `//seq.concat`,`//seq.join`, `//seq.sub`, `//seq.split`, `//seq.trim_prefix`, `//seq.trim_suffix`.

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
