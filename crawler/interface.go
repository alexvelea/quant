package crawler

import (
	"quant/model"
)

type Crawler interface {
	GetCandles(symbol string, interval model.TimeInterval) []*model.Candle
}
