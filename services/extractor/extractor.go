package extractor

import (
	"errors"
	DB "kolibra/models"
)

func Extract(book *DB.Book) error {
	switch book.Extension {
	case ".txt":
		return extractTxt(book)
	}
	return errors.New("unsupported file type")
}
