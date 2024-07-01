package extractor

import (
	"bufio"
	"kolibra/config"
	"kolibra/database"
	"kolibra/tools"
	"log"
	"os"
	"strings"
)

func getBookReg(book *database.Book) string {
	if book.TitleRegex != "" {
		return book.TitleRegex
	}
	return config.Settings.DefaultTitleRegex
}

func extractTxt(book *database.Book) error {
	fileEncoded, err := tools.GetFileEncoded(book.Path)
	if err != nil {
		return err
	}
	isUTF8 := fileEncoded == "utf-8"
	f, err := os.OpenFile(book.Path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	reg := getBookReg(book)

	var currentPos int64 = 0
	var currentChapter, previousChapter *database.Chapter
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		rawBytes := scanner.Bytes()
		bytesLength := len(rawBytes)
		currentPos += int64(bytesLength)
		var line string
		if isUTF8 {
			line, _ = tools.Gbk2utf8String(rawBytes)
		} else {
			line = string(rawBytes)
		}

		if match, _ := tools.IsMatchString(line, reg); !match {
			continue
		}

		line = strings.Trim(line, " ")
		log.Printf("Found Title: %s", line)

		currentChapter = &database.Chapter{
			Title:  line,
			BookID: book.ID,
			Start:  int64(currentPos),
		}
		if previousChapter != nil {
			previousChapter.End = currentPos - int64(bytesLength)
			database.CreateChapter(previousChapter)
		}
		previousChapter = currentChapter
	}

	if currentChapter != nil {
		currentChapter.End = currentPos
		database.CreateChapter(currentChapter)
	}

	book.Coding = fileEncoded
	book.Ready = true
	return database.UpdateBook(book)
}
