package pkg

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// ExportToExcel writes a generic table to an xlsx file at path.
// headers: slice of column headers
// rows: slice of rows, each row is a slice of string values. Length of each row may be <= len(headers).
func ExportToExcel(path string, headers []string, rows [][]string) error {
	f := excelize.NewFile()
	sheet := "Sheet1"

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(sheet, cell, h); err != nil {
			return fmt.Errorf("set header cell failed: %w", err)
		}
	}

	for rIdx, row := range rows {
		excelRow := rIdx + 2
		for cIdx, val := range row {
			cell, _ := excelize.CoordinatesToCellName(cIdx+1, excelRow)
			if err := f.SetCellValue(sheet, cell, val); err != nil {
				return err
			}
		}
	}

	if err := f.SaveAs(path); err != nil {
		return fmt.Errorf("save excel failed: %w", err)
	}
	return nil
}

// NOTE: Teacher-specific conversion moved to service layer. Use ExportToExcel for generic exports.
