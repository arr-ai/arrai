# Import

These examples demonstrate how importing works in arr.ai

There are two ways to import.

## Relative Import

Relative imports are relative to the source file that imports the other file.
The path will have the prefix `./`.

## Module Import

Module imports are relative to the root of the module. The root of the module
is determined by the file `go.mod`. The file itself can be located in the same
directory of the source script or the parent directory. This is meant to avoid
importing outside of the root.

If no module root is found, module-relative imports will fail.

The import path of this import has the prefix `/`.
