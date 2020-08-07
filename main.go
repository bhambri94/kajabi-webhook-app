package main

import (
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction()
	sugar     = logger.Sugar()
)

func handleKajabiWebhook(ctx *fasthttp.RequestCtx) {
	sugar.Infof("received a webhook call!")
	ctx.Response.Header.Set("Content-Type", "application/json")
	sugar.Infof(string(ctx.Request.Body()))
}

func main() {
	sugar.Infof("starting kajabi webhook app server...")
	defer logger.Sync() // flushes buffer, if any

	router := fasthttprouter.New()
	router.POST("/v1/kajabi/send", handleKajabiWebhook)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
