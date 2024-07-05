package extractor

import (
	"bufio"
	"bytes"
	"kolibra/config"
	DB "kolibra/database"
	"kolibra/tools"
	"log"
	"os"
	"strings"
)

func getBookReg(book *DB.Book) string {
	if book.TitleRegex != "" {
		return book.TitleRegex
	}
	return config.Settings.DefaultTitleRegex
}

func txtSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, data[0 : i+1], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func extractTxt(book *DB.Book) error {
	fileEncoded, err := tools.GetFileEncoded(book.Path)
	if err != nil {
		return err
	}
	log.Printf("Book[%s] encoded: %s", book.Title, fileEncoded)
	isUTF8 := fileEncoded == "utf-8"
	f, err := os.OpenFile(book.Path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	reg := getBookReg(book)

	var currentPos int64 = 0
	var currentChapter, previousChapter *DB.Chapter
	scanner := bufio.NewScanner(f)
	scanner.Split(txtSplitFunc)
	for scanner.Scan() {
		rawBytes := scanner.Bytes()
		bytesLength := len(rawBytes)
		currentPos += int64(bytesLength)
		var line string
		if isUTF8 {
			line = string(rawBytes)
		} else {
			line, _ = tools.Gbk2utf8String(rawBytes)
		}

		// log.Printf("==== %s", line)
		if match, _ := tools.IsMatchString(line, reg); !match {
			continue
		}

		line = strings.Trim(line, " \r\n")
		log.Printf("Found Title: %s", line)

		currentChapter = &DB.Chapter{
			ModelBase: DB.ModelBase{ID: DB.GenerateShortUUID()},
			Title:     line,
			BookID:    book.ID,
			Start:     int64(currentPos),
		}
		if previousChapter != nil {
			currentChapter.PreviousChapterID = previousChapter.ID
			previousChapter.End = currentPos - int64(bytesLength)
			previousChapter.Length = previousChapter.End - previousChapter.Start
			previousChapter.NextChapterID = currentChapter.ID
			DB.CreateChapter(previousChapter)
		}
		previousChapter = currentChapter
	}

	if currentChapter != nil {
		currentChapter.End = currentPos
		currentChapter.Length = currentChapter.End - currentChapter.Start
		DB.CreateChapter(currentChapter)
	}

	book.Coding = fileEncoded
	book.Ready = true
	return DB.UpdateBook(book)
}
