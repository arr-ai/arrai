# Import

These examples demonstrate how importing works in arr.ai

There are two ways to import

## Relative Import

Relative imports are relative to the source file that imports the other file.
The path will have the prefix `./`.

## Module Import

Module imports are relative to the root of the project. The root of the project
is determined by the file `go.mod`. The file itself can be located in the same
directory of the source script or the parent directory. This is meant to avoid
importing outside of the root.

If there is no `go.mod`, the directory of the source file is considered the root
of the project.

The import path of this import has the prefix `/`.
