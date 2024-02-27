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
	Ready      bool
}

type Chapter struct {
	gorm.Model
	Title  string
	Start  string
	End    string
	Length int
	BookID uint
}

// Book CRUD
func CreateBook(book *Book) error {
	return db.Create(book).Error
}
