# Standard Library

The scripts `stdlib-safe.arrai` and `stdlib-unsafe.arrai` provide a way to write some arr.ai standard libraries functions in arr.ai. Each script is a function that takes the standard library written in golang as a tuple. The scripts can add more fields to add more standard library functions in arr.ai. If the function is an expression that takes in inputs and returns a value, it is considered safe and should be added to `stdlib-safe.arrai`. If the function interacts with the outside world, for example by creating or modifying files, it is considered unsafe and should be added to `stdlib-unsafe.arrai`. A quick example:

```arrai
\stdlib
    stdlib +> (str: stdlib.str +> (trimPrefix: \string \prefix  (string &~ prefix) rank (:.@)))
```

It is important to note that the expressions inside the scripts can not use the standard library expressions `//` as it will create stack overflow.

It is also important that any asset written arr.ai that uses the standard library must be loaded lazily during compilation as it can also create stack overflow. You can do this by just returning the created expression instead of the value evaluated from the expression.
