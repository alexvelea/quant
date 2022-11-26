package simulator

import (
	"log"
)

type Portfolio struct {
	Balance        map[string]float64
	AvailableQuote float64
	InvestedQuote  float64
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
		e.AvailableQuote -= quoteQTY
		e.addQuantity(o.Symbol, o.GetQuantity())
	} else if o.Side == SELL {
		e.AvailableQuote += quoteQTY
		e.addQuantity(o.Symbol, -o.GetQuantity())
	} else {
		panic("order side not specified")
	}

	if o.OnExecuted != nil {
		o.OnExecuted()
	}
}

func (e *Portfolio) Invest(quote float64) {
	e.AvailableQuote += quote
	e.InvestedQuote += quote
}

func (e *Portfolio) LogProfit(market MarketViewer) {
	balance := e.AvailableQuote
	for symbol, quantity := range e.Balance {
		balance += market.GetPrice(symbol).GetQuote(quantity)
	}
	log.Printf("Invested:%v\n", e.InvestedQuote)
	log.Printf("Available:%v\n", e.AvailableQuote)
	for symbol, amount := range e.Balance {
		log.Printf("%v -- %v @ %v\n", symbol, amount, market.GetPrice(symbol))
	}
	log.Printf("Profit margin:%v", (balance/e.InvestedQuote-1.0)*100)
}
