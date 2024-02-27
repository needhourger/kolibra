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

func extractTxt(book *database.Book) error {
	f, err := os.Open(book.Path)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(f)
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
