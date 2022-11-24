package simulator

import "quant/model"

type Portfolio struct {
	Balance map[string]float64
	Quote   float64
}

func (e *Portfolio) getQuantity(symbol string) float64 {
	qty := e.Balance[symbol]
	return qty
}
func (e *Portfolio) addQuantity(symbol string, diff float64) {
	current := e.getQuantity(symbol)
	e.Balance[symbol] = current + diff
}

func (e *Portfolio) ExecuteOrder(o *Order) {
	quoteQTY := o.GetQuoteQty()
	if o.Side == BUY {
		e.Quote -= quoteQTY
		e.addQuantity(o.Symbol, o.GetQuantity())
	} else {
		e.Quote += quoteQTY
		e.addQuantity(o.Symbol, -o.GetQuantity())
	}

	if o.OnExecuted != nil {
		o.OnExecuted()
	}
}

func (e *Portfolio) TotalBalance(getPrice func(symbol string) model.Price) float64 {
	total := e.Quote
	for symbol, quantity := range e.Balance {
		total += getPrice(symbol).GetQuote(quantity)
	}
	return total
}
