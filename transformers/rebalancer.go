package transformers

import (
	"quant/core"
	"quant/model"
	"quant/utils"
)

type Rebalancer struct {
	Leverage  float64
	symbol    string
	portfolio *core.Portfolio
}

var _ Transformer = (*Rebalancer)(nil)

func (r *Rebalancer) rebalanceOnPrice(newPrice model.Price) {
	// sell all assets at newPrice
	r.portfolio.ExecuteOrder(&core.Order{
		Side:     core.SELL,
		Symbol:   r.symbol,
		Quantity: utils.Float64P(r.portfolio.Balance[r.symbol]),
		Price:    newPrice,
	})

	// repay all debt
	r.portfolio.Repay(r.portfolio.BorrowedQuote)

	positionSize := r.portfolio.AvailableQuote * r.Leverage

	// achieve desired leverage position
	if r.Leverage > 1 {
		r.portfolio.Borrow(r.portfolio.AvailableQuote * (r.Leverage - 1))
	}

	// invest with desired quote
	r.portfolio.ExecuteOrder(&core.Order{
		Side:   core.BUY,
		Symbol: r.symbol,
		Quote:  utils.Float64P(positionSize),
		Price:  newPrice,
	})
}

func (r *Rebalancer) TransformCandles(candles []*model.Candle) []*model.Candle {
	r.symbol = candles[0].Symbol
	r.portfolio = core.NewPortfolio()
	r.portfolio.Invest(1.0)

	// buy all assets at initial price
	r.portfolio.ExecuteOrder(&core.Order{
		Side:   core.BUY,
		Symbol: r.symbol,
		Quote:  utils.Float64P(r.portfolio.AvailableQuote),
		Price:  candles[0].Open,
	})

	transformed := make([]*model.Candle, 0, len(candles))
	for _, candle := range candles {
		newCandle := &model.Candle{
			Symbol: candle.Symbol,
			Time:   candle.Time,
		}

		r.rebalanceOnPrice(candle.Open)
		newCandle.Open = model.Price(r.portfolio.TotalAssets(&core.FixedMarketViewer{Price: candle.Open}))

		r.rebalanceOnPrice(candle.Close)
		newCandle.Close = model.Price(r.portfolio.TotalAssets(&core.FixedMarketViewer{Price: candle.Close}))

		transformed = append(transformed, newCandle)
	}

	return transformed
}
