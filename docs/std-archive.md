# archive

The `archive` library contains helper functions that related to outputting data into
certain formats of archive.

## `//archive.tar <: tuple`

`tar` is a tuple of functions related to the `tar` format.

### `//archive.tar.tar(data <: dict) <: array_of_bytes`

`tar` encodes `data`, representing a directory tree and its files, as the bytes of a tar archive.

`data` must be a dictionary with all keys of type `string`

Usage:

| example | equals |
|:-|:-|
| `//archive.tar.tar({"lightsaber": "lightsaber noises"})` | `lightsaber0000600000000000000000000000002100000000000011235 0ustar0000000000000000lightsaber noises` |

## `//archive.zip <: tuple`

`zip` is a tuple of functions related to the `zip` format.

### `//archive.zip.zip(data <: dict) <: array_of_bytes`

`zip` encodes `data`, representing a directory tree and its files, as the bytes of a ZIP archive.

`data` must be a dictionary with all keys of type `string`

Usage:

| example | equals |
|:-|:-|
| `//archive.zip.zip({"sidious" : "so it is treason then"})` | `sidious*�W�,Q�,V()JM,��S(�H���P˘l(˘l(sidiousPK5P` |
