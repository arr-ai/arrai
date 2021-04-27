package rel

import (
	"context"
	"testing"

	"github.com/arr-ai/arrai/pkg/arraictx"
	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
)

func TestOffsetExprArray(t *testing.T) {
	t.Parallel()

	AssertExprsEvalToSameValue(t,
		NewOffsetArray(
			5,
			NewNumber(float64(1)),
			NewNumber(float64(2)),
			NewNumber(float64(3)),
		),
		NewOffsetExpr(
			*parser.NewScanner(""),
			NewNumber(float64(5)),
			NewArray(
				NewNumber(float64(1)),
				NewNumber(float64(2)),
				NewNumber(float64(3)),
			),
		),
	)

	AssertExprsEvalToSameValue(t,
		NewOffsetArray(
			-3,
			NewNumber(float64(1)),
			NewNumber(float64(2)),
			NewNumber(float64(3)),
		),
		NewOffsetExpr(
			*parser.NewScanner(""),
			NewNumber(float64(-5)),
			NewOffsetArray(
				2,
				NewNumber(float64(1)),
				NewNumber(float64(2)),
				NewNumber(float64(3)),
			),
		),
	)

	AssertExprsEvalToSameValue(t,
		NewArray(
			NewNumber(float64(1)),
			NewNumber(float64(2)),
			NewNumber(float64(3)),
		),
		NewOffsetExpr(
			*parser.NewScanner(""),
			NewNumber(float64(0)),
			NewArray(
				NewNumber(float64(1)),
				NewNumber(float64(2)),
				NewNumber(float64(3)),
			),
		),
	)
}

func TestOffsetExprBytes(t *testing.T) {
	t.Parallel()

	AssertExprsEvalToSameValue(t,
		NewOffsetBytes(
			[]byte("random string"),
			10,
		),
		NewOffsetExpr(
			*parser.NewScanner(""),
			NewNumber(float64(10)),
			NewBytes([]byte("random string")),
		),
	)

	AssertExprsEvalToSameValue(t,
		NewOffsetBytes([]byte("random string"), -12),
		NewOffsetExpr(
			*parser.NewScanner(""),
			NewNumber(float64(-6)),
			NewOffsetBytes([]byte("random string"), -6),
		),
	)

	AssertExprsEvalToSameValue(t,
		NewOffsetBytes([]byte("random string"), 2),
		NewOffsetExpr(
			*parser.NewScanner(""),
			NewNumber(float64(0)),
			NewOffsetBytes([]byte("random string"), 2),
		),
	)
}

func TestOffsetExprString(t *testing.T) {
	t.Parallel()

	AssertExprsEvalToSameValue(t,
		NewOffsetString(
			[]rune("random string"),
			3,
		),
		NewOffsetExpr(
			*parser.NewScanner(""),
			NewNumber(float64(3)),
			NewString([]rune("random string")),
		),
	)

	AssertExprsEvalToSameValue(t,
		NewOffsetString([]rune("random string"), -10),
		NewOffsetExpr(
			*parser.NewScanner(""),
			NewNumber(float64(-10)),
			NewString([]rune("random string")),
		),
	)

	AssertExprsEvalToSameValue(t,
		NewOffsetString([]rune("random string"), -2),
		NewOffsetExpr(
			*parser.NewScanner(""),
			NewNumber(float64(0)),
			NewOffsetString([]rune("random string"), -2),
		),
	)
}

func TestOffsetExprEmptySet(t *testing.T) {
	t.Parallel()

	AssertExprsEvalToSameValue(t,
		None,
		NewOffsetExpr(*parser.NewScanner(""), NewNumber(float64(-10)), None),
	)
}

func TestOffsetExprEvalFail(t *testing.T) {
	t.Parallel()

	// None in LHS instead of a Number
	_, err := NewOffsetExpr(
		*parser.NewScanner(""), None, None).Eval(arraictx.InitRunCtx(context.Background()), EmptyScope)
	expected := errors.Errorf("offset must be a number, not set").Error()
	assert.EqualError(t, errors.New(err.Error()[:len(expected)]), expected)

	// Randomg set in RHS instead of an Array
	_, err = NewOffsetExpr(
		*parser.NewScanner(""),
		Number(float64(0)),
		MustNewSet(Number(float64(0))),
	).Eval(arraictx.InitRunCtx(context.Background()), EmptyScope)
	expected = errors.Errorf("offset not applicable to set").Error()
	assert.EqualError(t, errors.New(err.Error()[:len(expected)]), expected)
}
