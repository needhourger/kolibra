package reader

import (
	"errors"
	"io"
	"kolibra/config"
	"kolibra/database/model"
	"kolibra/tools"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
)

var readerCache *cache.Cache

var (
	ErrFileType          = errors.New("Pool already exists!")
	ErrReaderPoolExisted = errors.New("Pool already exists!")
)

func CreateReaderCache() error {
	if readerCache == nil {
		readerCache = cache.New(
			time.Duration(config.Settings.Advance.ReaderCachedMinutes)*time.Minute,
			time.Duration(config.Settings.Advance.ReaderCachedMinutes)*time.Minute,
		)
		return nil
	}
	return ErrReaderPoolExisted
}

func ReadChapter(book *model.Book, chapter *model.Chapter) (any, error) {
	switch book.Extension {
	case ".txt":
		return ReadChapterTXT(book, chapter)
	}
	return "", ErrFileType
}

func ReadChapterTXT(book *model.Book, chapter *model.Chapter) (string, error) {
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
	if book.Encoding == "utf-8" {
		return string(buf), nil
	}
	return tools.Gbk2utf8String(buf)
}
