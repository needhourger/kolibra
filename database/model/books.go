package model

type Book struct {
	ModelBase
	ISBN       string
	Title      string `gorm:"index"`
	Author     string
	Extension  string
	UploaderID string
	Chapters   []Chapter
	Size       int64
	Path       string
	TitleRegex string `gorm:"default:''"`
	Hash       string `gorm:"unique"`
	Coding     string `gorm:"default:'utf-8'"`
	Ready      bool
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

func CheckBookFileHash(hash string) (*Book, bool) {
	book := &Book{}
	err := GetInstance().Where("hash = ?", hash).First(book).Error
	return book, err == nil
}

func (book *Book) GetChapterByID(cid any) (*Chapter, error) {
	chapter := &Chapter{}
	err := GetInstance().Model(book).Where("id = ?", cid).Association("Chapters").Find(chapter, "id = ?", cid)
	return chapter, err
}
