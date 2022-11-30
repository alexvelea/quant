package main

import (
	"fmt"
	"math"
	"quant/model"
	"quant/storage"
)

func getExpectedGrowth(difs []float64, size float64) float64 {
	total := 0.0
	for _, dif := range difs {
		total += math.Log(1.0 + dif*size)
	}

	total /= float64(len(difs))

	return math.Exp(total)
}

func main() {
	symbol := model.SPX
	db := storage.NewStorage("./data/storage.db")
	candles := db.GetCandles(symbol)
	lastPrice := float64(candles[0].Open)
	difs := make([]float64, 0)
	for _, candle := range candles {
		price := float64(candle.Open)

		if lastPrice != price {
			difs = append(difs, price/lastPrice-1.0)
		}

		lastPrice = price
	}

	for _, size := range []float64{0.1, 0.2, 0.5, 1, 2, 3, 4, 4.5, 4.75, 5, 5.25, 5.5, 6, 7, 8, 9, 10, 30} {
		expected := getExpectedGrowth(difs, size)

		days := float64(365)
		total := math.Pow(expected, days)
		fmt.Printf("%v %v %v\n", size, expected, total)
	}
}
