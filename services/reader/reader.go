package reader

import (
	"errors"
	"io"
	"kolibra/database"
	"kolibra/tools"

	"github.com/gen2brain/go-fitz"
)

func ReadChapter(book *database.Book, chapter *database.Chapter) (string, error) {
	switch book.Extension {
	case ".txt":
		return ReadChapterTXT(book, chapter)
	case ".epub":
		return ReadChapterEPUB_PDF(book, chapter)
	case ".pdf":
		return ReadChapterEPUB_PDF(book, chapter)
	}
	return "", errors.New("Unsupported file type")
}

func ReadChapterEPUB_PDF(book *database.Book, chapter *database.Chapter) (string, error) {
	doc, err := fitz.New(book.Path)
	if err != nil {
		return "", err
	}
	defer doc.Close()

	return doc.HTML(int(chapter.Start), false)
}

func ReadChapterTXT(book *database.Book, chapter *database.Chapter) (string, error) {
	f, reader, err := tools.OpenFile(book.Path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.Seek(chapter.Start, io.SeekCurrent); err != nil {
		return "", err
	}
	bytes := make([]byte, chapter.Length)
	_, err = io.ReadFull(reader, bytes)
	return string(bytes), err
}
