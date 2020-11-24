The `re` library contains functions that are used for regular expression matching.
The library uses [RE2 syntax](https://github.com/google/re2/wiki/Syntax).

## `//re.compile(re <: string) <: tuple`

`compile` takes the string `re` representing a regular expression and returns a tuple
containing functions whose attributes are the following:

### `//re.compile(re <: string).match(str <: string) <: array_of_arrays`

`match` finds all matches of `re` in `str`, returning an
array of arrays. Each top-level array is an instance of a match against `re`.
Each second-level array is an array of captured submatches, with the first
element being the full match.

Submatches are expressed as offset-strings. That is, a match of the text `a(b)c`
at position 10 in `str` would produce the following match: `[10\"abc", 11\"b"]`,

Usage:

| example | equals |
|:-|:-|
| `//re.compile('.(\\d)').match('a1b2c3')` | `[['a1', 1\'1'], [2\'b2', 3\'2'], [4\'c3', 5\'3']]` |

### `//re.compile(re <: string).sub(replace <: string, str <: string) <: string`

`sub` replaces all matches of `re` in `str` with `replace`. Within `replace`,
any occurrences of `$n` will be replaced by capturing group `n` from the match.

Usage:

| example | equals |
|:-|:-|
| `//re.compile('.(\\d)').sub('-$1', 'a1b2c3')` | `'-1-2-3'` |
| `let nodigits = //re.compile('\\d+').sub('');`<br/>`nodigits('a1b2c3')` | `'abc'` |

### `//re.compile(re <: string).subf(f <: function, substring <: string) <: string`

`subf` replaces all matches of `re` in `str` with the result of calling `f(match)`.

Usage:

| example | equals |
|:-|:-|
| `//re.compile('.(\\d)').subf(//str.upper, 'a1b2c3')` | `'A1B2C3'` |
