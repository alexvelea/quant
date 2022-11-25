package model

import (
	"encoding/json"
	"quant/utils"
	"strconv"
	"time"
)

var (
	BTC  = "BTC"
	AAPL = "AAPL"
	SPX  = "SPX"
)

type TimeInterval struct {
	Start    time.Time     `gorm:"primaryKey"`
	End      time.Time     ``
	Duration time.Duration `gorm:"primaryKey"`
}

func (ti *TimeInterval) NumTicks() int64 {
	return (ti.End.UnixNano() - ti.Start.UnixNano()) / ti.Duration.Nanoseconds()
}

func ParseTimeFromJSON(v interface{}) time.Time {
	num, err := v.(json.Number).Int64()
	utils.PanicIfErr(err)

	if num < 2000000000 {
		return time.Unix(num, 0)
	} else if num < 2000000000000 {
		return time.UnixMilli(num)
	} else {
		panic("too big number")
	}
}

type Price float64

func ParsePriceFromJSON(v interface{}) Price {
	switch t := v.(type) {
	case string:
		if t[0] == '$' {
			t = t[1:]
		}
		p, err := strconv.ParseFloat(t, 64)
		utils.PanicIfErr(err)
		return Price(p)
	default:
		panic("unknown type")
	}
}

func (p Price) GetQuote(quanty float64) float64 {
	return quanty * float64(p)
}
func (p Price) FromQuote(quote float64) float64 {
	return quote / float64(p)
}
