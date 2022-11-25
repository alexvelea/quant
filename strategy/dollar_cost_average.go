package strategy

import (
	"log"
	"quant/model"
	"quant/simulator"
	"quant/utils"
	"time"
)

type dollarCostAverage struct {
}

var _ simulator.Consumer = (*dollarCostAverage)(nil)

func (d dollarCostAverage) OnNewCandle(sim simulator.Interactor, start time.Time) {
	sim.MarketOrder(&simulator.Order{
		Quote: utils.Float64P(utils.GetNormalizedMedianIncome(start)),
		OnExecuted: func() {
			log.Printf("Bought some at %v\n", sim.CurrentPrice())
		},
	})
}

func (d dollarCostAverage) OnPriceUpdate(sim simulator.Interactor, newPrice model.Price) {
}

func NewDollarCostAverageStrategy() simulator.Consumer {
	return &dollarCostAverage{}
}
