package api

import (
	"kolibra/database"
	"log"

	"github.com/gin-gonic/gin"
)

func GetAllBooks(c *gin.Context) {
	books, err := database.GetAllBooks()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Books count: %d", len(books))
	c.JSON(200, books)
}

func GetBook(c *gin.Context) {
	bookID := c.Param("id")
	book, err := database.GetBookByID(bookID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Book: %v", book)
	c.JSON(200, book)
}

func GetBookChapters(c *gin.Context) {
	bookID := c.Param("id")
	if _, err := database.GetBookByID(bookID); err != nil {
		c.JSON(404, gin.H{"error": "No such book"})
		return
	}
	chapters, err := database.GetChaptersByBookID(bookID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Chapters count: %d", len(chapters))
	c.JSON(200, chapters)
}

func DeleteBookByID(c *gin.Context) {
	bookID := c.Param("id")
	if _, err := database.GetBookByID(bookID); err != nil {
		c.JSON(404, gin.H{"error": "No such book"})
		return
	}
	err := database.DeleteBookByID(bookID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Book deleted"})
}
