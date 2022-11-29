package app

import (
	"log"
	"quant/crawler"
	"quant/model"
	"quant/storage"
	"quant/utils"
	"time"
)

func ReadFromAPIAndSave(crawler crawler.Crawler, symbol string, startTime time.Time, db *storage.Storage) {
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
		candles := crawler.ReadFromAPI(symbol, interval)

		err := db.InsertCandles(candles)
		utils.PanicIfErr(err)

		startTime = interval.End

		if end {
			break
		}
	}
}

func ReadFromCSVAndSave(crawler crawler.Crawler, symbol string, path string, db *storage.Storage) {
	candles := crawler.ReadCSV(symbol, path)
	for _, candle := range candles {
		log.Printf("%+v\n", candle)
	}
	utils.PanicIfErr(db.InsertCandles(candles))
}
