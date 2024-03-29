package extractor

import (
	"kolibra/database"
	"log"

	"github.com/gen2brain/go-fitz"
)

func extractEPUB_PDF(book *database.Book) error {
	doc, err := fitz.New(book.Path)
	if err != nil {
		return err
	}
	defer doc.Close()

	tocs, err := doc.ToC()
	if err != nil {
		return err
	}

	for index, toc := range tocs {
		var page int64
		if toc.Page == -1 {
			page = int64(index) + 1
		} else {
			page = int64(toc.Page)
		}
		log.Printf("Title: %s, Page: %d, URI: %s, Level: %d", toc.Title, toc.Page, toc.URI, toc.Level)
		chapter := database.Chapter{
			BookID: book.ID,
			Title:  toc.Title,
			Start:  page,
			URI:    toc.URI,
			Level:  toc.Level,
		}
		err := database.CreateChapter(&chapter)
		if err != nil {
			log.Printf("Error creating chapter: %v", err)
			continue
		}
	}

	book.Ready = true
	err = database.UpdateBook(book)
	if err != nil {
		return err
	}
	return nil
}
