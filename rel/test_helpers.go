package rel

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// intSet returns a new set from the given elements.
func intSet(elts ...interface{}) Set {
	result, err := NewSetFrom(elts...)
	if err != nil {
		panic(err)
	}
	return result
}

func equalValues(expected, actual Value) bool {
	return expected == nil && actual == nil || expected.Equal(actual)
}

// AssertEqualValues asserts that the two values are Equal.
func AssertEqualValues(t *testing.T, expected, actual Value) bool {
	if !equalValues(expected, actual) {
		return assert.Fail(t, "values not equal", "expected: %s\nactual:   %s", expected, actual)
	}
	return true
}

// RequireEqualValues requires that the two values are Equal.
func RequireEqualValues(t *testing.T, expected, actual Value) {
	if !AssertEqualValues(t, expected, actual) {
		t.FailNow()
	}
}

// AssertExprsEvalToSameValue asserts that the exprs evaluate to the same value.
func AssertExprsEvalToSameValue(t *testing.T, expected, expr Expr) bool {
	expectedValue, err := expected.Eval(EmptyScope)
	if !assert.NoError(t, err, "evaluating expected: %s", expected) {
		return false
	}
	value, err := expr.Eval(EmptyScope)
	if !assert.NoError(t, err, "evaluating expr: %s", expr) {
		return false
	}
	return equalValues(expectedValue, value) ||
		assert.Failf(t, "values not equal",
			"\nexpected: %v\nactual:   %v\nexpr:     %v",
			Repr(expectedValue), Repr(value), expr)
}

// RequireExprsEvalToSameValue requires that the exprs evaluate to the same
// value.
func RequireExprsEvalToSameValue(t *testing.T, expected, expr Expr) {
	if !AssertExprsEvalToSameValue(t, expected, expr) {
		t.FailNow()
	}
}

// AssertExprEvalsToType asserts that the exprs evaluate to the same value.
func AssertExprEvalsToType(t *testing.T, expected interface{}, expr Expr) bool {
	value, err := expr.Eval(EmptyScope)
	if !assert.NoError(t, err, "evaluating expr: %s", expr) {
		return false
	}
	if reflect.TypeOf(expected) != reflect.TypeOf(value) {
		t.Logf("\nexpected: %T\nvalue:    %v\nexpr:     %v", expected, Repr(value), expr)
		return false
	}
	return true
}

// AssertExprPanics asserts that the expr panics when evaluates.
func AssertExprPanics(t *testing.T, expr Expr) bool {
	return assert.Panics(t, func() { expr.Eval(EmptyScope) }) //nolint:errcheck
}
