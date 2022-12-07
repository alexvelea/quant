package transformers

import (
	"quant/model"
)

type Transformer interface {
	TransformCandles([]*model.Candle) []*model.Candle
}
