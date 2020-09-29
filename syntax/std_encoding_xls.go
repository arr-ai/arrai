package syntax

import (
	"bytes"
	"context"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/go-errors/errors"
	"github.com/iancoleman/strcase"
	"regexp"

	"github.com/arr-ai/arrai/rel"
)

func stdEncodingXlsx() rel.Attr {
	return rel.NewTupleAttr(
		"xlsx",
		createFunc3Attr("decodeToRelation", func(_ context.Context, x, i, h rel.Value) (rel.Value, error) {
			var bs []byte
			switch b := x.(type) {
			case rel.String:
				bs = []byte(b.String())
			case rel.Bytes:
				bs = b.Bytes()
			default:
				return nil, errors.Errorf("first arg to xlsx.decodeToRelation must be string or bytes, not %T", b)
			}

			iv, ok := i.(rel.Number)
			if !ok {
				return nil, errors.Errorf("second arg to xlsx.decodeToRelation must be integer, not %T", i)
			}
			ix, ok := iv.Int()
			if !ok {
				return nil, errors.Errorf("second arg to xlsx.decodeToRelation must be integer, not %v", i)
			}
			hv, ok := h.(rel.Number)
			if !ok {
				return nil, errors.Errorf("third arg to xlsx.decodeToRelation must be integer, not %T", h)
			}
			hx, ok := hv.Int()
			if !ok {
				return nil, errors.Errorf("third arg to xlsx.decodeToRelation must be integer, not %v", h)
			}
			return bytesXlsxToArrai(bs, ix, hx)
		}),
	)
}

func bytesXlsxToArrai(bs []byte, sheetIndex int, headerRow int) (rel.Value, error) {
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
	if len(rows) <= headerRow {
		return rel.None, nil
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
			if col+1 >= scol && col+1 <= ecol && row+1 >= srow && row+1 <= erow {
				return m.GetCellValue(), nil
			}
		}
		return "", nil
	}

	badChars, err := regexp.Compile(`[,.'"?:/\s(){}\[\]]`)
	if err != nil {
		return nil, err
	}
	cols := []string{}
	for _, cell := range rows[headerRow] {
		// Commas in names are difficult to use in arr.ai (e.g. in nest attr lists).
		name := strcase.ToSnake(cell)
		name = badChars.ReplaceAllString(name, "_")
		cols = append(cols, name)
	}
	rowTuples := []rel.Value{}
	for i := headerRow + 1; i < len(rows); i++ {
		attrs := []rel.Attr{rel.NewIntAttr("@row", i)}
		hasVals := false
		for j, name := range cols {
			if name == "" {
				continue
			}
			val, err := getValue(rows, i, j)
			if err != nil {
				return nil, err
			}
			if val != "" {
				hasVals = true
			}
			attrs = append(attrs, rel.NewStringAttr(name, []rune(val)))
		}
		if hasVals && len(attrs) > 0 {
			rowTuples = append(rowTuples, rel.NewTuple(attrs...))
		}
	}
	return rel.NewSet(rowTuples...)
}
