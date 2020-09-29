package syntax

import (
	"testing"
)

func TestXlsxDecode(t *testing.T) {
	t.Parallel()

	// foo.xlsx tests multiple edge cases, such as empty cells, rows, columns, names, commas, etc.
	AssertCodesEvalToSameValue(t,
		`{
			|@row, customer_id, account_type, account__balance|
			(2, 'foo', 'Checking', '100'), 
			(4, 'bar', 'Checking', '200'),
			(5, 'bar', 'Savings', '300'),
		}`,
		`//encoding.xlsx.decodeToRelation(//os.file('../examples/xlsx/foo.xlsx'), 0, 1)`)
}

func TestXlsxDecode_error(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t,
		"first arg to xlsx.decodeToRelation must be string or bytes, not *rel.GenericTuple",
		`//encoding.xlsx.decodeToRelation((), 0, 1)`)

	AssertCodeErrors(t,
		"second arg to xlsx.decodeToRelation must be integer, not *rel.GenericTuple",
		`//encoding.xlsx.decodeToRelation(<<"1">>, (), 1)`)
	AssertCodeErrors(t,
		"second arg to xlsx.decodeToRelation must be integer, not 1.5",
		`//encoding.xlsx.decodeToRelation(<<"1">>, 1.5, 1)`)

	AssertCodeErrors(t,
		"third arg to xlsx.decodeToRelation must be integer, not *rel.GenericTuple",
		`//encoding.xlsx.decodeToRelation(<<"1">>, 0, ())`)
	AssertCodeErrors(t,
		"third arg to xlsx.decodeToRelation must be integer, not 1.5",
		`//encoding.xlsx.decodeToRelation(<<"1">>, 0, 1.5)`)
}
