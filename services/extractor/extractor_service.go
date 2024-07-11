package extractor

import (
	"errors"
	"kolibra/database/model"
)

func Extract(book *model.Book) error {
	switch book.Extension {
	case ".txt":
		return extractTxt(book)
	}
	return errors.New("unsupported file type")
}
