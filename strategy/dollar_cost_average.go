package strategy

import (
	"log"
	"quant/model"
	"quant/simulator"
	"quant/utils"
	"time"
)

type dollarCostAverage struct {
	Symbol string
}

var _ simulator.Consumer = (*dollarCostAverage)(nil)

func (d dollarCostAverage) OnNewCandle(sim simulator.Interactor, start time.Time) {
	toInvest := utils.GetNormalizedMedianIncome(start)
	sim.GetPortfolio().Invest(toInvest)

	sim.MarketOrder(&simulator.Order{
		Side:   simulator.BUY,
		Symbol: d.Symbol,
		Quote:  utils.Float64P(toInvest),
		OnExecuted: func() {
			log.Printf("Bought some at %v\n", sim.GetPrice(d.Symbol))
		},
	})
}

func (d dollarCostAverage) OnPriceUpdate(sim simulator.Interactor, newPrice model.Price) {
}

func NewDollarCostAverageStrategy(symbol string) simulator.Consumer {
	return &dollarCostAverage{
		Symbol: symbol,
	}
}
