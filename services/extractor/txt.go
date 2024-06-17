package extractor

import (
	"bufio"
	"io"
	"kolibra/config"
	"kolibra/database"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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
	f, err := os.OpenFile(book.Path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	dumpedBytes, err := reader.Peek(1024)
	f.Seek(0, io.SeekStart)
	reader.Reset(f)
	if err != nil {
		return err
	}

	var scanner *bufio.Scanner
	_, codingName, _ := charset.DetermineEncoding(dumpedBytes, "text/plain")
	log.Printf("Detected coding: %s", codingName)
	if codingName != "utf-8" {
		transformedReader := transform.NewReader(f, simplifiedchinese.GBK.NewDecoder())
		scanner = bufio.NewScanner(transformedReader)
	} else {
		scanner = bufio.NewScanner(reader)
	}

	var preChapter, curChapter *database.Chapter
	var reg string
	if book.TitleRegex == "" {
		reg = config.Settings.DefaultTitleRegex
	} else {
		reg = book.TitleRegex
	}
	for scanner.Scan() {
		line := scanner.Text()
		if isStringTitle(reg, line) {
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

	book.Coding = codingName
	book.Ready = true
	err = database.UpdateBook(book)
	if err != nil {
		return err
	}
	return nil
}
