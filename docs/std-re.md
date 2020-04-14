# re library

The re library contains functions that are used for regular expression matching.
The library uses [RE2 syntax](https://github.com/google/re2/wiki/Syntax).

## match

`//.re.compile(re).match(str)` Finds all matches of `re` in `str`, returning an
array of arrays. Each top-level array is an instance of a match against `re`.
Each second-level array is an array of captured submatches, with the first
element being the full match.

Submatches are expressed as offset-strings. That is, a match of the text `a(b)c`
at position 10 in `str` would produce the following match: `[10\"abc", 11\"b"]`,

| example | equals |
|:-|:-|
| `//.re.compile('.(\\d)').match('a1b2c3')` | `[['a1', 1\'1'], [2\'b2', 3\'2'], [4\'c3', 5\'3']]` |

## sub

`//.re.compile(re).sub(replace, str)` replaces all matches of `re` in `str` with
`replace`. Within `replace`, any occurrences of `$n` will be replaced by
capturing group `n` from the match.

| example | equals |
|:-|:-|
| `//.re.compile('.(\\d)').sub('-$1', 'a1b2c3')` | `'-1-2-3'` |
| `let nodigits = //.re.compile('\\d+').sub('');`<br/>`nodigits('a1b2c3')` | `'abc'` |

## subf

`//.str.compile(re).subf(f, substring)` replaces all matches of `re` in `str` with
the result of calling `f(match)`.

| example | equals |
|:-|:-|
| `//.re.compile('.(\\d)').subf(//.str.upper, 'a1b2c3')` | `'A1B2C3'` |
