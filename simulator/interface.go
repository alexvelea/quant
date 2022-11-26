package simulator

import (
	"quant/model"
	"time"
)

type MarketViewer interface {
	GetPrice(symbol string) model.Price
}

type Interactor interface {
	AddOrder(order *Order)
	MarketOrder(order *Order)

	MarketViewer
	GetPortfolio() *Portfolio
}

type Consumer interface {
	OnNewCandle(sim Interactor, start time.Time)
	OnPriceUpdate(sim Interactor, newPrice model.Price)
}
