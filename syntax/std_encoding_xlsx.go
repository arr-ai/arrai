package syntax

import (
	"bytes"
	"context"
	"regexp"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/go-errors/errors"
	"github.com/iancoleman/strcase"

	"github.com/arr-ai/arrai/rel"
)

func stdEncodingXlsx() rel.Attr {
	return rel.NewTupleAttr(
		"xlsx",
		createFunc2Attr("decodeToRelation", func(_ context.Context, config rel.Value, x rel.Value) (rel.Value, error) {
			fn := "xlsx.decodeToRelation"
			var bs []byte
			c, ok := config.(*rel.GenericTuple)
			if !ok {
				return nil, errors.Errorf("first arg to %s must be tuple, not %s", fn, rel.ValueTypeAsString(config))
			}
			b, ok := x.(rel.Bytes)
			if !ok {
				return nil, errors.Errorf("second arg to %s must be string or bytes, not %s", fn, rel.ValueTypeAsString(x))
			}
			bs = b.Bytes()

			getConfigInt := func(key string, defaultVal int) (int, error) {
				if vv, ok := c.Get(key); ok {
					vn, ok := vv.(rel.Number)
					if !ok {
						return 0, errors.Errorf("%s config param to %s must be integer, not %s", key, fn, rel.ValueTypeAsString(vv))
					}
					v, ok := vn.Int()
					if !ok {
						return 0, errors.Errorf("%s config param to %s must be integer, not %v", key, fn, vn)
					}
					return v, nil
				}
				return defaultVal, nil
			}

			i, err := getConfigInt("sheet", 0)
			if err != nil {
				return nil, err
			}
			h, err := getConfigInt("headRow", 0)
			if err != nil {
				return nil, err
			}
			return bytesXlsxToArrai(bs, i, h)
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
