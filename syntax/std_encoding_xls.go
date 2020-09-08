package syntax

import (
	"bytes"
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/go-errors/errors"
	"github.com/iancoleman/strcase"
	"strings"

	"github.com/arr-ai/arrai/rel"
)

func stdEncodingXls() rel.Attr {
	return rel.NewTupleAttr(
		"xls",
		createFunc2Attr("decode", func(_ context.Context, x, i rel.Value) (rel.Value, error) {
			var bs []byte
			switch b := x.(type) {
			case rel.String:
				bs = []byte(b.String())
			case rel.Bytes:
				bs = b.Bytes()
			default:
				return nil, errors.Errorf("first arg to xls.decode must be string or bytes, not %T", b)
			}

			switch iv := i.(type) {
			case rel.Number:
				ix, ok := iv.Int()
				if !ok {
					return nil, errors.Errorf("second arg to xls.decode must be integer, not %v", i)
				}
				return bytesXlsxToArrai(bs, ix)
			default:
				return nil, errors.Errorf("second arg to xls.decode must be integer, not %T", i)
			}
		}),
	)
}

func bytesXlsxToArrai(bs []byte, sheetIndex int) (rel.Value, error) {
	f, err := excelize.OpenReader(bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}
	sheets := f.GetSheetList()
	if len(sheets) < sheetIndex+1 {
		return nil, fmt.Errorf("no sheet at index %d", sheetIndex)
	}
	rows, err := f.GetRows(sheets[sheetIndex])
	if err != nil {
		return nil, err
	}

	merges, err := f.GetMergeCells(sheets[sheetIndex])
	if err != nil {
		return nil, err
	}
	getValue := func(sheet [][]string, row, col int) string {
		if row >= len(sheet) {
			return ""
		}
		r := sheet[row]
		if col >= len(r) {
			return ""
		}
		cell := r[col]
		if cell != "" {
			return cell
		}
		for _, m := range merges {
			scol, srow, _ := excelize.CellNameToCoordinates(m.GetStartAxis())
			ecol, erow, _ := excelize.CellNameToCoordinates(m.GetEndAxis())
			if col+1 >= scol && col+1 <= ecol && row >= srow && row <= erow {
				return m.GetCellValue()
			}
		}
		return ""
	}

	cols := []string{}
	for _, cell := range rows[0] {
		// Commas in names are difficult to use in arr.ai (e.g. in nest attr lists).
		cols = append(cols, strings.ReplaceAll(strcase.ToSnake(cell), ",", ""))
	}
	rowTuples := []rel.Value{}
	for i := 1; i < len(rows); i++ {
		attrs := []rel.Attr{rel.NewIntAttr("@row", i)}
		for j, name := range cols {
			if name == "" {
				continue
			}
			val := getValue(rows, i, j)
			attrs = append(attrs, rel.NewStringAttr(name, []rune(val)))
		}
		if len(attrs) > 0 {
			rowTuples = append(rowTuples, rel.NewTuple(attrs...))
		}
	}
	return rel.NewSet(rowTuples...), nil
}
