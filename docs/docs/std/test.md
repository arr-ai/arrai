The `test` library contains functions that are used to write tests.

## `//test.suite(tests <: set)`

`suite` checks the outcome of each test and produces a report of failed tests.

Usage:

| example | equals |
|:-|:-|
| `//str.suite(//test.assert.equal(42)(6 * 7))` | `{}` |
| `//str.suite(//test.assert.equal(42)(6 * 9))` | Failure |

## `//test.assert <: tuple`

`assert` has a range of assertions.

### `//test.assert.equal(expected <: any, actual <: any)`

`equal` checks that `expected = actual`, otherwise it triggers a failure report
in the containing `//test.suite({...})` call.

Usage:

| example | equals |
|:-|:-|
| `//str.suite(//test.assert.equal(42)(6 * 7))` | `{}` |

### `//test.assert.false(value <: any)`

`true` checks that `cond (value: true) = false`, otherwise it triggers a failure
report in the containing `//test.suite({...})` call.

Usage:

| example | equals |
|:-|:-|
| `//str.suite(//test.assert.false(1 > 2))` | `{}` |

### `//test.assert.true(value <: any)`

`true` checks that `cond (value: true) = true`, otherwise it triggers a failure
report in the containing `//test.suite({...})` call.

Usage:

| example | equals |
|:-|:-|
| `//str.suite(//test.assert.true(1 < 2))` | `{}` |

### `//test.assert.unequal(unexpected <: any, actual <: any)`

`unequal` checks that `expected != actual`, otherwise it triggers a failure report
in the containing `//test.suite({...})` call.

Usage:

| example | equals |
|:-|:-|
| `//str.suite(//test.assert.unequal(42)(6 * 9))` | `{}` |
