# //.os library

## args

`//.os.args` returns an array of strings representing the arguments provided on
the command line when running an arrai program. `//.os.args(0)` is the path of
the arr.ai program that was invoked from the command line.

## cwd

`//.os.cwd` returns a string representing the current user directory.

## file

`//.os.file()` is a function that returns the content of a file represented by a
filepath in the form of an array of bytes.

## get_env

`//.os.get_env()` is a function that returns the environment variable that
corresponds to the provided key in the form of a string.

## path_separator

`//.os.path_separator` returns the path separator of the current operating
system. `/` for UNIX-like machine and `\` for Windows machine. It returns a
string.

## path_list_separator

`//.os.path_list_separator` returns the path list separator of the current
operating system. `:` for UNIX-like machine and `;` for Windows machine. It
returns a string.

## stdin

`//.os.stdin` holds all the bytes read from stdin. The bytes are read when
`//.os.stdin` is accessed for the first time.
