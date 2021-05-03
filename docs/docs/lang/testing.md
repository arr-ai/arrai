---
id: testing
title: Writing tests
---

import useBaseUrl from '@docusaurus/useBaseUrl';

Arr.ai is a functional language designed for the representation and transformation of data. As such, testing arr.ai code is different from testing more stateful, imperative code.

Arr.ai's approach to testing is for test files to produce a data structure within which all leaves are the literal `true` or `false` to indicate if a test has passed or failed. Any other value is considered an invalid test, and will fail the test run.

```arrai
(
    passingTest: //math.cos(//math.pi) = -1,
    failingTest: 5 > 10,
    invalidTest: 7,
)
```
A test run against the above will produce the following output:

<img alt="Test run output" src={useBaseUrl('img/test_output1.png')} />

### Test files and directories

Arr.ai test files follow a similar convention to Go, and must end with `_test.arrai`. Files without that suffix will not be treated as test files. There is no such requirement for directories, and all directories will be searched, except for hidden directories (prefixed with `.`).

### Command line
You can start a test run in one of three ways:
1. By not specifying a target. This will run all test files in the current working directory and all subdirectories recursively.
2. By specifying a directory target. This will run all test files in the specified directory, recursively.
3. By specifying a specific file as the target. This will run all tests in that file.

Examples of the three methods are given below, respectively.

```shell
$ arrai test
$ arrai test unit_tests/math
$ arrai test unit_tests/math/subtraction.arrai
```

### Test Containers

Test can be organised inside _test containers_ which logically group tests together. Named test containers, which require naming each test, are expressed as tuples and dictionaries. Unnamed test containers, which don't require naming tests are expressed using arrays. Sets are not allowed as test containers, since they can only contain one true test and one false test, and will be reported as invalid tests.

```arrai
(
    multiplication:
    [
        1 * 1 = 1,
        1 * 2 = 2,
        2 * 2 = 4,
        2 * 3 = 6
    ],
    division:
    {
        "positiveInfinity": 5/0 > 99999,
        3.141: //math.pi / 3 > 1,
    }
)
```

Which will result in:

<img alt="Test run output" src={useBaseUrl('img/test_output2.png')} />

### Exit code

If the test run concludes with any tests that have failed or are invalid, the test run will fail, and the exit code will be non-zero.

If evaluating the test files causes an error, be it a parsing error or runtime error, the test execution will stop immediately, the error will be printed, no test results will be reported, and the exit code will be non-zero. This includes `//testing.assert.*` functions that fail the assertion, and are thus deprecated. 

### Generative tests

If you have repetitive tests that differ only in data, you may generate a collection of results. For example:

```arrai
let isEven = \i i % 2 = 0;
(
    even: [-2, 0, 2, 10, 100] >> isEven(.),
    odd: [-1, 1, 3, 11, 101] >> !isEven(.)
)
```

Which will produce the output:

<img alt="Test run output" src={useBaseUrl('img/test_output3.png')} />


### Roadmap

- Replace leaf comparison expression with detailed comparisons to show what the actual vs expected was (plus a diff), instead of just true vs false. This is similar to what the `//testing.assert.*` functions provided, but without causing an error.
- Support ignoring individual tests, test containers, test files and directories (using a `__` prefix, or a `.testignore` file). Ignored tests will still be run and reported, but they don't affect the test run outcome and exit code.
- Run each test in isolation, so a failure in one test doesn't cause the test run to fail. Measure runtime of individual tests.
