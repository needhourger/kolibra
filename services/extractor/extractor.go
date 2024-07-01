package extractor

import (
	"errors"
	"kolibra/database"
)

func Extract(book *database.Book) error {
	switch book.Extension {
	case ".txt":
		return extractTxt(book)
	}
	return errors.New("unsupported file type")
}
