# Writing tests

Arr.ai has a rudimentary testing library. The API is documented
[here](../std-test.md).

Writing tests is as simple as calling `//test.suite({...})` with a set of tests:

```arrai
@> //test.suite({
 >     //test.assert.equal(42)(6 * 7),
 >     //test.assert.unequal(42)(6 * 9),
 > })
@> //test.suite({
 >     //test.assert.false({1, 2, 3} where . < 2)
 > })
```
