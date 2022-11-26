package simulator

import (
	"fmt"
	"quant/model"
	"quant/utils"
)

type Book struct {
	Symbol       string
	ActiveOrders []*Order
	Executor     *Portfolio
	LastPrice    model.Price
}

func (b *Book) AddOrder(o *Order) {
	utils.PanicIf(o.Symbol != b.Symbol, fmt.Errorf("symbol of order (%v) doesn't match symbol of book (%v)", o.Symbol, b.Symbol))

	b.ActiveOrders = append(b.ActiveOrders, o)
}

func (b *Book) OnPrice(prices []model.Price) {
	// see if any orders are triggered
	toExecute := make([]*Order, 0)
	for i := range b.ActiveOrders {
		order := b.ActiveOrders[i]
		ok := false
		for _, newPrice := range prices {
			if order.IsTriggeredBy(newPrice) {
				ok = true
			}
		}

		if ok {
			toExecute = append(toExecute, order)
			b.ActiveOrders[i] = b.ActiveOrders[len(b.ActiveOrders)-1]
			b.ActiveOrders = b.ActiveOrders[:len(b.ActiveOrders)-1]
			i -= 1
		}
	}

	for _, order := range toExecute {
		b.Executor.ExecuteOrder(order)
	}
}

func (b *Book) UpdatePrice(newPrice model.Price) {
	b.LastPrice = newPrice
}
