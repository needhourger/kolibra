package database

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title      string
	Author     string
	Extension  string
	UploaderID uint
	Chapters   []Chapter
	Size       int64
	Path       string
	TitleRegex string `gorm:"default:''"`
	Ready      bool
}

type Chapter struct {
	gorm.Model
	Title  string
	Start  int64
	End    int64
	Length int64
	BookID uint
}

// Book CRUD
func CreateBook(book *Book) error {
	return db.Create(book).Error
}

func UpdateBook(book *Book) error {
	return db.Save(book).Error
}

func CreateChapter(chapter *Chapter) error {
	return db.Create(chapter).Error
}
