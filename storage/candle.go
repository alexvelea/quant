package storage

import (
	"gorm.io/gorm/clause"
	"quant/model"
)

func (s *Storage) InsertCandles(candles []*model.Candle) error {
	result := s.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&candles)

	return result.Error
}

func (s *Storage) GetCandles(symbol string) []*model.Candle {
	query := s.db.Where(`symbol = ?`, symbol)
	var candles []*model.Candle
	result := query.Find(&candles)
	if result.Error != nil {
		panic(result.Error)
	}

	return candles
}
