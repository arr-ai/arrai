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
		`//encoding.xlsx.decodeToRelation((sheet: 0, headRow: 1), //os.file('../examples/xlsx/foo.xlsx'))`)
}

func TestXlsxDecode_error(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t,
		"second arg to xlsx.decodeToRelation must be string or bytes, not *rel.GenericTuple",
		`//encoding.xlsx.decodeToRelation((sheet: 0, headRow: 1), ())`)

	AssertCodeErrors(t,
		"sheet config param to xlsx.decodeToRelation must be integer, not *rel.GenericTuple",
		`//encoding.xlsx.decodeToRelation((sheet: ()), <<"1">>)`)
	AssertCodeErrors(t,
		"sheet config param to xlsx.decodeToRelation must be integer, not 1.5",
		`//encoding.xlsx.decodeToRelation((sheet: 1.5), <<"1">>)`)

	AssertCodeErrors(t,
		"headRow config param to xlsx.decodeToRelation must be integer, not *rel.GenericTuple",
		`//encoding.xlsx.decodeToRelation((headRow: ()), <<"1">>)`)
	AssertCodeErrors(t,
		"headRow config param to xlsx.decodeToRelation must be integer, not 1.5",
		`//encoding.xlsx.decodeToRelation((headRow: 1.5), <<"1">>)`)
}
