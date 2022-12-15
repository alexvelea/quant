package main

import (
	"fmt"
	"math"
	"time"

	"quant/model"
	"quant/storage"
)

func addDiffs(diffs map[time.Time][]float64, candles []*model.Candle) {
	lastPrice := float64(candles[0].Open)
	for _, candle := range candles {
		price := float64(candle.Open)

		arr := diffs[candle.Time.Start.UTC()]
		if price != lastPrice {
			diffs[candle.Time.Start.UTC()] = append(arr, price/lastPrice-1.0)
		}
		lastPrice = price
	}
}

func getExpectedGrowthMultiDimensions(diffs [][]float64, size []float64) float64 {
	total := 0.0
	for _, diff := range diffs {
		sum := 0.0
		for i := range diff {
			sum += diff[i] * size[i]
		}
		total += math.Log(1.0 + sum)
	}
	total /= float64(len(diffs))

	return math.Exp(total)
}

func main() {
	diffsByTime := make(map[time.Time][]float64)
	symbols := []string{model.SPX, model.BTC}
	//symbols := []string{model.SPX}
	db := storage.NewStorage("./data/storage.db")

	for _, symbol := range symbols {
		candles := db.GetCandles(symbol)
		addDiffs(diffsByTime, candles)
	}

	var diffs [][]float64
	for _, diff := range diffsByTime {
		if len(diff) == len(symbols) {
			diffs = append(diffs, diff)
		}
	}

	fmt.Printf("got %v diffs\n", len(diffs))

	for _, a := range []float64{0, 0.1, 0.2, 0.5, 1, 2, 3, 4, 4.5, 4.75, 5, 5.25, 5.5, 6, 7, 8, 9, 10} {
		fmt.Printf("[")
		for _, b := range []float64{0, 0.1, 0.2, 0.5, 1, 2, 3, 4, 4.5, 4.75, 5, 5.25, 5.5, 6, 7, 8, 9, 10} {
			expected := getExpectedGrowthMultiDimensions(diffs, []float64{a, b})

			days := float64(365)
			total := math.Pow(expected, days)
			if math.IsNaN(total) {
				fmt.Printf("None,")
			} else {
				fmt.Printf("%v,", total)
			}
			//fmt.Printf("%v-%v %v %v\n", a, b, expected, total)
		}
		fmt.Printf("],\n")
	}
}
