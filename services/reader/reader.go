package reader

import (
	"errors"
	"io"
	"kolibra/database"
	"kolibra/tools"
	"time"

	"github.com/gen2brain/go-fitz"
	"github.com/patrickmn/go-cache"
)

var readerCache *cache.Cache

func CreateReaderCache() error {
	if readerCache == nil {
		readerCache = cache.New(10*time.Minute, 10*time.Minute)
		return nil
	}
	return errors.New("Pool already exists!")
}

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
	doc, found := readerCache.Get(book.Path)
	if !found {
		doc, err := fitz.New(book.Path)
		if err != nil {
			defer doc.Close()
			return "", err
		}
		readerCache.Set(book.Path, doc, cache.DefaultExpiration)
		return doc.HTML(int(chapter.Start), false)
	}

	return doc.(*fitz.Document).HTML(int(chapter.Start), false)
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
