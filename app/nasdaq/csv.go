package main

import (
	"quant/app"
	"quant/crawler"
	"quant/model"
	"quant/storage"
)

func main() {
	nasdaq := crawler.NewNasdaq()
	db := storage.NewStorage("./data/storage.db")

	app.ReadFromCSVAndSave(nasdaq, model.AAPL, "./data/AAPL.csv", db)
	app.ReadFromCSVAndSave(nasdaq, model.SPX, "./data/SPX.csv", db)
}
