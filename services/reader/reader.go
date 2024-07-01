package reader

import (
	"errors"
	"io"
	"kolibra/config"
	"kolibra/database"
	"kolibra/tools"
	"os"
	"time"

	"github.com/gen2brain/go-fitz"
	"github.com/patrickmn/go-cache"
)

var readerCache *cache.Cache

func CreateReaderCache() error {
	if readerCache == nil {
		readerCache = cache.New(
			time.Duration(config.Settings.Advance.ReaderCachedMinutes)*time.Minute,
			time.Duration(config.Settings.Advance.ReaderCachedMinutes)*time.Minute,
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
		// return doc.HTML(int(chapter.Start), false)
		return doc.SVG(int(chapter.Start))
	}

	// return doc.(*fitz.Document).HTML(int(chapter.Start), false)
	return doc.(*fitz.Document).SVG(int(chapter.Start))
}

func ReadChapterTXT(book *database.Book, chapter *database.Chapter) (string, error) {
	var f *os.File
	var err error
	foundedF, found := readerCache.Get(book.Path)
	if !found {
		f, err = os.OpenFile(book.Path, os.O_RDONLY, 0)
		if err != nil {
			return "", err
		}
	} else {
		f = foundedF.(*os.File)
	}

	if _, err := f.Seek(chapter.Start, io.SeekStart); err != nil {
		return "", err
	}
	buf := make([]byte, chapter.Length)
	_, err = io.ReadFull(f, buf)
	if err != nil {
		return "", err
	}
	if book.Coding == "utf-8" {
		return string(buf), nil
	}
	return tools.Gbk2utf8String(buf)
}
