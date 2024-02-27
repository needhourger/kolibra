package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Database          string             `yaml:"database"`
	Library           string             `yaml:"library"`
	Port              uint               `yaml:"port"`
	Host              string             `yaml:"host"`
	BookExtension     []string           `yaml:"book_extension"`
	FileSizeThreshold int64              `yaml:"file_size_threshold"`
	FileNameMethod    FileNameMethodType `yaml:"file_name_method"`
}

type FileNameMethodType string

const (
	DIR_AUTHOR  FileNameMethodType = "DIR_AUTHOR"
	FILE_AUTHOR FileNameMethodType = "FILE_AUTHOR"
)

var Config *Settings

func load(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(content, &Config)
	log.Printf("Config: %v", Config)
	return err
}

func LoadConfig() {
	log.Printf("Start to load config from ./config.yaml")
	err := load("config.yaml")
	if err != nil {
		log.Panicf("Failed to load config: %s", err)
	}
}
