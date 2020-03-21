package main

import (
	"github.com/arr-ai/arrai/rel"
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
)

func handleXlsx(data []byte) ([]byte, error) {
	file, err := xlsx.OpenBinary(data)
	if err != nil {
		return nil, err
	}
	sheets := rel.None
	for _, sheet := range file.Sheets {
		log.Printf("  Sheet %q", sheet.Name)
		cells := rel.None
		for r, row := range sheet.Rows {
			rowAttr := rel.NewIntAttr("row", r)
			for c, cell := range row.Cells {
				var value rel.Value
				switch cell.Type() {
				case xlsx.CellTypeString:
					value = rel.NewString([]rune(cell.Value))
				case xlsx.CellTypeStringFormula:
					value = rel.NewString([]rune(cell.String()))
				case xlsx.CellTypeNumeric:
					if cell.Value == "" {
						value = rel.None
					} else {
						n, err := cell.Float()
						if err != nil {
							return nil, errors.WrapPrefix(err, "reading number", 0)
						}
						value = rel.NewNumber(n)
					}
				case xlsx.CellTypeBool:
					value = rel.NewBool(cell.Bool())
				case xlsx.CellTypeDate:
					t, err := cell.GetTime(file.Date1904)
					if err != nil {
						return nil, errors.WrapPrefix(err, "reading date", 0)
					}
					year, month, day := t.Date()
					hour, minute, second := t.Clock()
					nanosecond := t.Nanosecond()
					value = rel.NewTuple(
						rel.NewTupleAttr("@time",
							rel.NewIntAttr("year", year),
							rel.NewIntAttr("month", int(month)),
							rel.NewIntAttr("day", day),
							rel.NewIntAttr("hour", hour),
							rel.NewIntAttr("minute", minute),
							rel.NewIntAttr("second", second),
							rel.NewIntAttr("nanosecond", nanosecond),
						),
					)
				case xlsx.CellTypeError:
					s := cell.String()
					value = rel.NewTuple(rel.NewStringAttr("@error", []rune(s)))
				default:
					return nil, errors.Errorf(
						"Unhandled cell type: %d", cell.Type())
				}
				cells = cells.With(rel.NewTuple(
					rowAttr,
					rel.NewIntAttr("col", (c)),
					rel.NewAttr("cell", value),
				))
			}
		}
		sheets = sheets.With(rel.NewTuple(
			rel.NewStringAttr("name", []rune(sheet.Name)),
			rel.NewAttr("cells", cells),
			// rel.NewAttr("cols", ),
			rel.NewBoolAttr("hidden", sheet.Hidden),
			rel.NewBoolAttr("selected", sheet.Selected),
		))
	}
	return []byte(sheets.String()), nil
}
