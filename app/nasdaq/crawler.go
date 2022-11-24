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

	nasdaq := crawler.NewNasdaq()
	db := storage.NewStorage("./data/storage.db")

	app.CrawlAndSave(nasdaq, model.AAPL, startTime, db)
	app.CrawlAndSave(nasdaq, model.SPX, startTime, db)
}
