package main

import (
	"log"
	"quant/core"
	"quant/model"
	"quant/storage"
	"quant/strategy"
	"quant/utils"
	"time"
)

func getMinPriceInFuture(candles []*model.Candle, maxDuration *time.Duration) model.Price {
	bestPrice := &candles[0].Open

	if maxDuration == nil {
		return *bestPrice
	}

	t := candles[0].Time.Start.Add(*maxDuration)
	maxTime := &t

	for _, candle := range candles {
		if maxTime != nil && candle.Time.Start.After(*maxTime) {
			break
		}

		if bestPrice == nil || *bestPrice > candle.Open {
			bestPrice = &candle.Low
		}
	}

	return *bestPrice
}

func DurationP(duration time.Duration) *time.Duration {
	return &duration
}

func main() {
	symbol := model.SPX
	db := storage.NewStorage("./data/storage.db")
	candles := db.GetCandles(symbol)

	startTime := candles[0].Time.Start.Format("2006-01-02")
	endTime := candles[len(candles)-1].Time.Start.Format("2006-01-02")
	log.Printf("symbol:%v", symbol)
	log.Printf("startTime:%v endTime:%v", startTime, endTime)
	log.Println()

	configs := []struct {
		context     string
		maxDuration *time.Duration
	}{
		{context: "none (buy at open price daily) -- dollar cost averaging", maxDuration: nil},
		{context: "1 day (buy at min price daily)", maxDuration: DurationP(time.Duration(24) * time.Hour)},
		{context: "1 month", maxDuration: DurationP(time.Duration(24*30) * time.Hour)},
		{context: "3 months", maxDuration: DurationP(time.Duration(3*24*30) * time.Hour)},
		{context: "1 year", maxDuration: DurationP(time.Duration(12*24*30) * time.Hour)},
		{context: "3 year", maxDuration: DurationP(time.Duration(3*12*24*30) * time.Hour)},
	}

	for _, cfg := range configs {
		portfolio := core.NewPortfolio()

		sim := core.NewSimulator([]string{symbol})
		dca := strategy.NewDollarCostAverageStrategyOnDown(symbol)
		sim.Consumers = append(sim.Consumers, dca)

		var lastPrice model.Price
		for i := range candles {
			candle := candles[i]
			lastPrice = candle.Open
			toInvest := utils.GetNormalizedMedianIncome(candle.Time.Start)
			portfolio.Invest(toInvest)
			price := getMinPriceInFuture(candles[i:], cfg.maxDuration)

			portfolio.ExecuteOrder(&core.Order{
				Side:   core.BUY,
				Quote:  utils.Float64P(toInvest),
				Price:  price,
				Symbol: symbol,
			})
		}

		log.Printf("maxDuration: %v", cfg.context)
		portfolio.LogProfit(&core.FixedMarketViewer{Price: lastPrice})
		log.Println()
	}
}
