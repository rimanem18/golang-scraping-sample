package sheets

import (
	"context"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// SheetsService wraps the Google Sheets API service
type SheetsService struct {
	srv *sheets.Service
}

// NewSheetsService creates a new instance of SheetsService
func NewSheetsService(ctx context.Context, credentialsPath string) (*SheetsService, error) {
	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, err
	}

	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return &SheetsService{srv: srv}, nil
}

// WriteData writes data to the specified range in the spreadsheet
func (s *SheetsService) WriteData(spreadsheetId, writeRange string, values [][]interface{}) error {
	var vr sheets.ValueRange
	vr.Values = values

	_, err := s.srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}

	log.Println("Data successfully written to sheet")
	return nil
}
