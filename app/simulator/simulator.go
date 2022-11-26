package main

import (
	"quant/model"
	"quant/simulator"
	"quant/storage"
	"quant/strategy"
)

func main() {
	symbol := model.SPX
	sim := simulator.NewSimulator([]string{symbol})

	db := storage.NewStorage("./data/storage.db")
	dca := strategy.NewDollarCostAverageStrategy(symbol)

	sim.Consumers = append(sim.Consumers, dca)

	candles := db.GetCandles(symbol)
	for _, candle := range candles {
		sim.ProcessCandle(candle)
	}

	sim.Portfolio.LogProfit(sim)
}
