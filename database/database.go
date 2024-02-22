package database

import (
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetInstance() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		db, err = gorm.Open(sqlite.Open("./data.db"), &gorm.Config{})
		if err == nil {
			db.AutoMigrate(&User{})
		}
	})
	return db, err
}
