package database

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
	Book              Book
}

// Book CRUD
func CreateBook(book *Book) error {
	return db.Create(book).Error
}

func UpdateBook(book *Book) error {
	return db.Save(book).Error
}

func GetAllBooks() (*[]Book, error) {
	books := &[]Book{}
	err := db.Find(books).Error
	return books, err
}

func GetBookByID(id string) (*Book, error) {
	book := &Book{}
	err := db.First(book, "id = ?", id).Error
	return book, err
}

func CheckBookFileHash(hash string) (*Book, bool) {
	book := &Book{}
	err := db.Where("hash = ?", hash).First(book).Error
	return book, err == nil
}

func DeleteBookByID(id string) error {
	return db.Delete(&Book{}, id).Error
}

func CreateChapter(chapter *Chapter) error {
	return db.Create(chapter).Error
}

func GetChaptersByBookID(bookID any) (*[]Chapter, error) {
	chapters := &[]Chapter{}
	err := db.Where("book_id = ?", bookID).Find(chapters).Error
	return chapters, err
}

func (book *Book) GetChapterByID(cid any) (*Chapter, error) {
	chapter := &Chapter{}
	err := db.Model(book).Where("id = ?", cid).Association("Chapters").Find(chapter, "id = ?", cid)
	return chapter, err
}
