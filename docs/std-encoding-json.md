# str library

str library contains functions that are used for string manipulations.

## `//.encoding.json.decode(string|bytes)`

Decodes a JSON-encoded value. Because empty sets are indisinguishable, `decode`
maps incoming JSON values as follows:

| JSON encoding | maps to&hellip; | comments |
|:-|:-|:-|
| `"abc"` | `(s: "abc")` |
| `[1, 2, 3]` | `(a: [1, 2, 3])` |
| `false`/`true` | `(b: false)`/`(b: true)` |
| `null` | `(null: {})` |
| `{"a": [2, 4, 8]}` | `{"a": (a: [2, 4, 8])}` | Being so common, objects are mapped directly to dicts. |
| `42` | `42` | Numbers, including zero, cannot be confused with other values. |
