package database

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title string
	Author string
	Extension string
	UploaderID uint
	Chapters []Chapter
}

type Chapter struct {
	gorm.Model
	Title string
	Start string
	End string
	Length int
	BookID uint
}