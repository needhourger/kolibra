package models

type ReadingRecord struct {
	ModelBase
	UserID    string `gorm:"index:user_book_index,unique"`
	User      User
	BookID    string `gorm:"index:user_book_index,unique"`
	Book      Book
	ChapterID string
	Chapter   Chapter
}
