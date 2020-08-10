# Writing tests

Arr.ai is a functional language designed for the representation and transformation of data. As such, testing arr.ai code is different from testing more stateful, imperative code.

Arr.ai's approach to testing is for test files to produce a data structure within which all leaves are `true` (not just truthy). If any leaf is not `true`, the test is said to have failed.

## Example

Image we've written a function to calculate the final score of a game of ten pin bowling (as per [Uncle Bob's famous TDD kata][kata]).

`bowling.arrai`:

```arrai
let bowl = \game \pins
    game
;
```

`bowling_test.arrai`:

```arrai
let rollMany = \n \pins

;

{

}
```

  private void rollMany(int n, int pins) {
    for (int i = 0; i < n; i++)
      g.roll(pins);
  }

  ...
  public void testGutterGame() throws Exception {
    rollMany(20, 0);
    assertEquals(0, g.score());
  }

  public void testAllOnes() throws Exception {
    rollMany(20,1);
    assertEquals(20, g.score());
  }

  public void testOneSpare() throws Exception {
    rollSpare();
    g.roll(3);
    rollMany(17,0);
    assertEquals(16,g.score());
  }

  public void testOneStrike() throws Exception {
    rollStrike();
    g.roll(3);
    g.roll(4);
    rollMany(16, 0);
    assertEquals(24, g.score());
  }

  public void testPerfectGame() throws Exception {
    rollMany(12,10);
    assertEquals(300, g.score());
  }

  private void rollStrike() {
    g.roll(10);
  }

  private void rollSpare() {
    g.roll(5);
    g.roll(5);
  }
}



[kata]: http://butunclebob.com/ArticleS.UncleBob.TheBowlingGameKata
