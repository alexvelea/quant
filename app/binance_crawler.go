package main

import (
	"fmt"
	"time"

	"quant/crawler"
	"quant/model"
)

func main() {
	binance := crawler.NewBinance()
	startTime := time.UnixMilli(0)

	numCandles := 1
	duration := time.Hour
	interval := model.TimeInterval{
		Open:     startTime,
		Close:    startTime.Add(time.Duration(numCandles) * duration),
		Duration: duration,
	}
	candles := binance.GetCandles(model.BTC, interval)
	fmt.Printf("%v num %v\n", len(candles), candles[0])
}
