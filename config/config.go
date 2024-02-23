package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Database      string   `yaml:"database"`
	Library       string   `yaml:"library"`
	Port          uint     `yaml:"port"`
	Host          string   `yaml:"host"`
	BookExtension []string `yaml:"book_extension"`
}

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
