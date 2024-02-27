package library

import (
	"io/fs"
	"kolibra/config"
	"kolibra/database"
	"log"
	"path/filepath"
	"strings"
)

func extractFileName(path string, info fs.FileInfo) (string, string) {
	switch config.Config.FileNameMethod {
	case config.DIR_AUTHOR:
		return filepath.Base(filepath.Dir(path)), info.Name()
	case config.FILE_AUTHOR:
		return strings.Split(info.Name(), "-")[0], strings.Split(info.Name(), "-")[1]
	default:
		return "", info.Name()
	}
}

func LoadBookByPath(path string, info fs.FileInfo) {
	author, title := extractFileName(path, info)
	book := database.Book{
		Title:     title,
		Author:    author,
		Extension: filepath.Ext(info.Name()),
		Size:      info.Size(),
		Path:      path,
		Ready:     false,
	}
	err := database.CreateBook(&book)
	if err != nil {
		log.Printf("Failed to load book: %s", err)
		return
	}

	// todo: extract chapters
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
