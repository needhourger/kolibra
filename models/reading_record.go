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

func GetOrCreateReadingRecord(record *ReadingRecord) {
	db.Where(ReadingRecord{UserID: record.UserID, BookID: record.BookID}).Assign(ReadingRecord{ChapterID: record.ChapterID}).FirstOrCreate(record)
}
