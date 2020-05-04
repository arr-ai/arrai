# os

The `os` contains functions that are related to the operating system.

## `os.args <: array_of_strings`

`//os.args` returns an array of strings representing the arguments provided on
the command line when running an arrai program. `//os.args(0)` is the path of
the arr.ai program that was invoked from the command line.

Usage:

| example | equals |
|:-|:-|
|`//os.args` | `["arg0", "arg1", "arg2", ...]` |

## `os.cwd <: string`

`//os.cwd` returns a string representing the current user directory.

## `os.file(filepath <: string) <: array_of_bytes`

Returns the content of a file located at `filepath` in the form of an array of bytes.

Usage:

| example | equals |
|:-|:-|
|`//os.file("path/to/file")` | `{ |@, @byte| ... }` |
|`//os.file("/absolute/path/to/file")` | `{ |@, @byte| ... }` |

## `os.get_env(key <: string) <: string`

Returns the environment variable that corresponds to `key` in the form of a `string`.

Usage:

| example | equals |
|:-|:-|
| `//os.get_env("KEY")` | `"string_value"` |

## `os.path_separator <: string`

Returns the path separator of the current operating system.
`/` for UNIX-like machine and `\` for Windows machine. It returns a string.

Usage:

| example | equals |
|:-|:-|
| `//os.path_separator` | `"/"` or `"\"` |

## `os.path_list_separator <: string`

Returns the path list separator of the current operating system.
`:` for UNIX-like machines and `;` for Windows machines. It returns a string.

Usage:

| example | equals |
|:-|:-|
| `//os.path_list_separator` | `":"` or `";"` |

## `os.stdin <: array_of_bytes`

`//os.stdin` holds all the bytes read from stdin. The bytes are read when
`//os.stdin` is accessed for the first time.

| example | equals |
|:-|:-|
| `//os.stdin` | `{ |@, @byte| ... }` |
