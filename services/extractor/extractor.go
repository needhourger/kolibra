package extractor

import (
	"bufio"
	"errors"
	"io"
	"kolibra/config"
	"kolibra/database"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

func Extract(book *database.Book) error {
	switch book.Extension {
	case ".txt":
		return extractTxt(book)
	}
	return errors.New("unsupported file type")
}

func isStringTitle(reg string, str string) bool {
	isTitle, err := regexp.MatchString(reg, str)
	if err != nil {
		log.Printf("Failed to match title: %s", err)
		return false
	}
	return isTitle
}

func openFile(path string) (*bufio.Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(f)
	dumpedBytes, err := reader.Peek(1024)
	if err != nil {
		return nil, err
	}
	e, name, _ := charset.DetermineEncoding(dumpedBytes, "text/plain")
	if name == "utf-8" {
		return reader, nil
	}
	newDecodedReader := transform.NewReader(f, e.NewDecoder().Transformer)
	return bufio.NewReader(newDecodedReader), nil
}

func extractTxt(book *database.Book) error {
	reader, err := openFile(book.Path)
	if err != nil {
		return err
	}
	var preChapter, curChapter *database.Chapter
	for {
		bytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		line := string(bytes)
		var reg string
		if book.TitleRegex == "" {
			reg = config.Config.DefaultTitleRegex
		} else {
			reg = book.TitleRegex
		}

		isTitle := isStringTitle(reg, line)
		if isTitle {
			pos, err := f.Seek(0, io.SeekCurrent)
			if err != nil {
				log.Printf("Failed to get position: %s", err)
				continue
			}
			preChapter = curChapter
			curChapter = &database.Chapter{
				Title:  strings.Trim(line, " "),
				Start:  pos,
				BookID: book.ID,
			}
			if preChapter != nil {
				preChapter.End = pos - 1
				preChapter.Length = preChapter.End - preChapter.Start
				err = database.CreateChapter(preChapter)
				if err != nil {
					log.Printf("Failed to create chapter: %s", err)
				}
			}
		}
	}
	if curChapter != nil {
		curChapter.End = book.Size
		curChapter.Length = curChapter.End - curChapter.Start
		err = database.CreateChapter(curChapter)
		if err != nil {
			log.Printf("Failed to create chapter: %s", err)
		}
	}

	book.Ready = true
	err = database.UpdateBook(book)
	if err != nil {
		return err
	}
	return nil
}
