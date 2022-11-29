package crawler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"quant/model"
	"quant/utils"
	"time"
)

const (
	NasdaqStock = "stocks"
	NasdaqIndex = "index"
)

type nasdaqCrawler struct {
	baseURL string
	client  *http.Client
}

var _ Crawler = (*nasdaqCrawler)(nil)

func NewNasdaq() *nasdaqCrawler {
	return &nasdaqCrawler{
		baseURL: "https://api.nasdaq.com",
		client:  &http.Client{},
	}
}

// getSpecificInterval returns the duration in the format required by the origin, from a time.Duration
func (bc *nasdaqCrawler) getSpecificInterval(interval time.Duration) string {
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

func (bc *nasdaqCrawler) getSpecificTime(date time.Time) string {
	return date.UTC().Format("2006-01-02")
}

// getSpecificSymbol returns the origin symbol which should be used for a raw symbol, like "BTC"
func (bc *nasdaqCrawler) getSpecificSymbol(symbol string) string {
	if symbol == model.BTC {
		return "BTCBUSD"
	}

	panic(fmt.Sprintf("unknown symbol %v", symbol))
}

func (bc *nasdaqCrawler) getSymbolAssetClass(symbol string) string {
	if symbol == model.AAPL {
		return NasdaqStock
	}

	if symbol == model.SPX {
		return NasdaqIndex
	}

	panic(fmt.Sprintf("unknown symbol %v", symbol))
}

// parseCandleFromJson got from a JSON response into a generic model.Candle
func (bc *nasdaqCrawler) parseCandleFromJson(raw map[string]string) *model.Candle {
	start, err := time.Parse("01/02/2006", raw["date"])
	utils.PanicIfErr(err)

	c := &model.Candle{
		Time: model.TimeInterval{
			Start: start,
			End:   start.Add(time.Duration(24) * time.Hour),
		},
		Open:  model.ParsePriceFromJSON(raw["open"]),
		Close: model.ParsePriceFromJSON(raw["close"]),
		High:  model.ParsePriceFromJSON(raw["high"]),
		Low:   model.ParsePriceFromJSON(raw["low"]),
	}
	c.Time.Duration = c.Time.End.Sub(c.Time.Start)

	return c
}

// parseCandleFromJson got from a JSON response into a generic model.Candle
func (bc *nasdaqCrawler) parseCandleFromCsv(raw []string) *model.Candle {
	// Date,Close/Last,Volume,Open,High,Low
	start, err := time.Parse("01/02/2006", raw[0])
	utils.PanicIfErr(err)

	c := &model.Candle{
		Time: model.TimeInterval{
			Start: start,
			End:   start.Add(time.Duration(24) * time.Hour),
		},
		Open:  model.ParsePriceFromJSON(raw[3]),
		Close: model.ParsePriceFromJSON(raw[1]),
		High:  model.ParsePriceFromJSON(raw[4]),
		Low:   model.ParsePriceFromJSON(raw[5]),
	}
	c.Time.Duration = c.Time.End.Sub(c.Time.Start)

	return c
}

func (bc *nasdaqCrawler) emptyCandle(start time.Time, price model.Price) *model.Candle {
	c := &model.Candle{
		Time: model.TimeInterval{
			Start: start,
			End:   start.Add(time.Duration(24) * time.Hour),
		},
		Open:  price,
		Close: price,
		High:  price,
		Low:   price,
	}
	c.Time.Duration = c.Time.End.Sub(c.Time.Start)

	return c
}

type nasdaqCandleResponse struct {
	Data struct {
		Symbol      string
		TradesTable struct {
			Rows []map[string]string
		}
	}
}

func (bc *nasdaqCrawler) ReadFromAPI(symbol string, interval model.TimeInterval) []*model.Candle {
	if interval.Duration != time.Duration(24)*time.Hour {
		panic("invalid duration for nasdaq crawler: only day accepted")
	}

	numTicks := int(interval.NumTicks())

	params := url.Values{}
	params.Add("assetclass", bc.getSymbolAssetClass(symbol))
	params.Add("fromdate", bc.getSpecificTime(interval.Start))
	params.Add("todate", bc.getSpecificTime(interval.End.Add(-time.Hour)))
	params.Add("limit", fmt.Sprintf("%v", numTicks))

	url := bc.baseURL + fmt.Sprintf("/api/quote/%s/historical?", symbol) + params.Encode()
	log.Printf(`performing nasdaq API call %v`, url)

	req, err := http.NewRequest("GET", url, nil)
	utils.PanicIfErr(err)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	resp, err := bc.client.Do(req)
	utils.PanicIfErr(err)

	defer func(Body io.ReadCloser) {
		utils.PanicIfErr(Body.Close())
	}(resp.Body)

	decoded := new(nasdaqCandleResponse)
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(decoded)
	if err != nil {
		log.Fatalln(err)
	}

	rows := decoded.Data.TradesTable.Rows
	utils.ReverseSlice(rows)

	lastPrice := model.ParsePriceFromJSON(rows[0]["open"])
	nextTime := interval.Start

	candles := make([]*model.Candle, 0, len(rows))

	for i, nextCandle := 0, 0; i < numTicks; i += 1 {
		var candle *model.Candle
		if nextCandle < len(rows) {
			candle = bc.parseCandleFromJson(rows[nextCandle])
		}

		// correct time
		if candle != nil && candle.Time.Start == nextTime {
			candle.Symbol = symbol
			candles = append(candles, candle)
			lastPrice = candle.Close
			nextCandle += 1
		} else {
			// missing candle, adding fake one
			candle := bc.emptyCandle(nextTime, lastPrice)
			candle.Symbol = symbol
			candles = append(candles, candle)
		}
		nextTime = nextTime.Add(time.Duration(24) * time.Hour)
	}

	return candles
}

func (bc *nasdaqCrawler) ReadCSV(symbol string, path string) []*model.Candle {
	// open file
	f, err := os.Open(path)
	utils.PanicIfErr(err)

	// remember to close the file at the end of the program
	defer func(f *os.File) {
		utils.PanicIfErr(f.Close())
	}(f)

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	utils.PanicIfErr(err)

	// omit first line
	data = data[1:]
	utils.ReverseSlice(data)

	initialCandle := bc.parseCandleFromCsv(data[0])
	lastPrice := initialCandle.Open
	nextTime := initialCandle.Time.Start

	candles := make([]*model.Candle, 0, len(data))

	for i, nextCandle := 0, 0; i < len(data); i += 1 {
		var candle *model.Candle
		if nextCandle < len(data) {
			candle = bc.parseCandleFromCsv(data[nextCandle])
		}

		// correct time
		if candle != nil && candle.Time.Start == nextTime {
			// check if candle is actually empty
			if candle.Open == 0 {
				candle = bc.emptyCandle(nextTime, lastPrice)
			}

			nextCandle += 1
		} else {
			candle = bc.emptyCandle(nextTime, lastPrice)
		}
		candle.Symbol = symbol
		candles = append(candles, candle)
		lastPrice = candle.Close
		nextTime = nextTime.Add(time.Duration(24) * time.Hour)
	}

	return candles
}
