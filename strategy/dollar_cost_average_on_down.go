package strategy

import (
	"log"
	"quant/model"
	"quant/simulator"
	"quant/utils"
	"time"
)

type dollarCostAverageOnDown struct {
	Symbol     string
	pastPrices []model.Price
}

var _ simulator.Consumer = (*dollarCostAverageOnDown)(nil)

func (d *dollarCostAverageOnDown) OnNewCandle(sim simulator.Interactor, start time.Time) {
	toInvest := utils.GetNormalizedMedianIncome(start)
	sim.GetPortfolio().Invest(toInvest)
	available := sim.GetPortfolio().AvailableQuote
	currentPrice := sim.GetPrice(d.Symbol)

	d.pastPrices = append(d.pastPrices, currentPrice)
	pastSize := 5
	if len(d.pastPrices) > pastSize {

		numSmaller := 0.0
		weight := 0.0
		for i := len(d.pastPrices) - pastSize; i < len(d.pastPrices); i += 1 {
			weight += 1
			price := d.pastPrices[i]
			if currentPrice < price {
				numSmaller += weight * weight * (1.0 - (float64(currentPrice)/float64(price))*(float64(currentPrice)/float64(price)))
			}
		}

		size := (available / 3.0) * (numSmaller / float64(pastSize))
		sim.MarketOrder(&simulator.Order{
			Side:   simulator.BUY,
			Symbol: d.Symbol,
			Quote:  utils.Float64P(size),
			OnExecuted: func() {
				log.Printf("Bought based on discount metric size:%v price:%v available:%v\n", size, currentPrice, available)
			},
		})
	}
}

func (d *dollarCostAverageOnDown) OnPriceUpdate(sim simulator.Interactor, newPrice model.Price) {
}

func NewDollarCostAverageStrategyOnDown(symbol string) simulator.Consumer {
	return &dollarCostAverageOnDown{
		Symbol:     symbol,
		pastPrices: make([]model.Price, 0),
	}
}
