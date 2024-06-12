package reader

import (
	"errors"
	"io"
	"kolibra/config"
	"kolibra/database"
	"kolibra/tools"
	"time"

	"github.com/gen2brain/go-fitz"
	"github.com/patrickmn/go-cache"
)

var readerCache *cache.Cache

func CreateReaderCache() error {
	if readerCache == nil {
		readerCache = cache.New(
			time.Duration(config.Config.Advance.ReaderCachedMinutes)*time.Minute,
			time.Duration(config.Config.Advance.ReaderCachedMinutes)*time.Minute,
		)
		return nil
	}
	return errors.New("Pool already exists!")
}

func ReadChapter(book *database.Book, chapter *database.Chapter) (any, error) {
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

func ReadChapterEPUB_PDF(book *database.Book, chapter *database.Chapter) (any, error) {
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
	reader, found := readerCache.Get(book.Path)
	if !found {
		txtReader, err := tools.OpenTxtByEncode(book.Path)
		if err != nil {
			return "", err
		}
		if _, err := txtReader.F.Seek(chapter.Start, io.SeekStart); err != nil {
			return "", err
		}
		bytes := make([]byte, chapter.Length)
		_, err = io.ReadFull(txtReader.Reader, bytes)
		readerCache.Set(book.Path, txtReader, cache.DefaultExpiration)
		return string(bytes), err
	}

	txtReader := reader.(*tools.TxtReader)
	if _, err := txtReader.F.Seek(chapter.Start, io.SeekStart); err != nil {
		return "", err
	}
	bytes := make([]byte, chapter.Length)
	_, err := io.ReadFull(txtReader.Reader, bytes)
	return string(bytes), err
}
