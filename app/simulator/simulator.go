package main

import (
	"quant/core"
	"quant/model"
	"quant/storage"
	"quant/strategy"
)

func main() {
	symbol := model.SPX
	sim := core.NewSimulator([]string{symbol})

	db := storage.NewStorage("./data/storage.db")
	dca := strategy.NewDollarCostAverageStrategy(symbol)

	sim.Consumers = append(sim.Consumers, dca)

	candles := db.GetCandles(symbol)
	for _, candle := range candles {
		sim.ProcessCandle(candle)
	}

	sim.Portfolio.LogProfit(sim)
}
