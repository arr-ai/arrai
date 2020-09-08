package syntax

import (
	"testing"
)

func TestXlsDecode(t *testing.T) {
	t.Parallel()

	// foo.xlsx tests multiple edge cases, such as empty cells, rows, columns, names, commas, etc.
	AssertCodesEvalToSameValue(t,
		`{
			|@row, customer_id, account_type, account_balance|
			(1, 'foo', 'Checking', '100'), 
			(3, 'bar', 'Checking', '200'),
			(4, 'bar', 'Savings', '300'),
		}`,
		`//encoding.xls.decode(//os.file('../examples/xlsx/foo.xlsx'), 0)`)
}

func TestXlsDecode_error(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "first arg to xls.decode must be string or bytes, not *rel.GenericTuple", `//encoding.xls.decode((), 0)`)

	AssertCodeErrors(t, "second arg to xls.decode must be integer, not *rel.GenericTuple", `//encoding.xls.decode(<<"1">>, ())`)
	AssertCodeErrors(t, "second arg to xls.decode must be integer, not 1.5", `//encoding.xls.decode(<<"1">>, 1.5)`)
}
