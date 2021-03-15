# Standard Library

The script `stdlib.arrai` provides a way to write some arr.ai standard libraries functions in arr.ai. The script itself is a function that takes the standard library written in golang as a tuple. The script then can add more fields to add more standard library functions in arr.ai. A quick example:

```arrai
\stdlib
    stdlib +> (str: stdlib.str +> (trimPrefix: \string \prefix  (string &~ prefix) rank (:.@)))
```

It is important to note that the expressions inside the `stdlib.arrai` can not use the standard library expressions `//` as it will create stack overflow.

It is also important that any asset written arr.ai that uses the standard library must be loaded lazily during compilation as it can also create stack overflow. You can do this by just returning the created expression instead of the value evaluated from the expression.
