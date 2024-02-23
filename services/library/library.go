package library

import (
	"io/fs"
	"kolibra/config"
	"log"
	"path/filepath"
	"strings"
)


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
				}
			}
			return nil
		},
	)
	if err != nil {
		log.Printf("Failed to scan library: %s", err)
	}
}