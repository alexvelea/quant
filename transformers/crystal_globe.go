package transformers

import (
	"fmt"
	"time"

	"quant/model"
	"quant/utils"
)

type CrystalGlobe struct {
	MaxDuration time.Duration
	UseLowPrice bool

	pendingCandles []*model.Candle
}

var _ Transformer = (*CrystalGlobe)(nil)

func (r *CrystalGlobe) getBestPendingCandle() (candle *model.Candle) {
	first := r.pendingCandles[0]
	candle = &model.Candle{
		Symbol: first.Symbol,
		Time:   first.Time,
		Open:   first.Open,
		Close:  first.Close,
	}

	maxTime := candle.Time.Start.Add(r.MaxDuration)

	for _, current := range r.pendingCandles {
		if maxTime.After(current.Time.Start) {

			var price model.Price
			if r.UseLowPrice {
				price = current.Low
			} else {
				price = current.Open
				if price > current.Close {
					price = current.Close
				}
			}

			if candle.Open > price {
				candle.Open = price
				candle.Close = price
			}
		} else {
			break
		}
	}

	return candle
}

func (r *CrystalGlobe) ToString() string {
	nano := r.MaxDuration.Nanoseconds()
	if nano == 0 {
		return "none"
	}
	if nano > utils.Year.Nanoseconds() {
		return fmt.Sprintf("%v year", nano/utils.Year.Nanoseconds())
	}
	if nano > utils.Month.Nanoseconds() {
		return fmt.Sprintf("%v month", nano/utils.Month.Nanoseconds())
	}
	if nano > utils.Day.Nanoseconds() {
		return fmt.Sprintf("%v day", nano/utils.Day.Nanoseconds())
	}
	return fmt.Sprintf("%v hour", nano/utils.Hour.Nanoseconds())
}

func (r *CrystalGlobe) TransformCandles(candles []*model.Candle) []*model.Candle {
	transformed := make([]*model.Candle, 0, len(candles))

	for index := range candles {
		r.pendingCandles = candles[index:]
		transformed = append(transformed, r.getBestPendingCandle())
	}

	return transformed
}
