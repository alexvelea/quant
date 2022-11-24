package model

import (
	"encoding/json"
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
	if err != nil {
		panic(err)
	}

	if num < 2000000000 {
		return time.Unix(num, 0)
	} else if num < 2000000000000 {
		return time.UnixMilli(num)
	} else {
		panic("too big number")
	}
}

type Price string

func ParsePriceFromJSON(v interface{}) Price {
	switch t := v.(type) {
	case string:
		if t[0] == '$' {
			return Price(t[1:])
		}
		return Price(t)
	default:
		panic("unknown type")
	}
}

func (p Price) Cmp(oth Price) int {
	a, err := strconv.ParseFloat(string(p), 64)
	if err != nil {
		panic(err)
	}
	b, err := strconv.ParseFloat(string(oth), 64)
	if err != nil {
		panic(err)
	}
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return +1
}

func (p Price) GetQuote(quanty float64) float64 {
	a, err := strconv.ParseFloat(string(p), 64)
	if err != nil {
		panic(err)
	}
	return quanty * a
}
func (p *Price) FromQuote(quote float64) float64 {
	a, err := strconv.ParseFloat(string(*p), 64)
	if err != nil {
		panic(err)
	}
	return quote / a
}
