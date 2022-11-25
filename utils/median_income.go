package utils

import "time"

// based on https://www.statista.com/statistics/200838/median-household-income-in-the-united-states/
var (
	medianIncomes = []float64{
		30636,                             // 1990
		31241, 32264, 34076, 35492, 37005, // 1995
		38885, 40696, 41990, 42228, 42409, // 2000
		43318, 44334, 46326, 48201, 50233, // 2005
		50303, 49777, 49276, 50054, 51017, // 2010
		51939, 53585, 53657, 56516, 59039, // 2015
		61372, 61136, 63179, 68703, 68010, // 2020
		70784, // 2021
	}
)

const (
	startTimestamp = 631152000
	endTimestamp   = 1609459200
	tickSize       = 32610240
)

func GetMedianIncome(moment time.Time) float64 {
	timestamp := moment.Unix()
	if timestamp <= startTimestamp {
		return medianIncomes[0]
	} else if timestamp >= endTimestamp {
		return medianIncomes[30]
	}

	tick := (timestamp - startTimestamp) / tickSize
	remaining := float64((timestamp-startTimestamp)%tickSize) / tickSize
	return medianIncomes[tick]*(remaining) + medianIncomes[tick+1]*(1.0-remaining)
}

// GetNormalizedMedianIncome returns median incomes scaled to [0, 1]
func GetNormalizedMedianIncome(moment time.Time) float64 {
	return GetMedianIncome(moment) / medianIncomes[30]
}
