package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"app/config"
	"app/pkg/scraper"
	"app/pkg/sheets"
)

func main() {
	// 設定ファイルの読み込み
	cfg := config.MustLoadConfig("./config.json")

	// スクレイパーの作成と実行
	scraper := scraper.NewScraper(cfg)
	scraper.Scrape()

	// Google Sheets APIの設定とデータ書き込み
	ctx := context.Background()
	CREDENTIAL_FILE := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	sheetsService, err := sheets.NewSheetsService(ctx, CREDENTIAL_FILE)
	if err != nil {
		log.Fatalf("シートサービスの作成に失敗しました。: %v", err)
	}

	// スクレイピング結果をまとめる
	values := [][]interface{}{
		{"Title", "URL"},
	}
	for _, article := range scraper.Results {
		values = append(values, []interface{}{article.Title, article.URL})
	}

	spreadsheetId := os.Getenv("SPREAD_SHEETS_ID")
	sheetName := cfg.Category
	sheetsService.CreateSheet(spreadsheetId, sheetName)

	// 保持しておいた情報をシートに書き込む
	writeRange := fmt.Sprintf("%s!A1", sheetName)
	err = sheetsService.WriteData(spreadsheetId, writeRange, values)
	if err != nil {
		log.Fatalf("データをシートに書き込むことができませんでした。: %v", err)
	}
}
