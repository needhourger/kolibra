package model

import "time"

type Book struct {
	ModelBase
	ISBN        string
	Title       string `gorm:"index"`
	Author      string
	Extension   string
	UploaderID  string
	Chapters    []Chapter
	Size        int64
	Path        string
	TitleRegex  string `gorm:"default:''"`
	Hash        string `gorm:"unique"`
	Encoding    string `gorm:"default:'utf-8'"`
	Cover       string
	Publisher   string
	PublishDate time.Time
	Description string
	Language    string
	Ready       bool
}

type Chapter struct {
	ModelBase
	Title             string
	Start             int64
	End               int64
	Length            int64
	Level             int `gorm:"default:1"`
	URI               string
	PreviousChapterID string `gorm:"default:''"`
	NextChapterID     string `gorm:"default:''"`
	BookID            string
}

func (book *Book) GetChapterByID(cid any) (*Chapter, error) {
	chapter := &Chapter{}
	err := GetInstance().Model(book).Where("id = ?", cid).Association("Chapters").Find(chapter, "id = ?", cid)
	return chapter, err
}
