---
id: import
title: Importing external code and data
---

It is possible to import external arr.ai scripts. Every imported arr.ai script
must be a valid arr.ai expression. There are multiple ways to import arr.ai scripts.

### Package Imports

Package imports allow for importing files relative to a "root". The root for a script is the location of the nearest `go.mod` file, in either the same directory as the importing script or an ancestor directory.

You may define multiple `go.mod` files in the same repository to create multiple import roots for different parts of the codebase.

Example:

```arrai
//{/path/to/script}
```

With the following directories structure:

```bash
project
├── go.mod
├── nested
│   ├── child1.arrai
│   ├── go.mod
│   └── path
│       └── to
│           └── grandchildren.arrai
├── parent.arrai
└── path
    └── to
        └── child2.arrai
```

In that scenario, the root path of `parent.arrai` and `path/to/child2.arrai` is
`project/` and the root path for `nested/child1.arrai`and
`nested/path/to/grandchildren.arrai` is `project/nested/`. This means that any
Package imports in `parent.arrai` and `child2.arrai` will be relative to
`project/` and Package imports in `child1.arrai` and `grandchildren.arrai` will
be relative to `project/nested`.

### Relative Imports

It is possible to import a script relative to the importer's directory.

A relative import path begins with a `.`. For example:

```arrai
//{./path/to/script}
```

It is not possible to import anything from the parent directories.

```arrai
# this is not allowed
//{../script}
//{./../../script}
```

However, sometimes importing above the parent directory might be necessary for
testing locally. It is possible to work around this by using the Package import
and `go.mod`. Since `go.mod` also handles dependencies of arr.ai scripts, it is
possible to use `replace` to a dependency so that it points to a custom URL or
a filepath.

Example:

`go.mod`

```go
go 1.14

require some/dependency/path v1.2.3

replace some/dependency/path v1.2.3 => ../some/custom/dependency
```

**Note**: There are plans to deprecate the use of Go dependencies and `go.mod` for managing arr.ai imports. Whatever the new solution is will have an equivalent mechanism for `replace`.

### Remote Imports

Remote imports allow importing of an arr.ai script that is hosted in a GitHub
repository. This can be done by the following:

```arrai
//{github.com/username/repo_name/path/to/script}
```

TODO: versioning


### HTTP Imports

HTTP imports allow importing of a remotely-hosted arr.ai script given its URL. The
server must return an arr.ai script as a text. This can be done by the
following:

```arrai
//{http://www.somewebsite.com/path/to/arrai}

or

//{https://www.somewebsite.com/path/to/arrai}
```

### Non-arr.ai import

arr.ai supports importing of non-arr.ai files. The currently supported files are
the following file formats:

- JSON
- YAML

This can be done by adding the file extension to explicitly state file format
like the following:

```arrai
//{./path/to/file.json}
```

For importing arr.ai scripts, it is not necessary to add the `.arrai`
extension as it is always assumed that the files are arr.ai scripts, unless
explicitly stated.
