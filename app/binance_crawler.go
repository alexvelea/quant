package main

import (
	"time"

	"quant/crawler"
	"quant/model"
	"quant/storage"
)

func main() {
	// start time
	startTime := time.Unix(1577836800, 0)
	duration := time.Hour * time.Duration(24)
	numCandles := 1000

	binance := crawler.NewBinance()
	db := storage.NewStorage("./data/storage.db")

	for startTime.After(time.Now()) == false {
		interval := model.TimeInterval{
			Start:    startTime,
			End:      startTime.Add(time.Duration(numCandles) * duration),
			Duration: duration,
		}
		candles := binance.GetCandles(model.BTC, interval)

		err := db.InsertCandles(candles)
		if err != nil {
			panic(err)
		}

		startTime = interval.End
	}
}
