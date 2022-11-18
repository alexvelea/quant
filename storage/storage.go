package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"quant/model"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(sqliteDB string) *Storage {
	db, err := gorm.Open(sqlite.Open(sqliteDB), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.Candle{})
	if err != nil {
		panic(err)
	}

	return &Storage{
		db: db,
	}
}
