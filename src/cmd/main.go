package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {
	ctx := context.Background()

	// サービスアカウントの認証情報を読み込む
	CREADENTIALS_FILE := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	b, err := os.ReadFile(CREADENTIALS_FILE)
	if err != nil {
		log.Fatalf("Unable to read service account file: %v", err)
	}

	// Google Sheets APIクライアントを作成
	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// スプレッドシートIDと範囲
	spreadsheetId := os.Getenv("SPREAD_SHEETS_ID")
	writeRange := fmt.Sprintf("%s!A1:D5", os.Getenv("WORKSHEET_NAME"))

	// 書き込みデータ
	var vr sheets.ValueRange
	vr.Values = append(vr.Values, []interface{}{"Name", "Age", "Gender", "Location"})
	vr.Values = append(vr.Values, []interface{}{"山田 太郎", "30", "男", "東京"})
	vr.Values = append(vr.Values, []interface{}{"田中 花子", "25", "女", "大阪"})

	// データをスプレッドシートに書き込む
	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		// エラー詳細の出力
		log.Fatalf("Unable to write data to sheet: %v\nSpreadsheet ID: %s\nService Account Email: %s", err, spreadsheetId, config.Email)
	}

	fmt.Println("Data successfully written to sheet")
}
