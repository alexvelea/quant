package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"quant/model"
	"quant/utils"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(sqliteDB string) *Storage {
	db, err := gorm.Open(sqlite.Open(sqliteDB), &gorm.Config{})
	utils.PanicIfErr(err)

	err = db.AutoMigrate(&model.Candle{})
	utils.PanicIfErr(err)

	return &Storage{
		db: db,
	}
}
