//TODO: the context here maybe need to be initialized with proper values like fs
package rel

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"

	"github.com/arr-ai/arrai/pkg/arraictx"
)

// intSet returns a new set from the given elements.
func intSet(elts ...interface{}) Set {
	result, err := NewSetFrom(elts...)
	if err != nil {
		panic(err)
	}
	return result
}

func EqualValues(expected, actual Value) bool {
	return expected == nil && actual == nil || expected.Equal(actual)
}

// AssertEqualValues asserts that the two values are Equal.
func AssertEqualValues(t *testing.T, expected, actual Value) bool {
	t.Helper()

	if !EqualValues(expected, actual) {
		return assert.Fail(t, "values not equal", "expected: %s\nactual:   %s", expected, actual)
	}
	return true
}

// RequireEqualValues requires that the two values are Equal.
func RequireEqualValues(t *testing.T, expected, actual Value) {
	t.Helper()

	if !AssertEqualValues(t, expected, actual) {
		t.FailNow()
	}
}

// AssertExprsEvalToSameValue asserts that the exprs evaluate to the same value.
func AssertExprsEvalToSameValue(t *testing.T, expected, expr Expr) bool {
	t.Helper()

	ctx := arraictx.InitRunCtx(context.Background())
	expectedValue, err := expected.Eval(ctx, EmptyScope)
	if !assert.NoError(t, err, "evaluating expected: %s", expected) {
		return false
	}
	value, err := expr.Eval(ctx, EmptyScope)
	if !assert.NoError(t, err, "evaluating expr: %s", expr) {
		return false
	}
	return EqualValues(expectedValue, value) ||
		assert.Failf(t, "values not equal",
			"\nexpected: %v\nactual:   %v\nexpr:     %v",
			expectedValue, value, expr)
}

// RequireExprsEvalToSameValue requires that the exprs evaluate to the same
// value.
func RequireExprsEvalToSameValue(t *testing.T, expected, expr Expr) {
	t.Helper()

	if !AssertExprsEvalToSameValue(t, expected, expr) {
		t.FailNow()
	}
}

// AssertExprEvalsToType asserts that the exprs evaluate to the same value.
func AssertExprEvalsToType(t *testing.T, expected interface{}, expr Expr) bool {
	t.Helper()

	value, err := expr.Eval(arraictx.InitRunCtx(context.Background()), EmptyScope)
	if !assert.NoError(t, err, "evaluating expr: %s", expr) {
		return false
	}
	if reflect.TypeOf(expected) != reflect.TypeOf(value) {
		t.Logf("\nexpected: %T\nvalue:    %v\nexpr:     %v", expected, value, expr)
		return false
	}
	return true
}

// AssertExprErrors asserts that the expr returns an error when evaluated.
func AssertExprErrors(t *testing.T, expr Expr) bool {
	t.Helper()

	_, err := expr.Eval(arraictx.InitRunCtx(context.Background()), EmptyScope)
	return assert.Error(t, err)
}

// AssertExprErrorEquals asserts that the expr returns an error with the given message when evaluated.
func AssertExprErrorEquals(t *testing.T, expr Expr, msg string) bool {
	t.Helper()

	_, err := expr.Eval(arraictx.InitRunCtx(context.Background()), EmptyScope)

	return assert.EqualError(t, err, WrapContextErr(errors.Errorf(msg), expr, EmptyScope).Error())
}

// AssertExprPanics asserts that the expr panics when evaluated.
func AssertExprPanics(t *testing.T, expr Expr) bool {
	t.Helper()

	return assert.Panics(t, func() { expr.Eval(arraictx.InitRunCtx(context.Background()), EmptyScope) }) //nolint:errcheck
}
