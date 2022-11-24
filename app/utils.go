package app

import (
	"quant/crawler"
	"quant/model"
	"quant/storage"
	"time"
)

func CrawlAndSave(crawler crawler.Crawler, symbol string, startTime time.Time, db *storage.Storage) {
	duration := time.Hour * time.Duration(24)
	numCandles := 1000

	for startTime.After(time.Now()) == false {
		interval := model.TimeInterval{
			Start:    startTime,
			End:      startTime.Add(time.Duration(numCandles) * duration),
			Duration: duration,
		}
		candles := crawler.GetCandles(symbol, interval)

		err := db.InsertCandles(candles)
		if err != nil {
			panic(err)
		}

		startTime = interval.End
	}
}
