package core

import "quant/model"

type FixedMarketViewer struct{ Price model.Price }

func (mv FixedMarketViewer) GetPrice(string) model.Price { return mv.Price }
