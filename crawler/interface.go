package crawler

import (
	"quant/model"
)

type Crawler interface {
	ReadFromAPI(symbol string, interval model.TimeInterval) []*model.Candle
	ReadCSV(symbol string, path string) []*model.Candle
}
