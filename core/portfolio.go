package core

import (
	"fmt"
	"log"
	"math"
	"quant/utils"
)

type Portfolio struct {
	Balance        map[string]float64
	AvailableQuote float64
	InvestedQuote  float64
	BorrowedQuote  float64
}

func NewPortfolio() *Portfolio {
	return &Portfolio{Balance: make(map[string]float64)}
}

func (e *Portfolio) getQuantity(symbol string) float64 {
	qty := e.Balance[symbol]
	return qty
}
func (e *Portfolio) addQuantity(symbol string, diff float64) {
	current := e.getQuantity(symbol)
	e.Balance[symbol] = current + diff
}

func (e *Portfolio) AssetsValue(market MarketViewer) float64 {
	balance := 0.0
	for symbol, quantity := range e.Balance {
		balance += market.GetPrice(symbol).GetQuote(quantity)
	}
	return balance
}

func (e *Portfolio) TotalAssets(market MarketViewer) float64 {
	assets := e.AssetsValue(market)
	return e.AvailableQuote - e.BorrowedQuote + assets
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

func (e *Portfolio) Borrow(quote float64) {
	e.AvailableQuote += quote
	e.BorrowedQuote += quote
}

func (e *Portfolio) Repay(quote float64) {
	utils.PanicIf(quote > e.AvailableQuote, fmt.Errorf("not enough funds to repay"))
	e.AvailableQuote -= quote
	e.BorrowedQuote -= quote
}

func (e *Portfolio) LogProfit(market MarketViewer) {
	assets := e.AssetsValue(market)
	balance := e.AvailableQuote - e.BorrowedQuote + assets

	log.Printf("Invested:%v\n", e.InvestedQuote)
	log.Printf("Available:%v\n", e.AvailableQuote)
	log.Printf("Borrowed:%v\n", e.BorrowedQuote)
	log.Printf("Assets:%v\n", assets)
	log.Printf("Balance:%v\n", balance)
	for symbol, amount := range e.Balance {
		log.Printf("%v -- %v @ %v\n", symbol, amount, market.GetPrice(symbol))
	}
	log.Printf("Profit margin:%.2f%% / %.2f (log2)", (balance/e.InvestedQuote-1.0)*100, math.Log2(balance/e.InvestedQuote))
}
