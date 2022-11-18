package model

type Candle struct {
	Symbol string
	Time   TimeInterval
	Open   Price
	Close  Price
	High   Price
	Low    Price
}
