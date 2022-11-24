package main

import (
	"fmt"
	"quant/model"
	"quant/simulator"
	"quant/storage"
	"quant/strategy"
)

func main() {
	sim := simulator.NewSimulator()

	db := storage.NewStorage("./data/storage.db")
	dca := strategy.NewDollarCostAverageStrategy()

	sim.Consumers = append(sim.Consumers, dca)

	candles := db.GetCandles(model.AAPL)
	for _, candle := range candles {
		sim.ProcessCandle(candle)
	}
	fmt.Printf("profit %v\n", -100*sim.MarketValue()/sim.Portfolio.Quote)
}
