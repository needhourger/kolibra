package extractor

import (
	"io"
	"kolibra/config"
	"kolibra/database"
	"kolibra/tools"
	"log"
	"regexp"
	"strings"
)

func isStringTitle(reg string, str string) bool {
	isTitle, err := regexp.MatchString(reg, str)
	if err != nil {
		log.Printf("Failed to match title: %s", err)
		return false
	}
	return isTitle
}

func extractTxt(book *database.Book) error {
	txtReader, err := tools.OpenTxtByEncode(book.Path)
	if err != nil {
		return err
	}
	defer txtReader.F.Close()
	var preChapter, curChapter *database.Chapter
	var reg string
	if book.TitleRegex == "" {
		reg = config.Settings.DefaultTitleRegex
	} else {
		reg = book.TitleRegex
	}
	for {
		bytes, _, err := txtReader.Reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		line := string(bytes)
		if isStringTitle(reg, line) {
			pos, err := txtReader.F.Seek(0, io.SeekCurrent)
			if err != nil {
				log.Printf("Failed to get position: %s", err)
				continue
			}
			preChapter = curChapter
			curChapter = &database.Chapter{
				Title: strings.Trim(line, " "),
				// Start:  pos - int64(txtReader.Reader.Size()),
				Start:  pos,
				BookID: book.ID,
			}
			log.Printf("Title: %s, Start: %d", curChapter.Title, curChapter.Start)
			if preChapter != nil {
				preChapter.End = curChapter.Start - 1
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
