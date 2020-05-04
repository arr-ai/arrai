# unicode

The `unicode` library contains functions for operations related to unicode.

## `//unicode.utf8 <: tuple`

The `utf8` tuple contains functions that are used with utf8-encoded data.

### `//unicode.utf8.encode(s <: string) <: array_of_bytes`

`encode` encodes the string `s` as a UTF8 bytes sequence.

Usage:

| example |
|:-|
| `//unicode.utf8.encode("abc")` |
