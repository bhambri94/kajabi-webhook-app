package main

import (
	"encoding/csv"
	"os"
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

	// log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))

	values := ReadCSV(configs.Configurations.CSVName)
	SheetName := GetSheetName()
	sheets.ClearSheet(SheetName)
	sheets.BatchWrite(SheetName, values)
}

func ReadCSV(file string) [][]interface{} {
	f, err := os.Open(file)
	if err != nil {

	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {

	}
	var finalValues [][]interface{}
	for _, line := range lines {
		row := make([]interface{}, len(line))
		for i, v := range line {
			row[i] = v
		}
		finalValues = append(finalValues, row)
	}
	return finalValues
}

func GetSheetName() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02T01")
}
