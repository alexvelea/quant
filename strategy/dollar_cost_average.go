package strategy

import (
	"log"
	"quant/core"
	"quant/model"
	"quant/utils"
	"time"
)

type dollarCostAverage struct {
	Symbol string
}

var _ core.Consumer = (*dollarCostAverage)(nil)

func (d dollarCostAverage) OnNewCandle(sim core.Interactor, start time.Time) {
	toInvest := utils.GetNormalizedMedianIncome(start)
	sim.GetPortfolio().Invest(toInvest)

	sim.MarketOrder(&core.Order{
		Side:   core.BUY,
		Symbol: d.Symbol,
		Quote:  utils.Float64P(toInvest),
		OnExecuted: func() {
			log.Printf("Bought some at %v\n", sim.GetPrice(d.Symbol))
		},
	})
}

func (d dollarCostAverage) OnPriceUpdate(sim core.Interactor, newPrice model.Price) {
}

func NewDollarCostAverageStrategy(symbol string) core.Consumer {
	return &dollarCostAverage{
		Symbol: symbol,
	}
}
