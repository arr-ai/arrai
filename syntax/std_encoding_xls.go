package syntax

import (
	"bytes"
	"context"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/go-errors/errors"
	"github.com/iancoleman/strcase"

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
		return nil, errors.Errorf("no sheet at index %d", sheetIndex)
	}
	rows, err := f.GetRows(sheets[sheetIndex])
	if err != nil {
		return nil, err
	}

	merges, err := f.GetMergeCells(sheets[sheetIndex])
	if err != nil {
		return nil, err
	}
	getValue := func(sheet [][]string, row, col int) (string, error) {
		if row >= len(sheet) {
			return "", nil
		}
		r := sheet[row]
		if col >= len(r) {
			return "", nil
		}
		cell := r[col]
		if cell != "" {
			return cell, nil
		}
		for _, m := range merges {
			scol, srow, err := excelize.CellNameToCoordinates(m.GetStartAxis())
			if err != nil {
				return "", err
			}
			ecol, erow, err := excelize.CellNameToCoordinates(m.GetEndAxis())
			if err != nil {
				return "", err
			}
			if col+1 >= scol && col+1 <= ecol && row >= srow && row <= erow {
				return m.GetCellValue(), nil
			}
		}
		return "", nil
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
			val, err := getValue(rows, i, j)
			if err != nil {
				return nil, err
			}
			attrs = append(attrs, rel.NewStringAttr(name, []rune(val)))
		}
		if len(attrs) > 0 {
			rowTuples = append(rowTuples, rel.NewTuple(attrs...))
		}
	}
	return rel.NewSet(rowTuples...), nil
}
