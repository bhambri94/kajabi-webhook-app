package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bhambri94/kajabi-webhook-app/configs"
	"github.com/bhambri94/kajabi-webhook-app/sheets"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction()
	sugar     = logger.Sugar()
)

func handleKajabiWebhookPost(ctx *fasthttp.RequestCtx) {
	sugar.Infof("received a webhook call POST!")
	ctx.Response.Header.Set("Content-Type", "application/json")
	sugar.Infof(string(ctx.Request.Body()))
}

func handleKajabiWebhookGet(ctx *fasthttp.RequestCtx) {
	sugar.Infof("received a webhook call GET!")
	ctx.Response.Header.Set("Content-Type", "application/json")
	sugar.Infof(string(ctx.Request.Body()))
}

func handleKajabiWebhookPut(ctx *fasthttp.RequestCtx) {
	sugar.Infof("received a webhook call PUT!")
	ctx.Response.Header.Set("Content-Type", "application/json")
	sugar.Infof(string(ctx.Request.Body()))
}

func main() {
	configs.SetConfig()
	sugar.Infof("starting kajabi app server...")
	defer logger.Sync() // flushes buffer, if any

	router := fasthttprouter.New()
	router.POST("/v1/kajabi/send", handleKajabiWebhookPost)
	router.GET("/v1/kajabi/send", handleKajabiWebhookGet)
	router.PUT("/v1/kajabi/send", handleKajabiWebhookPut)
	router.GET("/v1/kajabi/uploadcontacts", handleGenerateCSV)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

func handleGenerateCSV(ctx *fasthttp.RequestCtx) {
	sugar.Infof("received a push contacts to Google Sheet Request!")
	URL := GenerateCSV("contactsCsv.csv")
	sugar.Infof("URL generated for CSV:=" + URL)
	values := readCSVFromUrl(URL)
	sugar.Infof("CSV file from Email:=")
	fmt.Println(values)
	finalGoogleSheetValues := ReadCSV(values)
	sugar.Infof("Final Values being pushed to Google Sheet:=")
	fmt.Println(finalGoogleSheetValues)
	SheetName := configs.Configurations.SheetNameWithRange
	// sheets.ClearSheet(SheetName)
	sheets.BatchAppend(SheetName, finalGoogleSheetValues)
	ctx.Response.Header.Set("Content-Type", "application/json")
	sugar.Infof(string(ctx.Request.Body()))
}

func ReadCSV(lines [][]string) [][]interface{} {
	currentTime := time.Now()
	RefreshTime := currentTime.String()
	// fmt.Println(currentTime)
	// f, err := os.Open(file)
	// if err != nil {

	// }
	// defer f.Close()

	// lines, err := csv.NewReader(f).ReadAll()
	// if err != nil {
	// }
	var finalValues [][]interface{}
	headerRow := true
	secondRow := false
	for _, line := range lines {
		row := make([]interface{}, len(line)+1)
		if headerRow {
			headerRow = false
			secondRow = true
			continue
			row[0] = "Refresh Time"
		} else {
			if secondRow {
				blankRow := make([]interface{}, len(line))
				finalValues = append(finalValues, blankRow)
				secondRow = false
				row[0] = RefreshTime[:19]
			} else {
				row[0] = RefreshTime[:19]
			}
		}
		for i, v := range line {
			row[i+1] = v
		}
		finalValues = append(finalValues, row)
	}
	return finalValues
}

func GenerateCSV(fileName string) string {
	dat, _ := ioutil.ReadFile(fileName)
	actual := strings.Index(string(dat), "https://kajabi-storefronts-production.s3.amazonaws.co")
	end := strings.Index(string(dat), "in the next 3 days")
	// fmt.Println(actual)
	// fmt.Println(end)
	filteredString := (string(dat)[actual : end-5])
	filteredString = strings.Replace(filteredString, "=\\r\\n", "", -1)
	filteredString = strings.Replace(filteredString, "3D", "", -1)
	return filteredString
}

func readCSVFromUrl(url string) [][]string {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	// fmt.Println(resp)
	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	// reader.Comma = ';'
	data, err := reader.ReadAll()
	if err != nil {
		return nil
	}
	// fmt.Println(data)
	return data
}
