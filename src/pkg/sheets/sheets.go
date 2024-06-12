package sheets

import (
	"context"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Service wraps the Google Sheets API service
type Service struct {
	srv *sheets.Service
}

// NewService creates a new instance of Service
func NewService(ctx context.Context, credentialsPath string) (*Service, error) {
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

	return &Service{srv: srv}, nil
}

// WriteData writes data to the specified range in the spreadsheet
func (s *Service) WriteData(spreadsheetId, writeRange string, values [][]interface{}) error {
	var vr sheets.ValueRange
	vr.Values = values

	_, err := s.srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}

	log.Println("Data successfully written to sheet")
	return nil
}
