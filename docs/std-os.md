# os

The `os` contains functions that are related to the operating system.

## `//os.args <: array_of_strings`

`args` contains an array of strings representing the arguments provided on
the command line when running an arrai program. `//os.args(0)` is the path of
the arr.ai program that was invoked from the command line.

Usage:

| example | equals |
|:-|:-|
|`//os.args` | `["arg0", "arg1", "arg2", ...]` |

## `//os.cwd <: string`

`cwd` contains a string representing the current user directory.

## `//os.exists(filepath <: string) <: boolean`

`exists` returns `true` if a file or directory exists at `filepath`, or `false` otherwise.

## `//os.file(filepath <: string) <: array_of_bytes`

`file` returns the content of a file located at `filepath` in the form of an array of bytes.

Usage:

| example | equals |
|:-|:-|
|`//os.file("path/to/file")` | `<<...>>` |
|`//os.file("/absolute/path/to/file")` | `<<...>>` |

## `//os.tree(dirPath <: string) <: set`

`tree` returns a set of `stat` details about all directories and files in the tree rooted at `dirPath` (including `dirPath` itself.

If `dirPath` is a file, `tree` returns a set containing details for just that file.

The `path` attribute is relative to the current working directory, *not* `dirPath`.

## `//os.get_env(key <: string) <: string`

`get_env` returns the environment variable that corresponds to `key` in the form of a `string`.

Usage:

| example | equals |
|:-|:-|
| `//os.get_env("KEY")` | `"string_value"` |

## `//os.isatty(fileDescriptor <: int) <: bool`

`isatty` returns `true` if the `fileDescriptor` (0 for `stdin`, 1 for `stdout`) is a terminal, or `false` otherwise (e.g. if input is piped in).

Usage:

| example | equals |
|:-|:-|
| `arrai eval "//os.isatty(0)"` | `true` |
| `echo "" | arrai eval "//os.isatty(0)"` | `false` |

## `//os.path_separator <: string`

`path_separator` contains the path separator of the current operating system.
`/` for UNIX-like machine and `\` for Windows machine.

Usage:

| example | equals |
|:-|:-|
| `//os.path_separator` | `"/"` or `"\"` |

## `//os.path_list_separator <: string`

`path_list_separator` contains the path list separator of the current operating system.
`:` for UNIX-like machines and `;` for Windows machines.

Usage:

| example | equals |
|:-|:-|
| `//os.path_list_separator` | `":"` or `";"` |

## `//os.stdin <: array_of_bytes`

`stdin` holds all the bytes read from stdin. The bytes are read when
`//os.stdin` is accessed for the first time.

| example | equals |
|:-|:-|
| `//os.stdin` | `<<...>>` |
