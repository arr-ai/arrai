package syntax

import (
	"testing"
)

func TestCsvDecode(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`[['a', 'b', 'c'], ['1', '2', '3']]`,
		`//encoding.csv.decode(<<'a,b,c\n1,2,3\n'>>)`)
}

func TestCsvDecode_comma(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`[['a', 'b', 'c'], ['1', '2', '3']]`,
		`//encoding.csv.decoder((comma: %|))(<<'a|b|c\n1|2|3\n'>>)`)
}

func TestCsvDecode_comment(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`[['a', 'b', 'c'], ['1', '2', '3']]`,
		`//encoding.csv.decoder((comment: %#))(<<'a,b,c\n#comment\n1,2,3\n'>>)`)
}

func TestCsvDecode_trimLeadingSpace(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`[['a', 'b', 'c'], ['1', '2', '3']]`,
		`//encoding.csv.decoder((trimLeadingSpace: true))(<<' a,b, c\n 1,2,3\n'>>)`)
}

func TestCsvDecode_fieldsPerRecord_negative_variableLengthFieldsPermitted(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`[['a', 'b', 'c'], ['1', '2']]`,
		`//encoding.csv.decoder((fieldsPerRecord: -1))(<<'a,b,c\n1,2\n'>>)`)
}

func TestCsvDecode_fieldsPerRecord_zero_linesMatchFirstRow(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`[['a', 'b', 'c'], ['1', '2', '3']]`,
		`//encoding.csv.decoder((fieldsPerRecord: 0))(<<'a,b,c\n1,2,3\n'>>)`)
}

func TestCsvDecode_fieldsPerRecord_zero_linesDoNotMatchFirstRow(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t,
		"record on line 2: wrong number of fields",
		`//encoding.csv.decoder((fieldsPerRecord: 0))(<<'a,b,c\n1,2\n'>>)`)
}

func TestCsvDecode_fieldsPerRecord_positive_linesMatchCount(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`[['a', 'b', 'c'], ['1', '2', '3']]`,
		`//encoding.csv.decoder((fieldsPerRecord: 3))(<<'a,b,c\n1,2,3\n'>>)`)
}

func TestCsvDecode_fieldsPerRecord_positive_linesDoNotMatchCount(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t,
		"record on line 1: wrong number of fields",
		`//encoding.csv.decoder((fieldsPerRecord: 2))(<<'a,b,c\n1,2,3\n'>>)`)
}

func TestCsvDecode_lazyQuotes(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`[['a', 'b', 'c'], ['1', '2', '3']]`,
		`//encoding.csv.decoder((lazyQuotes: true))(<<'"a",b,c\n1,"2",3\n'>>)`)
}

func TestCsvDecode_string(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`[['a', 'b', 'c']]`,
		`//encoding.csv.decode('a,b,c')`)
}

func TestCsvDecode_error(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t,
		"first arg to csv.decode must be string or bytes, not tuple",
		`//encoding.csv.decode(())`)

	AssertCodeErrors(t,
		"first arg to csv.decoder must be tuple, not string",
		`//encoding.csv.decoder('csv')('csv')`)

	AssertCodeErrors(t,
		"comma config param to csv.decoder must be integer, not tuple",
		`//encoding.csv.decoder((comma: ()))('csv')`)

	AssertCodeErrors(t,
		"comment config param to csv.decoder must be integer, not tuple",
		`//encoding.csv.decoder((comment: ()))('csv')`)

	AssertCodeErrors(t,
		"first arg to csv.decoder payload must be string or bytes, not tuple",
		`//encoding.csv.decoder(())(())`)

	AssertCodeErrors(t,
		"record on line 2: wrong number of fields",
		`//encoding.csv.decode('a,b,c\n# comment')`)
}

func TestCsvEncode(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`<<'a,b,c\n1,2,3\n'>>`,
		`//encoding.csv.encode([['a', 'b', 'c'], ['1', '2', '3']])`)
}

func TestCsvEncode_comma(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`<<'a|b|c\n1|2|3\n'>>`,
		`//encoding.csv.encoder((comma: %|))([['a', 'b', 'c'], ['1', '2', '3']])`)
}

func TestCsvEncode_crlf(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`<<'a,b,c\r\n1,2,3\r\n'>>`,
		`//encoding.csv.encoder((crlf: true))([['a', 'b', 'c'], ['1', '2', '3']])`)
}

func TestCsvEncode_error(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t,
		"first arg to csv.encode must be array, not tuple",
		`//encoding.csv.encode(())`)

	AssertCodeErrors(t,
		"first arg to csv.encoder must be tuple, not string",
		`//encoding.csv.encoder('csv')('csv')`)

	AssertCodeErrors(t,
		"comma config param to csv.encoder must be integer, not tuple",
		`//encoding.csv.encoder((comma: ()))([])`)

	AssertCodeErrors(t,
		"first arg to csv.encoder payload must be array, not tuple",
		`//encoding.csv.encoder(())(())`)

	AssertCodeErrors(t,
		"record 1 must be array, not number",
		`//encoding.csv.encode([['a'], 1])`)

	AssertCodeErrors(t,
		"value 0 of record 1 must be string, not num",
		`//encoding.csv.encode([['a'], [1]])`)
}
