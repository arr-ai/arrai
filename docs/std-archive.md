# archive

`archive` library contains helper functions that related to outputting data into
certain formats of archive.

## `tar <: tuple`

`tar` returns a tuple that contains functions related to the format `tar`. The
attributes are the following:

### `tar(data <: dict) <: array_of_bytes`

It takes the dictionary `data` and returns an array of bytes representing `data`
which has been compressed to the `tar` format.

**Important note**: `data` itself has to be a dictionary whose keys are all
`string`.

Usage:

| example | equals |
|:-|:-|
| `//archive.tar.tar({"lightsaber": "lightsaber noises"})` | `lightsaber0000600000000000000000000000002100000000000011235 0ustar0000000000000000lightsaber noises` |

## `zip <: tuple`

`zip` returns a tuple that contains functions related to the format `zip`. The
attributes are the following:

### `zip(data <: dict) <: array_of_bytes`

It takes the dictionary `data` and returns an array of bytes representing `data`
which has been compressed to the `zip` format.

**Important note**: `data` itself has to be a dictionary whose keys are all
`string`.

Usage:

| example | equals |
|:-|:-|
| `//archive.zip.zip({"sidious" : "so it is treason then"})` | `sidious*�W�,Q�,V()JM,��S(�H���P˘l(˘l(sidiousPK5P` |
