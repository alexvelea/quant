package simulator

import (
	"quant/model"
	"time"
)

type Interactor interface {
	AddOrder(order *Order)
	MarketOrder(order *Order)

	CurrentPrice() model.Price
}

type Consumer interface {
	OnNewCandle(sim Interactor, start time.Time)
	OnPriceUpdate(sim Interactor, newPrice model.Price)
}
