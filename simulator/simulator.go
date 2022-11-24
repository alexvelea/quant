package simulator

import (
	"log"
	"quant/model"
)

type Simulator struct {
	Book      *Book
	Portfolio *Portfolio
	Consumers []Consumer
}

func NewSimulator() *Simulator {
	portfolio := &Portfolio{Balance: make(map[string]float64)}
	return &Simulator{
		Book: &Book{
			ActiveOrders: make([]*Order, 0),
			Executor:     portfolio,
		},
		Portfolio: portfolio,
		Consumers: make([]Consumer, 0),
	}
}

func (s *Simulator) ProcessCandle(candle *model.Candle) {
	log.Printf("Processing candle time:%v", candle.Time.Start.Format("2006-01-02"))
	s.Book.OnPrice([]model.Price{candle.Open})
	s.Book.UpdatePrice(candle.Open)
	for _, consumer := range s.Consumers {
		consumer.OnPriceUpdate(s, candle.Open)
		consumer.OnNewCandle(s, candle.Time.Start)
	}

	s.Book.OnPrice([]model.Price{candle.Low, candle.High})
	s.Book.OnPrice([]model.Price{candle.Close})

	s.Book.UpdatePrice(candle.Close)
	for _, consumer := range s.Consumers {
		consumer.OnPriceUpdate(s, candle.Close)
	}
}

func (s *Simulator) CurrentPrice() model.Price {
	return s.Book.LastPrice
}
func (s *Simulator) MarketValue() float64 {
	return s.Portfolio.TotalBalance(func(symbol string) model.Price {
		return s.Book.LastPrice
	})
}

func (s *Simulator) AddOrder(order *Order) {
	s.Book.AddOrder(order)
}
func (s *Simulator) MarketOrder(order *Order) {
	order.Price = s.Book.LastPrice
	s.Portfolio.ExecuteOrder(order)
}
