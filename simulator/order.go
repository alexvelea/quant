package simulator

import "quant/model"

type OrderSide string

const (
	BUY  OrderSide = "Buy"
	SELL OrderSide = "Sell"
)

type Order struct {
	Side       OrderSide
	Symbol     string
	Quantity   *float64
	Quote      *float64
	Price      model.Price
	CustomID   string
	OnExecuted func()
}

func (o *Order) IsTriggeredBy(price model.Price) bool {
	if o.Side == BUY {
		return o.Price.Cmp(price) <= 0
	} else {
		return o.Price.Cmp(price) >= 0
	}
}

func (o *Order) GetQuantity() float64 {
	if o.Quantity != nil {
		return *o.Quantity
	}
	return o.Price.FromQuote(*o.Quote)
}

func (o *Order) GetQuoteQty() float64 {
	if o.Quantity != nil {
		return o.Price.GetQuote(*o.Quantity)
	}
	return *o.Quote
}
