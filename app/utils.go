package app

import (
	"quant/crawler"
	"quant/model"
	"quant/storage"
	"quant/utils"
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
		end := false
		if interval.End.After(time.Now()) {
			interval.End = time.Now()
			end = true
		}
		candles := crawler.GetCandles(symbol, interval)

		err := db.InsertCandles(candles)
		utils.PanicIfErr(err)

		startTime = interval.End

		if end {
			break
		}
	}
}
