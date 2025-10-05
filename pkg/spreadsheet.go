package pkg

import "google.golang.org/api/sheets/v4"

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