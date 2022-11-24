package main

import (
	"quant/app"
	"quant/crawler"
	"quant/model"
	"quant/storage"
	"time"
)

func main() {
	// start time -- 1st of Jan 2020
	startTime := time.Unix(1577836800, 0).UTC()

	binance := crawler.NewBinance()
	db := storage.NewStorage("./data/storage.db")

	app.CrawlAndSave(binance, model.BTC, startTime, db)
}
