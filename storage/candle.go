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
