package rel

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExprIndexedSequenceMap(t *testing.T) {
	t.Parallel()

	AssertExprsEvalToSameValue(t,
		NewArray(Number(1), Number(3), Number(5)),
		NewIndexedSequenceMapExpr(
			NewArray(Number(1), Number(2), Number(3)),
			NewFunction(
				"i",
				NewFunction(
					"n",
					NewAddExpr(
						NewIdentExpr("i"),
						NewIdentExpr("n"),
					),
				),
			),
		),
	)

	AssertExprsEvalToSameValue(t,
		NewDict(false,
			NewDictEntryTuple(
				NewString([]rune("stuff")),
				NewTuple(
					NewAttr("key", NewString([]rune("stuff"))),
					NewAttr("val", NewString([]rune("random"))),
				),
			),
			NewDictEntryTuple(
				NewString([]rune("ten")),
				NewTuple(
					NewAttr("key", NewString([]rune("ten"))),
					NewAttr("val", NewNumber(10)),
				),
			),
			NewDictEntryTuple(
				NewNumber(3),
				NewTuple(
					NewAttr("key", NewNumber(3)),
					NewAttr("val", NewSet(NewNumber(2))),
				),
			),
		),
		NewIndexedSequenceMapExpr(
			NewDict(false,
				NewDictEntryTuple(NewString([]rune("stuff")), NewString([]rune("random"))),
				NewDictEntryTuple(NewString([]rune("ten")), NewNumber(10)),
				NewDictEntryTuple(NewNumber(3), NewSet(NewNumber(2))),
			),
			NewFunction(
				"i",
				NewFunction(
					"n",
					NewTupleExpr(
						MustNewAttrExpr("key", NewIdentExpr("i")),
						MustNewAttrExpr("val", NewIdentExpr("n")),
					),
				),
			),
		),
	)
}

func TestExprIndexedSequenceMapFail(t *testing.T) {
	t.Parallel()

	val := NewSet(EmptyTuple.With("a", Number(1)).With("b", Number(2)))
	_, err := NewIndexedSequenceMapExpr(
		val,
		NewFunction(
			"i",
			NewFunction(
				"n",
				NewAddExpr(
					NewIdentExpr("i"),
					NewIdentExpr("n"),
				),
			),
		),
	).Eval(EmptyScope)

	assert.EqualError(t, err, fmt.Sprintf("=> not applicable to unindexed type %v", val))
}
