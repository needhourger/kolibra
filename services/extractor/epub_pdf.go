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

	for _, toc := range tocs {
		chapter := database.Chapter{
			BookID: book.ID,
			Title:  toc.Title,
			Start:  int64(toc.Page),
			Level:  toc.Level,
		}
		err := database.CreateChapter(&chapter)
		if err != nil {
			log.Printf("Error creating chapter: %v", err)
			continue
		}
	}
	return nil
}
