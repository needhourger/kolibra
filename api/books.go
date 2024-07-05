package api

import (
	DB "kolibra/models"
	"kolibra/services/reader"
	"log"

	"github.com/gin-gonic/gin"
)

func GetAllBooks(c *gin.Context) {
	books, err := DB.GetAllBooks()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Books count: %d", len(*books))
	c.JSON(200, gin.H{"data": *books})
}

func GetBook(c *gin.Context) {
	bookID := c.Param("id")
	book, err := DB.GetBookByID(bookID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Book: %v", book)
	c.JSON(200, gin.H{"data": book})
}

func GetBookChapters(c *gin.Context) {
	bookID := c.Param("id")
	if _, err := DB.GetBookByID(bookID); err != nil {
		c.JSON(404, gin.H{"error": "No such book"})
		return
	}
	chapters, err := DB.GetChaptersByBookID(bookID)
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
	book, err := DB.GetBookByID(bookID)
	if err != nil {
		c.JSON(404, gin.H{"error": "No such Book"})
		return
	}
	chapter, err := book.GetChapterByID(chapterID)
	if err != nil {
		c.JSON(404, gin.H{"error": "No such Chapter"})
		return
	}
	c.JSON(200, gin.H{"data": chapter})
}

func GetChapterContent(c *gin.Context) {
	bookID := c.Param("id")
	chapterID := c.Param("cid")
	book, err := DB.GetBookByID(bookID)
	if err != nil {
		c.JSON(404, gin.H{"error": "No such Book"})
		return
	}
	chapter, err := book.GetChapterByID(chapterID)
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
	if _, err := DB.GetBookByID(bookID); err != nil {
		c.JSON(404, gin.H{"error": "No such book"})
		return
	}
	err := DB.DeleteBookByID(bookID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Book deleted"})
}
