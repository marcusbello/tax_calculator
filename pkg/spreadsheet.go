package pkg

import (
	"google.golang.org/api/sheets/v4"
)

type Spreadsheet struct {
	sheet *sheets.Service
	ID    string
}

func NewSpreadsheet(sheet *sheets.Service, id string) *Spreadsheet {
	return &Spreadsheet{
		sheet: sheet,
		ID:    id,
	}
}

func (s *Spreadsheet) GetSheet(sheetName string, readRange string) (*sheets.ValueRange, error) {
	return s.sheet.Spreadsheets.Values.Get(s.ID, sheetName+"!"+readRange).Do()
}

func (s *Spreadsheet) AppendSheet(sheetName string, writeRange string, values [][]interface{}) (*sheets.AppendValuesResponse, error) {
	valueRange := &sheets.ValueRange{
		Values: values,
	}
	return s.sheet.Spreadsheets.Values.Append(s.ID, sheetName+"!"+writeRange, valueRange).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Do()
}

// lookupValueTaxByID looks up a value in the spreadsheet by ID and returns the corresponding tax amount.
func (s *Spreadsheet) LookupValueTaxByID(sheetName string, readRange string, id string) (uint64, error) {
	resp, err := s.GetSheet(sheetName, readRange)
	if err != nil {
		return 0, err
	}
	if len(resp.Values) == 0 {
		return 0, nil // No data found
	}
	for _, row := range resp.Values {
		if len(row) < 5 {
			continue // Skip rows that don't have enough columns
		}
		if row[0] == id {
			if taxAmount, ok := row[4].(int64); ok {

				return uint64(taxAmount), nil
			}
			return 0, nil // Tax amount is not an integer
		}
	}
	return 0, nil // ID not found
}