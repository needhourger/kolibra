package model

import (
	"kolibra/config"
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func InitDatabase() {
	once.Do(func() {
		var err error
		db, err = gorm.Open(sqlite.Open(config.Settings.Database), &gorm.Config{})
		if err == nil {
			db.AutoMigrate(&User{}, &Book{}, &Chapter{}, &ReadingRecord{})
			adminUser := User{Username: "admin", Password: "admin", Role: ADMIN, Email: ""}
			db.FirstOrCreate(&adminUser)
			log.Printf("Admin user: %v", adminUser)
		} else {
			log.Fatalf("Open database error: %v", err)
		}
	})
}

func GetInstance() *gorm.DB {
	if db != nil {
		return db
	}
	log.Fatalf("Database get instance nil: %v", db)
	return nil
}
