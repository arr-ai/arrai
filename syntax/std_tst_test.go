package syntax

import "testing"

func TestStdTest(t *testing.T) {
	AssertCodesEvalToSameValue(t, `{}`, `//test.suite({})`)
	AssertCodesEvalToSameValue(t, `{false}`, `//test.suite({//test.assert.equal(42)(6 * 7)})`)
	AssertCodeErrors(t, `//test.suite({//test.assert.equal(42)(6 * 9)})`, "not equal\nexpected: 42\nactual:   54")
	AssertCodesEvalToSameValue(t, `{false}`, `
		//test.suite({
			//test.assert.equal(42)(6 * 7),
			//test.assert.unequal(42)(6 * 9),
		})`,
	)

	AssertCodesEvalToSameValue(t, `{false}`, `//test.suite({//test.assert.false(0)})`)
	AssertCodesEvalToSameValue(t, `{false}`, `//test.suite({//test.assert.false({})})`)
	AssertCodesEvalToSameValue(t, `{false}`, `//test.suite({//test.assert.false({1, 2, 3} where . > 10)})`)

	AssertCodesEvalToSameValue(t, `{false}`, `//test.suite({//test.assert.true(1)})`)
	AssertCodesEvalToSameValue(t, `{false}`, `//test.suite({//test.assert.true({1})})`)
	AssertCodesEvalToSameValue(t, `{false}`, `//test.suite({//test.assert.true({1, 2, 3} where . < 2)})`)
}
