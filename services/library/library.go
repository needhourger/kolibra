package library

import (
	"io/fs"
	"kolibra/config"
	"kolibra/database"
	"kolibra/services/extractor"
	"kolibra/tools"
	"log"
	"path/filepath"
	"strings"
)

func extractFileName(path string, info fs.FileInfo) (string, string) {
	fileName := strings.Trim(strings.Split(info.Name(), ".")[0], " ")
	switch config.Config.FileNameMethod {
	case config.DIR_AUTHOR:
		return filepath.Base(filepath.Dir(path)), fileName
	case config.FILE_AUTHOR:
		author := strings.Trim(strings.Split(fileName, "-")[0], " ")
		name := strings.Trim(strings.Split(fileName, "-")[1], " ")
		return author, name
	default:
		return "", fileName
	}
}

func LoadBookByPath(path string, info fs.FileInfo) {
	hash, err := tools.CalculateFileHash(path)
	if err != nil {
		log.Printf("Failed to calculate file hash: %s", err)
		return
	}
	if database.CheckBookFileHash(hash) {
		log.Printf("Book already exists: %s", path)
		return
	}
	author, title := extractFileName(path, info)
	book := database.Book{
		Title:     title,
		Author:    author,
		Extension: filepath.Ext(info.Name()),
		Size:      info.Size(),
		Path:      path,
		Ready:     false,
		Hash:      hash,
	}
	err = database.CreateBook(&book)
	if err != nil {
		log.Printf("Failed to load book: %s", err)
		return
	}

	err = extractor.Extract(&book)
	if err != nil {
		log.Printf("Failed to extract book: %s", err)
	}
}

func ScanLibrary() {
	err := filepath.Walk(
		config.Config.Library,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			for _, suffix := range config.Config.BookExtension {
				if strings.HasSuffix(info.Name(), suffix) {
					log.Printf("Found book: %s", path)
					LoadBookByPath(path, info)
				}
			}
			return nil
		},
	)
	if err != nil {
		log.Printf("Failed to scan library: %s", err)
	}
}
