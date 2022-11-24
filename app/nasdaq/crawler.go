package main

import (
	"quant/app"
	"quant/crawler"
	"quant/model"
	"quant/storage"
	"quant/utils"
	"time"
)

func main() {
	// start time -- 1st of Jan 2020
	startTime, err := time.Parse("2006-01-02", "2018-01-01")
	utils.PanicIfErr(err)

	nasdaq := crawler.NewNasdaq()
	db := storage.NewStorage("./data/storage.db")

	app.CrawlAndSave(nasdaq, model.AAPL, startTime, db)
	app.CrawlAndSave(nasdaq, model.SPX, startTime, db)
}
