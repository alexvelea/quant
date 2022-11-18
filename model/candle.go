package model

type Candle struct {
	Symbol string       `gorm:"primaryKey"`
	Time   TimeInterval `gorm:"embedded"`
	Open   Price
	Close  Price
	High   Price
	Low    Price
}
