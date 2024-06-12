package main

import (
	"app/pkg/sheets"
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	CREADENTIALS_FILE := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	service, err := sheets.NewService(ctx, CREADENTIALS_FILE)
	if err != nil {
		log.Fatalf("Failed to create Sheets service: %v", err)
	}

	spreadsheetId := os.Getenv("SPREAD_SHEETS_ID")
	writeRange := fmt.Sprintf("%s!A1:D5", os.Getenv("WORKSHEET_NAME"))
	values := [][]interface{}{
		{"Name", "Age", "Gender", "Location"},
		{"小林 太郎", "30", "男", "長野"},
		{"佐藤 花子", "25", "女", "東京"},
	}

	err = service.WriteData(spreadsheetId, writeRange, values)
	if err != nil {
		log.Fatalf("Failed to write data to sheet: %v", err)
	}
}
