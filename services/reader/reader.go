package reader

import (
	"bufio"
	"errors"
	"io"
	"kolibra/config"
	"kolibra/database"
	"os"
	"time"

	"github.com/gen2brain/go-fitz"
	"github.com/patrickmn/go-cache"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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

type ReaderObject struct {
	F      *os.File
	Reader *bufio.Reader
}

func ReadChapterTXT(book *database.Book, chapter *database.Chapter) (string, error) {
	rdObject, found := readerCache.Get(book.Path)
	if !found {
		f, err := os.OpenFile(book.Path, os.O_RDONLY, 0)
		if err != nil {
			return "", err
		}
		if _, err := f.Seek(chapter.Start, io.SeekStart); err != nil {
			return "", err
		}
		var reader *bufio.Reader
		if book.Coding != "utf-8" {
			transformedReader := transform.NewReader(f, simplifiedchinese.GBK.NewDecoder())
			reader = bufio.NewReader(transformedReader)
		} else {
			reader = bufio.NewReader(f)
		}

		bytes := make([]byte, chapter.Length)
		_, err = io.ReadFull(reader, bytes)
		readerCache.Set(book.Path, &ReaderObject{F: f, Reader: reader}, cache.DefaultExpiration)
		return string(bytes), err
	}

	if _, err := rdObject.(*ReaderObject).F.Seek(chapter.Start, io.SeekStart); err != nil {
		return "", err
	}
	bytes := make([]byte, chapter.Length)
	_, err := io.ReadFull(rdObject.(*ReaderObject).Reader, bytes)
	return string(bytes), err
}
