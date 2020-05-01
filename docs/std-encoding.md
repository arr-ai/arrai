# encoding

The following functions are available by accessing the `//encoding` attribute.
It contains a function that decodes a JSON value into a built-in arrai values.

## `encoding.json.decode(json <: string|bytes) <: set`

It takes either a raw `string` or `bytes` that represents a JSON object. `json`
is then converted to a built-in arrai value.

Because empty sets are indistinguishable to `""`, `false`, and `[]`, `decode`
maps incoming JSON values as follows:

| JSON encoding | maps to&hellip; | notes |
|:-|:-|:-|
| `"abc"` | `(s: "abc")` |
| `[1, 2, 3]` | `(a: [1, 2, 3])` |
| `false`/`true` | `(b: false)`/`(b: true)` |
| `null` | `(null: {})` |
| `{"a": [2, 4, 8]}` | `{"a": (a: [2, 4, 8])}` | Objects are mapped directly to dicts. |
| `42` | `42` | Numbers, including zero, cannot be confused with other values. |

Usage:

| example | equals |
|:-|:-|
| `//encoding.json.decode('{"hi": "abc", "hello": 123}')` | `{'hello': 123, 'hi': (s: 'abc')}` |
