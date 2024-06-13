package sheets

import (
	"context"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Google Sheets API をラップする
type SheetsService struct {
	srv *sheets.Service
}

// インスタンスの作成
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

// 指定された範囲にデータを書き込む
func (s *SheetsService) WriteData(spreadsheetId, writeRange string, values [][]interface{}) error {
	var vr sheets.ValueRange
	vr.Values = values

	_, err := s.srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}

	log.Println("シートへの書き込みに成功しました。")
	return nil
}

// 新たなシートを作成する
func (s *SheetsService) CreateSheet(spreadsheetId, sheetName string) error {
	requests := []*sheets.Request{
		{
			AddSheet: &sheets.AddSheetRequest{
				Properties: &sheets.SheetProperties{
					Title: sheetName,
				},
			},
		},
	}

	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}

	_, err := s.srv.Spreadsheets.BatchUpdate(spreadsheetId, batchUpdateRequest).Do()
	if err != nil {
		return err
	}

	log.Printf("シート '%s' の作成に成功しました。\n", sheetName)
	return nil
}
