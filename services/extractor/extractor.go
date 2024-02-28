package extractor

import (
	"errors"
	"kolibra/database"
)

func Extract(book *database.Book) error {
	switch book.Extension {
	case ".txt":
		return extractTxt(book)
	case ".epub":
		return extractEPUB_PDF(book)
	case ".pdf":
		return extractEPUB_PDF(book)
	}
	return errors.New("unsupported file type")
}
