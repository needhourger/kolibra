package dao

import (
	"kolibra/database/model"
	"log"
)

var (
	BookDAO          *GenericDAO[model.Book]
	ChapterDAO       *GenericDAO[model.Chapter]
	UserDAO          *GenericDAO[model.User]
	ReadingRecordDAO *GenericDAO[model.ReadingRecord]
)

func InitDAO() {
	db := model.GetInstance()
	if db == nil {
		log.Fatalf("Dao init database instance is null: %v", db)
	}
	BookDAO = NewGenericDAO[model.Book](db)
	ChapterDAO = NewGenericDAO[model.Chapter](db)
	UserDAO = NewGenericDAO[model.User](db)
	ReadingRecordDAO = NewGenericDAO[model.ReadingRecord](db)
}
