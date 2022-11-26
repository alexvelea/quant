package simulator

import (
	"fmt"
	"log"
	"quant/model"
	"quant/utils"
)

type Simulator struct {
	Books     []*Book
	Portfolio *Portfolio
	Consumers []Consumer
}

var _ Interactor = (*Simulator)(nil)

func NewSimulator(symbols []string) *Simulator {
	portfolio := &Portfolio{Balance: make(map[string]float64)}
	books := make([]*Book, 0)

	for _, symbol := range symbols {
		books = append(books, &Book{
			Symbol:       symbol,
			ActiveOrders: make([]*Order, 0),
			Executor:     portfolio,
		})
	}

	return &Simulator{
		Books:     books,
		Portfolio: portfolio,
		Consumers: make([]Consumer, 0),
	}
}

func (s *Simulator) getBook(symbol string) (book *Book) {
	for _, b := range s.Books {
		if b.Symbol == symbol {
			book = b
		}
	}
	utils.PanicIf(book == nil, fmt.Errorf("unable to find book for symbol (%v)", symbol))
	return
}

func (s *Simulator) ProcessCandle(candle *model.Candle) {
	log.Printf("Processing candle time:%v", candle.Time.Start.Format("2006-01-02"))

	book := s.getBook(candle.Symbol)
	book.OnPrice([]model.Price{candle.Open})
	book.UpdatePrice(candle.Open)
	for _, consumer := range s.Consumers {
		consumer.OnPriceUpdate(s, candle.Open)
		consumer.OnNewCandle(s, candle.Time.Start)
	}

	book.OnPrice([]model.Price{candle.Low, candle.High})
	book.OnPrice([]model.Price{candle.Close})

	book.UpdatePrice(candle.Close)
	for _, consumer := range s.Consumers {
		consumer.OnPriceUpdate(s, candle.Close)
	}
}

func (s *Simulator) GetPrice(symbol string) model.Price {
	return s.getBook(symbol).LastPrice
}

func (s *Simulator) AddOrder(order *Order) {
	s.getBook(order.Symbol).AddOrder(order)
}

func (s *Simulator) MarketOrder(order *Order) {
	price := s.GetPrice(order.Symbol)
	order.Price = price
	s.Portfolio.ExecuteOrder(order)
}

func (s *Simulator) GetPortfolio() *Portfolio {
	return s.Portfolio
}
