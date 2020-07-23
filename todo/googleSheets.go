package todo

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// ReadGoogleSheets reads data from Google Sheets
func ReadGoogleSheets() (*sheets.ValueRange, error) {
	spreadsheetID := os.Getenv("SHEET_ID")
	if "" == spreadsheetID {
		return nil, fmt.Errorf("SHEET_ID does not set")
	}

	secret, err := ioutil.ReadFile("secret.json")
	if err != nil {
		log.Fatal(err)
	}

	conf, err := google.JWTConfigFromJSON(secret, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(context.Background())
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatal(err)
	}

	// 参照
	readRange := "Sheet1!A:E"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Values) == 0 {
		return nil, fmt.Errorf("No data found from Google Sheets")
	}

	var programs []RadioProgram
	for _, row := range resp.Values {
		programs = append(programs, RadioProgram{
			name:      fmt.Sprint(row[0]),
			weekday:   Weekday(fmt.Sprint(row[1])),
			startTime: fmt.Sprint(row[2]),
			endTime:   fmt.Sprint(row[3]),
			stationID: fmt.Sprint(row[4]),
		})
	}

	return resp, nil
}
