package api

import (
	"kolibra/database/dao"
	"kolibra/database/model"
	"kolibra/middleware"
	"kolibra/services/reader"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllBooks(c *gin.Context) {
	books, err := dao.BookDAO.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Books count: %d", len(*books))
	c.JSON(200, gin.H{"data": *books})
}

func GetBook(c *gin.Context) {
	bookID := c.Param("id")
	book, err := dao.BookDAO.GetByID(bookID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Book: %v", book)
	c.JSON(200, gin.H{"data": book})
}

func GetBookChapters(c *gin.Context) {
	bookID := c.Param("id")
	if _, exist := dao.BookDAO.ExistByID(bookID); !exist {
		c.JSON(404, gin.H{"error": "No such book"})
		return
	}
	chapters, err := dao.ChapterDAO.Gets(map[string]interface{}{"book_id": bookID})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Chapters count: %d", len(*chapters))
	c.JSON(200, gin.H{"data": chapters})
}

func GetBookChapter(c *gin.Context) {
	bookID := c.Param("id")
	chapterID := c.Param("cid")
	if _, exist := dao.BookDAO.ExistByID(bookID); !exist {
		c.JSON(404, gin.H{"error": "No such Book"})
		return
	}
	chapter, err := dao.ChapterDAO.GetByID(chapterID)
	if err != nil {
		c.JSON(404, gin.H{"error": "No such Chapter"})
		return
	}
	c.JSON(200, gin.H{"data": chapter})
}

func GetChapterContent(c *gin.Context) {
	bookID := c.Param("id")
	chapterID := c.Param("cid")
	book, err := dao.BookDAO.GetByID(bookID)
	if err != nil {
		c.JSON(404, gin.H{"error": "No such Book"})
		return
	}
	chapter, err := dao.ChapterDAO.GetByID(chapterID)
	if err != nil {
		c.JSON(404, gin.H{"error": "No such Chapter"})
		return
	}
	content, err := reader.ReadChapter(book, chapter)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": content})
}

func DeleteBookByID(c *gin.Context) {
	bookID := c.Param("id")
	if _, exist := dao.BookDAO.ExistByID(bookID); !exist {
		c.JSON(404, gin.H{"error": "No such book"})
		return
	}
	err := dao.BookDAO.DeleteByID(bookID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Book deleted"})
}

func GetBookReadingRecord(c *gin.Context) {
	bookID := c.Param("bookID")
	user := middleware.GetUserFromJWT(c)
	record, err := dao.ReadingRecordDAO.Get(map[string]any{"book_id": bookID})
	if err != nil && err == gorm.ErrRecordNotFound {
		chapter, err := dao.ChapterDAO.Get(map[string]any{"book_id": bookID, "user_id": user.ID})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		c.JSON(200, gin.H{"data": chapter.ID})
	} else if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"data": record.ChapterID})
	}
}

func UpdateBook(c *gin.Context) {
	user := middleware.GetUserFromJWT(c)
	bookID := c.Param("id")
	book := model.Book{}
	if book, exist := dao.BookDAO.ExistByID(bookID); !exist {
		c.JSON(404, gin.H{"error": "Book not found"})
		return
	} else if !(book.UploaderID == user.ID || user.Role == model.ADMIN) {
		c.JSON(403, gin.H{"error": "Access denied"})
	}

	err := c.BindJSON(book)
	if err != nil {
		log.Printf("Update book bind json error: %v", err)
		c.JSON(400, gin.H{"error": "Bad request"})
		return
	}
	book.ID = bookID
	err = dao.BookDAO.Update(&book)
	if err != nil {
		log.Printf("Update Book info error: %v", err)
		c.JSON(500, gin.H{"error": "Interval server error"})
		return
	}

	c.JSON(200, gin.H{"message": "Update success"})
}
