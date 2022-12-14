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

	// transform all candles by applying a 4x leverage rebalancer to it
	//candles = (&transformers.Rebalancer{Leverage: 4}).TransformCandles(candles)

	// transform all candles by replacing each one with the minimum price in the next MaxDuration
	//candles = (&transformers.CrystalGlobe{MaxDuration: utils.Month * time.Duration(3)}).TransformCandles(candles)

	for _, candle := range candles {
		sim.ProcessCandle(candle)
	}

	sim.Portfolio.LogProfit(sim)
}
