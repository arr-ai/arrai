package syntax

import (
	"testing"
)

func TestXlsDecode(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`{
			|@row, customer_id, account, balance|
			(1, 'foo', 'Checking', '100'), 
			(2, 'bar', 'Checking', '200'),
			(3, 'bar', 'Savings', '300'),
		}`,
		`//encoding.xls.decode(//os.file('../examples/xlsx/foo.xlsx'), 0)`)
}
