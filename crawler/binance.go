package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"quant/model"
)

type binanceCrawler struct {
	baseURL string
	client  *http.Client
}

func NewBinance() *binanceCrawler {
	return &binanceCrawler{
		baseURL: "https://api.binance.com",
		client:  &http.Client{},
	}
}

// getSpecificInterval returns the duration in the format required by the origin, from a time.Duration
func (bc *binanceCrawler) getSpecificInterval(interval time.Duration) string {
	/*
		Kline/Candlestick chart intervals:
		s-> seconds; m -> minutes; h -> hours; d -> days; w -> weeks; M -> months
		1s
		1m
		3m
		5m
		15m
		30m
		1h
		2h
		4h
		6h
		8h
		12h
		1d
		3d
		1w
		1M
	*/

	switch interval {
	case time.Minute:
		return "1m"
	case time.Hour:
		return "1h"
	case time.Duration(24) * time.Hour:
		return "1d"
	default:
		panic("unknown")
	}
}

// getSpecificSymbol returns the origin symbol which should be used for a raw symbol, like "BTC"
func (bc *binanceCrawler) getSpecificSymbol(symbol string) string {
	if symbol == model.BTC {
		return "BTCBUSD"
	}

	panic(fmt.Sprintf("unknown symbol %v", symbol))
}

// parseCandle got from a JSON response into a generic model.Candle
func (bc *binanceCrawler) parseCandle(raw []interface{}) *model.Candle {
	c := &model.Candle{
		Time: model.TimeInterval{
			Open:  model.ParseTimeFromJSON(raw[0]),
			Close: model.ParseTimeFromJSON(raw[6]).Add(time.Millisecond),
		},
		Open:  model.ParsePriceFromJSON(raw[1]),
		Close: model.ParsePriceFromJSON(raw[2]),
		High:  model.ParsePriceFromJSON(raw[3]),
		Low:   model.ParsePriceFromJSON(raw[4]),
	}
	c.Time.Duration = c.Time.Close.Sub(c.Time.Open)

	return c
}

func (bc *binanceCrawler) GetCandles(symbol string, interval model.TimeInterval) []*model.Candle {
	params := url.Values{}
	params.Add("symbol", bc.getSpecificSymbol(symbol))
	params.Add("interval", bc.getSpecificInterval(interval.Duration))
	params.Add("startTime", fmt.Sprintf("%v", interval.Open.UnixMilli()))
	params.Add("limit", fmt.Sprintf("%v", interval.NumTicks()))

	request := bc.baseURL + `/api/v3/klines?` + params.Encode()
	log.Printf(`performing bincance API call %v`, request)
	resp, err := http.Get(request)

	decoded := make([][]interface{}, 0)
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&decoded)
	if err != nil {
		log.Fatalln(err)
	}

	candles := make([]*model.Candle, len(decoded))
	for i := range decoded {
		candles[i] = bc.parseCandle(decoded[i])
		candles[i].Symbol = symbol
	}

	return candles
}
