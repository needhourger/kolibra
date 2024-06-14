package database

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

func GetInstance() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		db, err = gorm.Open(sqlite.Open(config.Settings.Database), &gorm.Config{})
		if err == nil {
			db.AutoMigrate(&User{}, &Book{}, &Chapter{})
			adminUser := User{Username: "admin", Password: "admin", Role: ADMIN, Email: ""}
			db.FirstOrCreate(&adminUser)
			log.Printf("Admin user: %v", adminUser)
		}
	})
	return db, err
}

func InitDatabase() error {
	_, err := GetInstance()
	return err
}
