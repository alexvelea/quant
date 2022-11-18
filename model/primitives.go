package model

import (
	"encoding/json"
	"time"
)

var (
	BTC = "BTC"
)

type TimeInterval struct {
	Open     time.Time
	Close    time.Time
	Duration time.Duration
}

func (ti *TimeInterval) NumTicks() int64 {
	return (ti.Close.UnixNano() - ti.Open.UnixNano()) / ti.Duration.Nanoseconds()
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
		return Price(t)
	default:
		panic("unknown type")
	}
}
