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

const (
	defaultSheetIndex = 0
	defaultHeadRow    = 0
)

func stdEncodingXlsx() rel.Attr {
	return rel.NewTupleAttr(
		"xlsx",
		rel.NewNativeFunctionAttr(decodeAttr, func(_ context.Context, v rel.Value) (rel.Value, error) {
			return arraiBytesXlsxToArrai(v, defaultSheetIndex, defaultHeadRow, "xlsx.decode")
		}),
		createFunc2Attr("decodeToRelation", func(_ context.Context, config rel.Value, x rel.Value) (rel.Value, error) {
			//TODO: replace this with decoder instead of decodeToRelation
			fn := "xlsx.decodeToRelation"
			c, ok := config.(*rel.GenericTuple)
			if !ok {
				return nil, errors.Errorf("first arg to %s must be tuple, not %s", fn, rel.ValueTypeAsString(config))
			}

			i, err := getConfigInt(c, fn, "sheet", defaultSheetIndex)
			if err != nil {
				return nil, err
			}
			h, err := getConfigInt(c, fn, "headRow", defaultHeadRow)
			if err != nil {
				return nil, err
			}
			return arraiBytesXlsxToArrai(x, i, h, fn)
		}),
	)
}

func arraiBytesXlsxToArrai(v rel.Value, sheetIndex, headRow int, fn string) (rel.Value, error) {
	b, ok := v.(rel.Bytes)
	if !ok {
		return nil, errors.Errorf("second arg to %s must be string or bytes, not %s", fn, rel.ValueTypeAsString(v))
	}
	return bytesXlsxToArrai(b.Bytes(), sheetIndex, headRow)
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

	badChars, err := regexp.Compile(`[,.'"%?:/\s(){}\[\]]`)
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
	b := rel.NewSetBuilder()
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
			b.Add(rel.NewTuple(attrs...))
		}
	}
	return b.Finish()
}
